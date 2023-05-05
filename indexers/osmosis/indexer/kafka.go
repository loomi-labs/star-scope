package indexer

import (
	"context"
	"fmt"
	"github.com/loomi-labs/star-scope/indexers/osmosis/common"
	"github.com/segmentio/kafka-go"
	"github.com/shifty11/go-logger/log"
)

var (
	topic = common.GetEnvX("KAFKA_TOPIC")
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
		Addr:                   kafka.TCP(k.addresses...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
		//Logger:                 kafka.LoggerFunc(log.Sugar.Debugf),
		//ErrorLogger:            kafka.LoggerFunc(log.Sugar.Errorf),
	}
}

func (k *KafkaProducer) close(w *kafka.Writer) {
	err := w.Close()
	if err != nil {
		log.Sugar.Fatal("failed to close writer:", err)
	}
}

func (k *KafkaProducer) Produce(msgs [][]byte) {
	w := k.writer()
	defer k.close(w)

	kafkaMsgs := make([]kafka.Message, len(msgs))
	for i, msg := range msgs {
		kafkaMsgs[i] = kafka.Message{Value: msg}
	}

	err := w.WriteMessages(context.Background(), kafkaMsgs...)
	if err != nil {
		panic(fmt.Sprintf("failed to write messages: %v", err))
	}
}
