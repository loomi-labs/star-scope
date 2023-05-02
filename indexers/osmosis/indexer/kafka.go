package indexer

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/kafka-go"
	indexEvent "github.com/shifty11/blocklog-backend/indexers/osmosis/index_event"
	"github.com/shifty11/go-logger/log"
)

var (
	topic = "new-event"
)

type KafkaProducer struct {
	addresses []string
}

func NewKafkaProducer(addresses ...string) *KafkaProducer {
	return &KafkaProducer{
		addresses: addresses,
	}
}

func (k *KafkaProducer) writer() *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(k.addresses...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func (k *KafkaProducer) close(w *kafka.Writer) {
	err := w.Close()
	if err != nil {
		log.Sugar.Fatal("failed to close writer:", err)
	}
}

func (k *KafkaProducer) Produce(txEvent *indexEvent.TxEvent) {
	msg, err := proto.Marshal(txEvent)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal txEvent: %v", err))
	}
	w := k.writer()
	defer k.close(w)

	err = w.WriteMessages(context.Background(), kafka.Message{
		Value: msg,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to write messages: %v", err))
	}
}
