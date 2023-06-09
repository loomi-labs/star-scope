package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/shifty11/go-logger/log"
)

type Topic string

var (
	WalletEvents Topic = "wallet-events"
)

type KafkaProducer struct {
	topic     Topic
	addresses []string
}

func NewKafkaProducer(topic Topic, addresses ...string) *KafkaProducer {
	return &KafkaProducer{
		topic:     topic,
		addresses: addresses,
	}
}

func (k *KafkaProducer) writer() *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(k.addresses...),
		Topic:                  string(k.topic),
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
