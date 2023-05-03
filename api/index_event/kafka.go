package index_event

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/kafka-go"
	"github.com/shifty11/blocklog-backend/database"
	"github.com/shifty11/go-logger/log"
	"time"
)

var (
	topic = "index-events"
)

type KafkaConsumer struct {
	addresses            []string
	eventListenerManager *database.EventListenerManager
}

func NewKafkaConsumer(dbManager *database.DbManagers, addresses ...string) *KafkaConsumer {
	return &KafkaConsumer{
		addresses:            addresses,
		eventListenerManager: dbManager.EventListenerManager,
	}
}

func (k *KafkaConsumer) reader() *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   k.addresses,
		Topic:     topic,
		GroupID:   topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		MaxWait:   1 * time.Second,
	})
	return r
}

func (k *KafkaConsumer) close(r *kafka.Reader) {
	err := r.Close()
	if err != nil {
		log.Sugar.Fatal("failed to close writer:", err)
	}
}

func (k *KafkaConsumer) StartConsuming() {
	r := k.reader()
	defer k.close(r)

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		var txEvent TxEvent
		err = proto.Unmarshal(msg.Value, &txEvent)
		if err != nil {
			log.Sugar.Error(err)
		} else {
			switch txEvent.GetEvent().(type) {
			case *TxEvent_CoinReceived:
				log.Sugar.Infof("%v received %v%v from %v", txEvent.WalletAddress, txEvent.GetCoinReceived().GetCoin().Amount, txEvent.GetCoinReceived().GetCoin().Denom, txEvent.GetCoinReceived().Sender)
			case *TxEvent_OsmosisPoolUnlock:
				log.Sugar.Infof("%v will unlock pool at %v", txEvent.WalletAddress, txEvent.GetOsmosisPoolUnlock().UnlockTime)
			}
		}
	}
}
