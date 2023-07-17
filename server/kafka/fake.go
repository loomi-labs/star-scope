package kafka

import (
	"context"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/event"
	kafkaevent "github.com/loomi-labs/star-scope/event"
	"github.com/loomi-labs/star-scope/grpc/event/eventpb"
	"github.com/shifty11/go-logger/log"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"time"
)

type FakeEventCreator struct {
	walletAddresses      []string
	eventListenerManager *database.EventListenerManager
	chainManager         *database.ChainManager
	kafka                *Kafka
}

func NewFakeEventCreator(dbManager *database.DbManagers, walletAddresses []string, kafkaBrokers ...string) *FakeEventCreator {
	return &FakeEventCreator{
		walletAddresses:      walletAddresses,
		eventListenerManager: dbManager.EventListenerManager,
		chainManager:         dbManager.ChainManager,
		kafka:                NewKafka(dbManager, kafkaBrokers),
	}
}

func (d *FakeEventCreator) convertedAddresses(bech32Prefix string) []string {
	var convertedAddresses = []string{}
	for _, address := range d.walletAddresses {
		convertedAddress, err := common.ConvertWithOtherPrefix(address, bech32Prefix)
		if err != nil {
			log.Sugar.Panicf("error converting wallet address: %v", err)
		}
		convertedAddresses = append(convertedAddresses, convertedAddress)
	}
	return convertedAddresses
}

func (d *FakeEventCreator) getEventListenerMap() map[string]*ent.EventListener {
	var elMap = make(map[string]*ent.EventListener)
	for _, el := range d.eventListenerManager.QueryWithChain(context.Background()) {
		if slices.Contains(d.convertedAddresses(el.Edges.Chain.Bech32Prefix), el.WalletAddress) {
			elMap[el.WalletAddress] = el
		}
	}
	return elMap
}

func createWalletEvent(mainAddress string, chains []*ent.Chain) kafkaevent.WalletEvent {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	chain := chains[r.Intn(len(chains))]
	walletAddress, err := common.ConvertWithOtherPrefix(mainAddress, chain.Bech32Prefix)
	if err != nil {
		log.Sugar.Panicf("error converting wallet address: %v", err)
	}
	return kafkaevent.WalletEvent{
		ChainId:       uint64(chain.ID),
		WalletAddress: walletAddress,
		Timestamp:     timestamppb.New(time.Now()),
		NotifyTime:    timestamppb.New(time.Now()),
	}
}

func createCoinReceivedEvent(walletAddress string, chains []*ent.Chain) *kafkaevent.WalletEvent {
	var walletEvent = createWalletEvent(walletAddress, chains)
	walletEvent.Event = &kafkaevent.WalletEvent_CoinReceived{
		CoinReceived: &kafkaevent.CoinReceivedEvent{
			Sender: "cosmos1h872wxm58laz23rld32hlsqq6067j257txh8j6",
			Coin:   &kafkaevent.Coin{Denom: "uatom", Amount: "1000000"},
		},
	}
	return &walletEvent
}

func createUnstakeEvent(walletAddress string, chains []*ent.Chain) *kafkaevent.WalletEvent {
	var txEvent = createWalletEvent(walletAddress, chains)
	txEvent.Event = &kafkaevent.WalletEvent_Unstake{
		Unstake: &kafkaevent.UnstakeEvent{
			CompletionTime: timestamppb.New(time.Now()),
			Coin:           &kafkaevent.Coin{Denom: "uatom", Amount: "1000000"},
		},
	}
	return &txEvent
}

func createOsmoPoolUnlockEvent(walletAddress string, chains []*ent.Chain) *kafkaevent.WalletEvent {
	var txEvent = createWalletEvent(walletAddress, chains)
	txEvent.Event = &kafkaevent.WalletEvent_OsmosisPoolUnlock{
		OsmosisPoolUnlock: &kafkaevent.OsmosisPoolUnlockEvent{
			Duration:   durationpb.New(time.Hour * 24 * 7),
			UnlockTime: timestamppb.New(time.Now()),
		},
	}
	return &txEvent
}

func (d *FakeEventCreator) createRandomEvents(chains []*ent.Chain) []*kafkaevent.WalletEvent {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var walletEvents []*kafkaevent.WalletEvent
	for _, walletAddress := range d.walletAddresses {
		var walletEvent kafkaevent.WalletEvent
		randNbr := r.Intn(3)
		switch randNbr {
		case 0:
			walletEvent = *createCoinReceivedEvent(walletAddress, chains)
		case 1:
			walletEvent = *createUnstakeEvent(walletAddress, chains)
		case 2:
			walletEvent = *createOsmoPoolUnlockEvent(walletAddress, chains)
		}
		log.Sugar.Debugf("create random event %v for %v", randNbr, walletEvent.WalletAddress)
		walletEvents = append(walletEvents, &walletEvent)
	}
	return walletEvents
}

func (d *FakeEventCreator) CreateFakeEvents() {
	log.Sugar.Info("Start creating fake events")

	chains := d.chainManager.QueryEnabled(context.Background())

	for {
		elMap := d.getEventListenerMap()
		msgs := d.createRandomEvents(chains)
		for _, msg := range msgs {
			if el, ok := elMap[msg.WalletAddress]; ok {
				var eventType event.EventType
				var dataType event.DataType
				switch msg.GetEvent().(type) {
				case *kafkaevent.WalletEvent_CoinReceived:
					eventType, dataType = event.EventTypeFUNDING, event.DataTypeWalletEvent_CoinReceived
				case *kafkaevent.WalletEvent_OsmosisPoolUnlock:
					eventType, dataType = event.EventTypeDEX, event.DataTypeWalletEvent_OsmosisPoolUnlock
				case *kafkaevent.WalletEvent_Unstake:
					eventType, dataType = event.EventTypeSTAKING, event.DataTypeWalletEvent_Unstake
				}
				_, err := d.eventListenerManager.UpdateAddWalletEvent(context.Background(), el, msg, eventType, dataType, false)
				if err != nil {
					log.Sugar.Panicf("failed to update event for %v: %v", msg.WalletAddress, err)
				} else {
					if msg.NotifyTime.AsTime().Before(time.Now()) {
						d.kafka.produceProcessedEvents([]*kafkaevent.EventProcessedMsg{{
							WalletAddress: msg.WalletAddress,
							EventType:     eventpb.EventType(eventpb.EventType_value[eventType.String()]),
						}})
					}
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
}
