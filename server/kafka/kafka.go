package kafka

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/contractproposal"
	"github.com/loomi-labs/star-scope/ent/event"
	"github.com/loomi-labs/star-scope/ent/proposal"
	"github.com/loomi-labs/star-scope/grpc/event/eventpb"
	"github.com/loomi-labs/star-scope/indexevent"
	"github.com/loomi-labs/star-scope/queryevent"
	"github.com/segmentio/kafka-go"
	"github.com/shifty11/go-logger/log"
	"golang.org/x/exp/slices"
	"time"
)

var (
	indexEventsTopic     = "index-events"
	queryEventsTopic     = "query-events"
	processedEventsTopic = "processed-events"
)

type Kafka struct {
	addresses            []string
	chainManager         *database.ChainManager
	eventListenerManager *database.EventListenerManager
}

func NewKafka(dbManager *database.DbManagers, addresses ...string) *Kafka {
	return &Kafka{
		addresses:            addresses,
		chainManager:         dbManager.ChainManager,
		eventListenerManager: dbManager.EventListenerManager,
	}
}

func (k *Kafka) indexedEventsReader() *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   k.addresses,
		Topic:     indexEventsTopic,
		GroupID:   indexEventsTopic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		MaxWait:   1 * time.Second,
	})
	return r
}

func (k *Kafka) queryEventsReader() *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   k.addresses,
		Topic:     queryEventsTopic,
		GroupID:   queryEventsTopic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		MaxWait:   1 * time.Second,
	})
	return r
}

func (k *Kafka) processedEventsReader() *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   k.addresses,
		Topic:     processedEventsTopic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		MaxWait:   1 * time.Second,
	})
	return r
}

func (k *Kafka) closeReader(r *kafka.Reader) {
	err := r.Close()
	if err != nil {
		log.Sugar.Fatal("failed to closeReader writer:", err)
	}
}

func (k *Kafka) writer() *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(k.addresses...),
		Topic:    processedEventsTopic,
		Balancer: &kafka.LeastBytes{},
	}
}

func (k *Kafka) closeWriter(w *kafka.Writer) {
	err := w.Close()
	if err != nil {
		log.Sugar.Fatal("failed to close writer:", err)
	}
}

func (k *Kafka) produce(msgs [][]byte) {
	w := k.writer()
	defer k.closeWriter(w)

	kafkaMsgs := make([]kafka.Message, len(msgs))
	for i, msg := range msgs {
		kafkaMsgs[i] = kafka.Message{Value: msg}
	}

	err := w.WriteMessages(context.Background(), kafkaMsgs...)
	if err != nil {
		panic(fmt.Sprintf("failed to write messages: %v", err))
	}
}

func (k *Kafka) getEventListenerMapForWallets() map[string]*ent.EventListener {
	var elMap = make(map[string]*ent.EventListener)
	for _, el := range k.eventListenerManager.QueryAll(context.Background()) {
		elMap[el.WalletAddress] = el
	}
	return elMap
}

func (k *Kafka) ProcessIndexedEvents() {
	log.Sugar.Info("Start consuming indexed events")
	r := k.indexedEventsReader()
	defer k.closeReader(r)

	elMap := k.getEventListenerMapForWallets()
	elMapUpdatedAt := time.Now()

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Sugar.Errorf("failed to read message: %v", err)
		}
		var txEvent indexevent.TxEvent
		err = proto.Unmarshal(msg.Value, &txEvent)
		if err != nil {
			log.Sugar.Error(err)
		} else {
			if time.Since(elMapUpdatedAt) > 5*time.Minute {
				elMap = k.getEventListenerMapForWallets()
				elMapUpdatedAt = time.Now()
			}
			if el, ok := elMap[txEvent.WalletAddress]; ok {
				var ctx = context.Background()
				var err2 error
				log.Sugar.Debugf("Processing event %v with address %v", msg.Offset, txEvent.WalletAddress)
				switch txEvent.GetEvent().(type) {
				case *indexevent.TxEvent_CoinReceived:
					_, err2 = k.eventListenerManager.UpdateAddEvent(ctx, el, event.EventTypeFUNDING, event.DataTypeTxEvent_CoinReceived, txEvent.NotifyTime.AsTime(), msg.Value, true)
				case *indexevent.TxEvent_OsmosisPoolUnlock:
					_, err2 = k.eventListenerManager.UpdateAddEvent(ctx, el, event.EventTypeDEX, event.DataTypeTxEvent_OsmosisPoolUnlock, txEvent.NotifyTime.AsTime(), msg.Value, true)
				case *indexevent.TxEvent_Unstake:
					_, err2 = k.eventListenerManager.UpdateAddEvent(ctx, el, event.EventTypeSTAKING, event.DataTypeTxEvent_Unstake, txEvent.NotifyTime.AsTime(), msg.Value, true)
				}
				if err2 != nil {
					log.Sugar.Errorf("failed to update event for %v: %v", txEvent.WalletAddress, err2)
				} else {
					if txEvent.NotifyTime.AsTime().Before(time.Now()) {
						k.produce([][]byte{msg.Value})
						log.Sugar.Debugf("Put event %v with address %v to `%v`", msg.Offset, txEvent.WalletAddress, processedEventsTopic)
					}
				}
			} else {
				log.Sugar.Debugf("Discard event %v with address %v", msg.Offset, txEvent.WalletAddress)
			}
		}
	}
}

func (k *Kafka) getEventListenerMapForChains() map[uint64][]*ent.EventListener {
	var elMap = make(map[uint64][]*ent.EventListener)
	for _, el := range k.eventListenerManager.QueryAllWithChain(context.Background()) {
		if _, ok := elMap[uint64(el.Edges.Chain.ID)]; !ok {
			elMap[uint64(el.Edges.Chain.ID)] = make([]*ent.EventListener, 0)
		}
		elMap[uint64(el.Edges.Chain.ID)] = append(elMap[uint64(el.Edges.Chain.ID)], el)
	}
	return elMap
}

func getProposalDataType(prop *ent.Proposal) event.DataType {
	switch prop.Status {
	case proposal.StatusPROPOSAL_STATUS_UNSPECIFIED, proposal.StatusPROPOSAL_STATUS_DEPOSIT_PERIOD, proposal.StatusPROPOSAL_STATUS_VOTING_PERIOD:
		return event.DataTypeQueryEvent_GovernanceProposal_Ongoing
	case proposal.StatusPROPOSAL_STATUS_PASSED, proposal.StatusPROPOSAL_STATUS_REJECTED, proposal.StatusPROPOSAL_STATUS_FAILED:
		return event.DataTypeQueryEvent_GovernanceProposal_Finished
	}
	log.Sugar.Panicf("Unknown proposal status: %v", prop.Status)
	return event.DataTypeQueryEvent_GovernanceProposal_Ongoing
}

func getContractProposalDataType(prop *ent.ContractProposal) event.DataType {
	switch prop.Status {
	case contractproposal.StatusOPEN:
		return event.DataTypeQueryEvent_GovernanceProposal_Ongoing
	case contractproposal.StatusREJECTED, contractproposal.StatusPASSED, contractproposal.StatusEXECUTED, contractproposal.StatusCLOSED, contractproposal.StatusEXECUTION_FAILED:
		return event.DataTypeQueryEvent_GovernanceProposal_Finished
	}
	log.Sugar.Panicf("Unknown proposal status: %v", prop.Status)
	return event.DataTypeQueryEvent_GovernanceProposal_Ongoing
}

func (k *Kafka) getChains() map[uint64]*ent.Chain {
	var chains = make(map[uint64]*ent.Chain)
	for _, chain := range k.chainManager.QueryAll(context.Background()) {
		chains[uint64(chain.ID)] = chain
	}
	return chains
}

func (k *Kafka) getChains() map[uint64]*ent.Chain {
	var chains = make(map[uint64]*ent.Chain)
	for _, chain := range k.chainManager.QueryAll(context.Background()) {
		chains[uint64(chain.ID)] = chain
	}
	return chains
}

func (k *Kafka) ProcessQueryEvents() {
	log.Sugar.Info("Start consuming query events")
	r := k.queryEventsReader()
	defer k.closeReader(r)

	elMap := k.getEventListenerMapForChains()
	elMapUpdatedAt := time.Now()
	chains := k.getChains()

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Sugar.Errorf("failed to read message: %v", err)
		}
		var queryEvent queryevent.QueryEvent
		err = proto.Unmarshal(msg.Value, &queryEvent)
		if err != nil {
			log.Sugar.Error(err)
		} else {
			if time.Now().Sub(elMapUpdatedAt) > 5*time.Minute {
				elMap = k.getEventListenerMapForChains()
				elMapUpdatedAt = time.Now()
			}
			chain, ok := chains[queryEvent.ChainId]
			if !ok {
				log.Sugar.Panicf("failed to find chain %v", queryEvent.ChainId)
			}
			log.Sugar.Debugf("Processing event %v for chain %v", msg.Offset, chain.PrettyName)
			var ctx = context.Background()
			var pbEvents [][]byte
			switch queryEvent.GetEvent().(type) {
			case *queryevent.QueryEvent_GovernanceProposal:
				prop, err := k.chainManager.UpdateProposal(ctx, chain, queryEvent.GetGovernanceProposal())
				if err != nil {
					log.Sugar.Panicf("failed to update prop %v: %v", queryEvent.GetGovernanceProposal().ProposalId, err)
				}
				var ignoredStates = []proposal.Status{
					proposal.StatusPROPOSAL_STATUS_UNSPECIFIED,
					proposal.StatusPROPOSAL_STATUS_DEPOSIT_PERIOD,
				}
				if slices.Contains(ignoredStates, prop.Status) {
					continue
				}
				if els, ok := elMap[queryEvent.ChainId]; ok {
					for _, el := range els {
						var err2 error
						log.Sugar.Debugf("Processing event %v with address %v for %v", msg.Offset, queryEvent.ChainId, el.WalletAddress)
						_, err2 = k.eventListenerManager.UpdateAddEvent(ctx, el, event.EventTypeGOVERNANCE, getProposalDataType(prop), queryEvent.NotifyTime.AsTime(), msg.Value, false)
						if err2 != nil {
							log.Sugar.Errorf("failed to update event for %v: %v", el.WalletAddress, err2)
						} else {
							if queryEvent.NotifyTime.AsTime().Before(time.Now()) {
								pbEvents = append(pbEvents, msg.Value)
								log.Sugar.Debugf("Put event %v with address %v to `%v`", msg.Offset, el.WalletAddress, processedEventsTopic)
							}
						}
					}
				}
			case *queryevent.QueryEvent_ContractGovernanceProposal:
				prop, err := k.chainManager.UpdateContractProposal(ctx, chain, queryEvent.GetContractGovernanceProposal())
				if err != nil {
					log.Sugar.Panicf("failed to update prop %v: %v", queryEvent.GetContractGovernanceProposal().ProposalId, err)
				}
				if els, ok := elMap[queryEvent.ChainId]; ok {
					for _, el := range els {
						var err2 error
						log.Sugar.Debugf("Processing event %v with address %v for %v", msg.Offset, queryEvent.ChainId, el.WalletAddress)
						_, err2 = k.eventListenerManager.UpdateAddEvent(ctx, el, event.EventTypeGOVERNANCE, getContractProposalDataType(prop), queryEvent.NotifyTime.AsTime(), msg.Value, false)
						if err2 != nil {
							log.Sugar.Errorf("failed to update event for %v: %v", el.WalletAddress, err2)
						} else {
							if queryEvent.NotifyTime.AsTime().Before(time.Now()) {
								pbEvents = append(pbEvents, msg.Value)
								log.Sugar.Debugf("Put event %v with address %v to `%v`", msg.Offset, el.WalletAddress, processedEventsTopic)
							}
						}
					}
				}
			}
			if len(pbEvents) > 0 {
				log.Sugar.Debugf("Produce %v events", len(pbEvents))
				k.produce(pbEvents)
			}
		}
	}
}

func (k *Kafka) ConsumeProcessedEvents(ctx context.Context, user *ent.User, eventsChannel chan *eventpb.EventList) {
	log.Sugar.Debugf("Start processed-events consumer for user %v", user.WalletAddress)
	r := k.processedEventsReader()
	defer k.closeReader(r)

	err := r.SetOffsetAt(context.Background(), time.Now())
	if err != nil {
		log.Sugar.Errorf("failed to set offset: %v", err)
		eventsChannel <- nil
	}

	chains := k.chainManager.QueryAll(context.Background())

	for {
		select {
		case <-ctx.Done():
			log.Sugar.Debugf("Stop the processed-events consumer for user %v", user.WalletAddress)
			return
		default:
			msg, err := r.ReadMessage(context.Background())
			if err != nil {
				break
			}
			txEvent, err := kafkaMsgToProto(msg.Value, chains)
			if err != nil {
				log.Sugar.Error(err)
			}
			eventsChannel <- &eventpb.EventList{Events: []*eventpb.Event{txEvent}}
		}
	}
}
