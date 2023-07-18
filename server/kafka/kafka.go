package kafka

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/contractproposal"
	"github.com/loomi-labs/star-scope/ent/event"
	"github.com/loomi-labs/star-scope/ent/eventlistener"
	"github.com/loomi-labs/star-scope/ent/proposal"
	kafkaevent "github.com/loomi-labs/star-scope/event"
	"github.com/loomi-labs/star-scope/grpc/event/eventpb"
	"github.com/loomi-labs/star-scope/kafka_internal"
	"github.com/loomi-labs/star-scope/types"
	"github.com/segmentio/kafka-go"
	"github.com/shifty11/go-logger/log"
	"golang.org/x/exp/slices"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type GroupTopic types.Topic

const (
	chainEventsTopic    = GroupTopic(types.ChainEventsTopic)
	contractEventsTopic = GroupTopic(types.ContractEventsTopic)
	walletEventsTopic   = GroupTopic(types.WalletEventsTopic)
)

type Topic string

const (
	processedEventsTopic = Topic(types.ProcessedEventsTopic)
)

type Kafka struct {
	addresses            []string
	chainManager         *database.ChainManager
	eventListenerManager *database.EventListenerManager
}

func NewKafka(dbManager *database.DbManagers, addresses []string) *Kafka {
	return &Kafka{
		addresses:            addresses,
		chainManager:         dbManager.ChainManager,
		eventListenerManager: dbManager.EventListenerManager,
	}
}

func (k *Kafka) groupReader(topic GroupTopic) *kafka.Reader {
	log.Sugar.Infof("Start consuming %v", topic)
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   k.addresses,
		Topic:     string(topic),
		GroupID:   string(topic),
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		MaxWait:   1 * time.Second,
	})
}

func (k *Kafka) reader(topic Topic) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   k.addresses,
		Topic:     string(topic),
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		MaxWait:   1 * time.Second,
	})
}

func (k *Kafka) closeReader(r *kafka.Reader) {
	err := r.Close()
	if err != nil {
		log.Sugar.Fatal("failed to closeReader writer:", err)
	}
}

func (k *Kafka) writer(topic Topic) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(k.addresses...),
		Topic:    string(topic),
		Balancer: &kafka.LeastBytes{},
	}
}

func (k *Kafka) closeWriter(w *kafka.Writer) {
	err := w.Close()
	if err != nil {
		log.Sugar.Fatal("failed to close writer:", err)
	}
}

func (k *Kafka) produceProcessedEvents(msgs []*kafkaevent.EventProcessedMsg) {
	w := k.writer(processedEventsTopic)
	defer k.closeWriter(w)

	kafkaMsgs := make([]kafka.Message, len(msgs))
	for i, msg := range msgs {
		value, err := proto.Marshal(msg)
		if err != nil {
			log.Sugar.Panicf("failed to marshal message: %v", err)
		}
		kafkaMsgs[i] = kafka.Message{Value: value}
	}

	err := w.WriteMessages(context.Background(), kafkaMsgs...)
	if err != nil {
		log.Sugar.Panicf("failed to write messages: %v", err)
	}
}

func (k *Kafka) getEventListenerMapForWallets() map[string]map[event.DataType]*ent.EventListener {
	var elMap = make(map[string]map[event.DataType]*ent.EventListener)
	var dt = []eventlistener.DataType{
		eventlistener.DataTypeWalletEvent_Unstake,
		eventlistener.DataTypeWalletEvent_CoinReceived,
		eventlistener.DataTypeWalletEvent_OsmosisPoolUnlock,
		eventlistener.DataTypeWalletEvent_NeutronTokenVesting,
		eventlistener.DataTypeWalletEvent_Voted,
		eventlistener.DataTypeWalletEvent_VoteReminder,
	}
	for _, el := range k.eventListenerManager.Query(context.Background(), dt...) {
		if _, ok := elMap[el.WalletAddress]; !ok {
			elMap[el.WalletAddress] = make(map[event.DataType]*ent.EventListener)
		}
		elMap[el.WalletAddress][event.DataType(el.DataType)] = el
	}
	return elMap
}

func (k *Kafka) ProcessWalletEvents() {
	r := k.groupReader(walletEventsTopic)
	defer k.closeReader(r)

	elMap := k.getEventListenerMapForWallets()
	elMapMutex := sync.Mutex{} // Mutex to synchronize access to elMap

	go func() {
		kafkaInternal := kafka_internal.NewKafkaInternal(k.addresses)

		ch := make(chan kafka_internal.DbChange)
		defer close(ch)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go kafkaInternal.ReadDbChanges(ctx, ch, []kafka_internal.DbChange{kafka_internal.EventListenerCreated, kafka_internal.EventListenerDeleted})

		for range ch {
			log.Sugar.Debugf("Updating event listener map")
			elMapMutex.Lock()
			elMap = k.getEventListenerMapForWallets()
			elMapMutex.Unlock()
		}
	}()

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Sugar.Errorf("failed to read message: %v", err)
		}
		var walletEvent kafkaevent.WalletEvent
		err = proto.Unmarshal(msg.Value, &walletEvent)
		if err != nil {
			log.Sugar.Error(err)
		} else {
			elMapMutex.Lock()
			if els, ok := elMap[walletEvent.WalletAddress]; ok {
				log.Sugar.Debugf("Processing event %v with address %v", msg.Offset, walletEvent.WalletAddress)
				var eventType event.EventType
				var dataType event.DataType
				var isBackground = false
				switch walletEvent.GetEvent().(type) {
				case *kafkaevent.WalletEvent_CoinReceived:
					eventType, dataType = event.EventTypeFUNDING, event.DataTypeWalletEvent_CoinReceived
				case *kafkaevent.WalletEvent_OsmosisPoolUnlock:
					eventType, dataType = event.EventTypeDEX, event.DataTypeWalletEvent_OsmosisPoolUnlock
				case *kafkaevent.WalletEvent_Unstake:
					eventType, dataType = event.EventTypeSTAKING, event.DataTypeWalletEvent_Unstake
				case *kafkaevent.WalletEvent_NeutronTokenVesting:
					eventType, dataType = event.EventTypeFUNDING, event.DataTypeWalletEvent_NeutronTokenVesting
				case *kafkaevent.WalletEvent_Voted:
					eventType, dataType, isBackground = event.EventTypeGOVERNANCE, event.DataTypeWalletEvent_Voted, true
				case *kafkaevent.WalletEvent_VoteReminder:
					eventType, dataType = event.EventTypeGOVERNANCE, event.DataTypeWalletEvent_VoteReminder
				default:
					log.Sugar.Errorf("unknown event type %v", reflect.TypeOf(walletEvent.GetEvent()))
				}
				if el, ok := els[dataType]; ok {
					_, err := k.eventListenerManager.UpdateAddWalletEvent(context.Background(), el, &walletEvent, eventType, dataType, isBackground)
					if err != nil {
						log.Sugar.Errorf("failed to update event for %v: %v", walletEvent.WalletAddress, err)
					} else {
						if walletEvent.NotifyTime.AsTime().Before(time.Now()) {
							k.produceProcessedEvents([]*kafkaevent.EventProcessedMsg{{
								WalletAddress: walletEvent.WalletAddress,
								EventType:     eventpb.EventType(eventpb.EventType_value[eventType.String()]),
							}})
							log.Sugar.Debugf("Put event %v with address %v to `%v`", msg.Offset, walletEvent.WalletAddress, processedEventsTopic)
						}
					}
				} else {
					log.Sugar.Debugf("No event listener for %v and %v", walletEvent.WalletAddress, dataType)
				}
			} else {
				log.Sugar.Debugf("Discard event %v with address %v", msg.Offset, walletEvent.WalletAddress)
			}
			elMapMutex.Unlock()
		}
	}
}

func (k *Kafka) getEventListenerMapForChains(forChainEvents bool) map[uint64]map[event.DataType][]*ent.EventListener {
	var elMap = make(map[uint64]map[event.DataType][]*ent.EventListener)
	var dt []eventlistener.DataType
	if forChainEvents {
		dt = append(dt, eventlistener.DataTypeChainEvent_GovernanceProposal_Ongoing, eventlistener.DataTypeChainEvent_GovernanceProposal_Finished)
	} else {
		dt = append(dt, eventlistener.DataTypeContractEvent_ContractGovernanceProposal_Ongoing, eventlistener.DataTypeContractEvent_ContractGovernanceProposal_Finished)
	}
	for _, el := range k.eventListenerManager.QueryWithChain(context.Background(), dt...) {
		chainId := uint64(el.Edges.Chain.ID)
		if _, ok := elMap[chainId]; !ok {
			elMap[chainId] = make(map[event.DataType][]*ent.EventListener)
		}
		elMap[chainId][event.DataType(el.DataType)] = append(elMap[chainId][event.DataType(el.DataType)], el)
	}
	return elMap
}

func getProposalDataType(status kafkaevent.ProposalStatus) event.DataType {
	switch status.String() {
	case proposal.StatusPROPOSAL_STATUS_UNSPECIFIED.String(), proposal.StatusPROPOSAL_STATUS_DEPOSIT_PERIOD.String(), proposal.StatusPROPOSAL_STATUS_VOTING_PERIOD.String():
		return event.DataTypeChainEvent_GovernanceProposal_Ongoing
	case proposal.StatusPROPOSAL_STATUS_PASSED.String(), proposal.StatusPROPOSAL_STATUS_REJECTED.String(), proposal.StatusPROPOSAL_STATUS_FAILED.String():
		return event.DataTypeChainEvent_GovernanceProposal_Finished
	}
	log.Sugar.Panicf("Unknown proposal status: %v", status.String())
	return event.DataTypeChainEvent_GovernanceProposal_Ongoing
}

func getContractProposalDataType(status kafkaevent.ContractProposalStatus) event.DataType {
	switch status.String() {
	case contractproposal.StatusOPEN.String():
		return event.DataTypeChainEvent_GovernanceProposal_Ongoing
	case contractproposal.StatusREJECTED.String(), contractproposal.StatusPASSED.String(), contractproposal.StatusEXECUTED.String(), contractproposal.StatusCLOSED.String(), contractproposal.StatusEXECUTION_FAILED.String():
		return event.DataTypeChainEvent_GovernanceProposal_Finished
	}
	log.Sugar.Panicf("Unknown proposal status: %v", status.String())
	return event.DataTypeChainEvent_GovernanceProposal_Ongoing
}

func (k *Kafka) getChains() map[uint64]*ent.Chain {
	var chains = make(map[uint64]*ent.Chain)
	for _, chain := range k.chainManager.QueryAll(context.Background()) {
		chains[uint64(chain.ID)] = chain
	}
	return chains
}

func (k *Kafka) ProcessChainEvents() {
	r := k.groupReader(chainEventsTopic)
	defer k.closeReader(r)

	elMap := k.getEventListenerMapForChains(true)
	chains := k.getChains()     // TODO: update chains when new chain is added
	elMapMutex := sync.Mutex{}  // Mutex to synchronize access to elMap
	chainsMutex := sync.Mutex{} // Mutex to synchronize access to chains

	go func() {
		kafkaInternal := kafka_internal.NewKafkaInternal(k.addresses)

		ch := make(chan kafka_internal.DbChange)
		defer close(ch)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var subscribe = []kafka_internal.DbChange{kafka_internal.EventListenerCreated, kafka_internal.EventListenerDeleted, kafka_internal.ChainEnabled, kafka_internal.ChainDisabled}
		go kafkaInternal.ReadDbChanges(ctx, ch, subscribe)

		for change := range ch {
			switch change {
			case kafka_internal.EventListenerCreated, kafka_internal.EventListenerDeleted:
				elMapMutex.Lock()
				elMap = k.getEventListenerMapForChains(true)
				elMapMutex.Unlock()
			case kafka_internal.ChainEnabled, kafka_internal.ChainDisabled:
				chainsMutex.Lock()
				chains = k.getChains()
				chainsMutex.Unlock()
			}
		}
		log.Sugar.Debugf("Stopped processing chain events")
	}()

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Sugar.Errorf("failed to read message: %v", err)
		}
		var chainEvent kafkaevent.ChainEvent
		err = proto.Unmarshal(msg.Value, &chainEvent)
		if err != nil {
			log.Sugar.Error(err)
		} else {
			elMapMutex.Lock()
			chainsMutex.Lock()
			chain, ok := chains[chainEvent.ChainId]
			if !ok {
				log.Sugar.Panicf("failed to find chain %v", chainEvent.ChainId)
			}
			log.Sugar.Debugf("Processing event %v for chain %v", msg.Offset, chain.PrettyName)
			var pbEvents []*kafkaevent.EventProcessedMsg
			var eventType event.EventType
			var dataType event.DataType
			switch chainEvent.GetEvent().(type) {
			case *kafkaevent.ChainEvent_GovernanceProposal:
				var ignoredStates = []string{
					proposal.StatusPROPOSAL_STATUS_UNSPECIFIED.String(),
					proposal.StatusPROPOSAL_STATUS_DEPOSIT_PERIOD.String(),
				}
				if slices.Contains(ignoredStates, chainEvent.GetGovernanceProposal().GetProposalStatus().String()) {
					continue
				}
				eventType = event.EventTypeGOVERNANCE
				dataType = getProposalDataType(chainEvent.GetGovernanceProposal().GetProposalStatus())
			case *kafkaevent.ChainEvent_ValidatorOutOfActiveSet:
				eventType, dataType = event.EventTypeSTAKING, event.DataTypeChainEvent_ValidatorOutOfActiveSet
			case *kafkaevent.ChainEvent_ValidatorSlash:
				eventType, dataType = event.EventTypeSTAKING, event.DataTypeChainEvent_ValidatorSlash
			default:
				log.Sugar.Errorf("Unknown event type: %v", reflect.TypeOf(chainEvent.GetEvent()))
			}
			if dTypes, ok := elMap[chainEvent.ChainId]; ok {
				if els, ok := dTypes[dataType]; ok {
					for _, el := range els {
						_, err := k.eventListenerManager.UpdateAddChainEvent(context.Background(), el, &chainEvent, eventType, dataType)
						if err != nil {
							log.Sugar.Errorf("failed to update event for %v: %v", el.WalletAddress, err)
						} else {
							if chainEvent.NotifyTime.AsTime().Before(time.Now()) {
								pbEvents = append(pbEvents, &kafkaevent.EventProcessedMsg{
									ChainId:       chainEvent.ChainId,
									WalletAddress: "",
									EventType:     eventpb.EventType(eventpb.EventType_value[eventType.String()]),
								})
								log.Sugar.Debugf("Put event %v with address %v to `%v`", msg.Offset, el.WalletAddress, processedEventsTopic)
							}
						}
					}
				}
			}
			if len(pbEvents) > 0 {
				log.Sugar.Debugf("Produce %v events", len(pbEvents))
				k.produceProcessedEvents(pbEvents)
			}
			elMapMutex.Unlock()
			chainsMutex.Unlock()
		}
	}
}

func (k *Kafka) ProcessContractEvents() {
	r := k.groupReader(contractEventsTopic)
	defer k.closeReader(r)

	elMap := k.getEventListenerMapForChains(false)
	chains := k.getChains()
	elMapMutex := sync.Mutex{} // Mutex to synchronize access to elMap
	chainsMutex := sync.Mutex{}

	go func() {
		kafkaInternal := kafka_internal.NewKafkaInternal(k.addresses)

		ch := make(chan kafka_internal.DbChange)
		defer close(ch)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var subscribe = []kafka_internal.DbChange{kafka_internal.EventListenerCreated, kafka_internal.EventListenerDeleted, kafka_internal.ChainEnabled, kafka_internal.ChainDisabled}
		go kafkaInternal.ReadDbChanges(ctx, ch, subscribe)

		for change := range ch {
			switch change {
			case kafka_internal.EventListenerCreated, kafka_internal.EventListenerDeleted:
				elMapMutex.Lock()
				elMap = k.getEventListenerMapForChains(false)
				elMapMutex.Unlock()
			case kafka_internal.ChainEnabled, kafka_internal.ChainDisabled:
				chainsMutex.Lock()
				chains = k.getChains()
				chainsMutex.Unlock()
			}
		}
	}()

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Sugar.Errorf("failed to read message: %v", err)
		}
		var contractEvent kafkaevent.ContractEvent
		err = proto.Unmarshal(msg.Value, &contractEvent)
		if err != nil {
			log.Sugar.Error(err)
		} else {
			elMapMutex.Lock()
			chainsMutex.Lock()
			chain, ok := chains[contractEvent.ChainId]
			if !ok {
				log.Sugar.Panicf("failed to find chain %v", contractEvent.ChainId)
			}
			log.Sugar.Debugf("Processing event %v for chain %v", msg.Offset, chain.PrettyName)
			var pbEvents []*kafkaevent.EventProcessedMsg
			var eventType event.EventType
			var dataType event.DataType
			switch contractEvent.GetEvent().(type) {
			case *kafkaevent.ContractEvent_ContractGovernanceProposal:
				dataType = getContractProposalDataType(contractEvent.GetContractGovernanceProposal().GetProposalStatus())
				eventType = event.EventTypeGOVERNANCE
			default:
				log.Sugar.Errorf("Unknown event type: %v", reflect.TypeOf(contractEvent.GetEvent()))
			}
			if dTypes, ok := elMap[contractEvent.ChainId]; ok {
				if els, ok := dTypes[dataType]; ok {
					for _, el := range els {
						_, err := k.eventListenerManager.UpdateAddContractEvent(context.Background(), el, &contractEvent, eventType, dataType)
						if err != nil {
							log.Sugar.Errorf("failed to update event for %v: %v", el.WalletAddress, err)
						} else {
							if contractEvent.NotifyTime.AsTime().Before(time.Now()) {
								pbEvents = append(pbEvents, &kafkaevent.EventProcessedMsg{
									ChainId:       contractEvent.ChainId,
									WalletAddress: "",
									EventType:     eventpb.EventType(eventpb.EventType_value[eventType.String()]),
								})
								log.Sugar.Debugf("Put event %v with address %v to `%v`", msg.Offset, el.WalletAddress, processedEventsTopic)
							}
						}
					}
				}
			}
			if len(pbEvents) > 0 {
				log.Sugar.Debugf("Produce %v events", len(pbEvents))
				k.produceProcessedEvents(pbEvents)
			}
			elMapMutex.Unlock()
			chainsMutex.Unlock()
		}
	}
}

func (k *Kafka) ConsumeProcessedEvents(ctx context.Context, user *ent.User, eventsChannel chan *eventpb.NewEvent) {
	log.Sugar.Debugf("Start processed-events consumer for user %v", user.ID)
	els := k.eventListenerManager.QueryByUser(ctx, user)
	subscriptions := map[string]interface{}{}
	for _, el := range els {
		if el.WalletAddress != "" {
			subscriptions[el.WalletAddress] = nil
		}
		subscriptions[strconv.Itoa(el.Edges.Chain.ID)] = nil
	}

	r := k.reader(processedEventsTopic)
	defer k.closeReader(r)

	err := r.SetOffsetAt(context.Background(), time.Now())
	if err != nil {
		log.Sugar.Errorf("failed to set offset: %v", err)
		eventsChannel <- nil
	}

	for {
		select {
		case <-ctx.Done():
			log.Sugar.Debugf("Stop the processed-events consumer for user %v", user.ID)
			return
		default:
			msg, err := r.ReadMessage(context.Background())
			if err != nil {
				break
			}
			var processedEvent kafkaevent.EventProcessedMsg
			err = proto.Unmarshal(msg.Value, &processedEvent)
			if err != nil {
				log.Sugar.Error(err)
				break
			}
			if _, ok := subscriptions[processedEvent.GetWalletAddress()]; ok {
				eventType := processedEvent.GetEventType()
				eventsChannel <- &eventpb.NewEvent{EventType: &eventType}
			} else if _, ok := subscriptions[strconv.Itoa(int(processedEvent.GetChainId()))]; ok {
				eventType := processedEvent.GetEventType()
				eventsChannel <- &eventpb.NewEvent{EventType: &eventType}
			}
		}
	}
}
