package kafka

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/event"
	"github.com/loomi-labs/star-scope/indexevent"
	"github.com/shifty11/go-logger/log"
	"golang.org/x/exp/slices"
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
		kafka:                NewKafka(dbManager, kafkaBrokers...),
	}
}

func (d *FakeEventCreator) getEventListenerMap() map[string]*ent.EventListener {
	var elMap = make(map[string]*ent.EventListener)
	for _, el := range d.eventListenerManager.QueryAll(context.Background()) {
		if slices.Contains(d.walletAddresses, el.WalletAddress) {
			elMap[el.WalletAddress] = el
		}
	}
	return elMap
}

func createTxEvent(walletAddress string, chains []*ent.Chain) indexevent.TxEvent {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	chain := chains[r.Intn(len(chains))]
	return indexevent.TxEvent{
		ChainId:       uint64(chain.ID),
		WalletAddress: walletAddress,
		Timestamp:     timestamppb.New(time.Now()),
		NotifyTime:    timestamppb.New(time.Now()),
	}
}

func createCoinReceivedEvent(walletAddress string, chains []*ent.Chain) *indexevent.TxEvent {
	var txEvent = createTxEvent(walletAddress, chains)
	txEvent.Event = &indexevent.TxEvent_CoinReceived{
		CoinReceived: &indexevent.CoinReceivedEvent{
			Sender: "cosmos1h872wxm58laz23rld32hlsqq6067j257txh8j6",
			Coin:   &indexevent.Coin{Denom: "uatom", Amount: "1000000"},
		},
	}
	return &txEvent
}

func createUnstakeEvent(walletAddress string, chains []*ent.Chain) *indexevent.TxEvent {
	var txEvent = createTxEvent(walletAddress, chains)
	txEvent.Event = &indexevent.TxEvent_Unstake{
		Unstake: &indexevent.UnstakeEvent{
			CompletionTime: timestamppb.New(time.Now()),
			Coin:           &indexevent.Coin{Denom: "uatom", Amount: "1000000"},
		},
	}
	return &txEvent
}

func createOsmoPoolUnlockEvent(walletAddress string, chains []*ent.Chain) *indexevent.TxEvent {
	var txEvent = createTxEvent(walletAddress, chains)
	txEvent.Event = &indexevent.TxEvent_Unstake{
		Unstake: &indexevent.UnstakeEvent{
			CompletionTime: timestamppb.New(time.Now()),
			Coin:           &indexevent.Coin{Denom: "uatom", Amount: "1000000"},
		},
	}
	return &txEvent
}

func (d *FakeEventCreator) createRandomTxEvents(chains []*ent.Chain) []*indexevent.TxEvent {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var txEvents []*indexevent.TxEvent
	for _, walletAddress := range d.walletAddresses {
		var txEvent indexevent.TxEvent
		randNbr := r.Intn(3)
		switch randNbr {
		case 0:
			txEvent = *createCoinReceivedEvent(walletAddress, chains)
		case 1:
			txEvent = *createUnstakeEvent(walletAddress, chains)
		case 2:
			txEvent = *createOsmoPoolUnlockEvent(walletAddress, chains)
		}
		log.Sugar.Debugf("create random event %v", randNbr)
		txEvents = append(txEvents, &txEvent)
	}
	return txEvents
}

func (d *FakeEventCreator) CreateFakeEvents() {
	log.Sugar.Info("Start creating fake events")

	chains := d.chainManager.QueryEnabled(context.Background())

	for {
		elMap := d.getEventListenerMap()
		msgs := d.createRandomTxEvents(chains)
		for _, msg := range msgs {
			if el, ok := elMap[msg.WalletAddress]; ok {
				var ctx = context.Background()
				var err error
				byteMsg, err := proto.Marshal(msg)
				if err != nil {
					log.Sugar.Panicf("failed to byteMsg event: %v", err)
				}
				switch msg.GetEvent().(type) {
				case *indexevent.TxEvent_CoinReceived:
					_, err = d.eventListenerManager.UpdateAddEvent(ctx, el, event.EventTypeFUNDING, event.DataTypeTxEvent_CoinReceived, msg.NotifyTime.AsTime(), byteMsg, true)
				case *indexevent.TxEvent_OsmosisPoolUnlock:
					_, err = d.eventListenerManager.UpdateAddEvent(ctx, el, event.EventTypeDEX, event.DataTypeTxEvent_OsmosisPoolUnlock, msg.NotifyTime.AsTime(), byteMsg, true)
				case *indexevent.TxEvent_Unstake:
					_, err = d.eventListenerManager.UpdateAddEvent(ctx, el, event.EventTypeSTAKING, event.DataTypeTxEvent_Unstake, msg.NotifyTime.AsTime(), byteMsg, true)
				}
				if err != nil {
					log.Sugar.Panicf("failed to update event for %v: %v", msg.WalletAddress, err)
				} else {
					if msg.NotifyTime.AsTime().Before(time.Now()) {
						d.kafka.produce([][]byte{byteMsg})
					}
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}
