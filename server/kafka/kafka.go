package kafka

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/event"
	"github.com/loomi-labs/star-scope/grpc/event/eventpb"
	"github.com/loomi-labs/star-scope/indexevent"
	"github.com/segmentio/kafka-go"
	"github.com/shifty11/go-logger/log"
	"time"
)

var (
	indexEventsTopic     = "index-events"
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

func (k *Kafka) getEventListenerMap() map[string]*ent.EventListener {
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

	for {
		// TODO: find a better way to get event listener map than querying db every time
		elMap := k.getEventListenerMap()
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Sugar.Errorf("failed to read message: %v", err)
		}
		var txEvent indexevent.TxEvent
		err = proto.Unmarshal(msg.Value, &txEvent)
		if err != nil {
			log.Sugar.Error(err)
		} else {
			if el, ok := elMap[txEvent.WalletAddress]; ok {
				var ctx = context.Background()
				var err error
				switch txEvent.GetEvent().(type) {
				case *indexevent.TxEvent_CoinReceived:
					_, err = k.eventListenerManager.UpdateAddEvent(ctx, el, event.TypeTxEvent_CoinReceived, txEvent.NotifyTime.AsTime(), msg.Value)
					log.Sugar.Infof("%v received %v%v from %v", txEvent.WalletAddress, txEvent.GetCoinReceived().GetCoin().Amount, txEvent.GetCoinReceived().GetCoin().Denom, txEvent.GetCoinReceived().Sender)
				case *indexevent.TxEvent_OsmosisPoolUnlock:
					_, err = k.eventListenerManager.UpdateAddEvent(ctx, el, event.TypeTxEvent_OsmosisPoolUnlock, txEvent.NotifyTime.AsTime(), msg.Value)
					log.Sugar.Infof("%v will unlock pool at %v", txEvent.WalletAddress, txEvent.GetOsmosisPoolUnlock().UnlockTime)
				case *indexevent.TxEvent_Unstake:
					_, err = k.eventListenerManager.UpdateAddEvent(ctx, el, event.TypeTxEvent_Unstake, txEvent.NotifyTime.AsTime(), msg.Value)
					log.Sugar.Infof("%v will unstake %v%v at %v", txEvent.WalletAddress, txEvent.GetUnstake().GetCoin().Amount, txEvent.GetUnstake().GetCoin().Denom, txEvent.GetUnstake().CompletionTime)
				}
				if err != nil {
					log.Sugar.Errorf("failed to update event for %v: %v", txEvent.WalletAddress, err)
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
