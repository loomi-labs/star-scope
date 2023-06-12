package kafka_internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/loomi-labs/star-scope/types"
	"github.com/segmentio/kafka-go"
	"github.com/shifty11/go-logger/log"
	"golang.org/x/exp/slices"
	"time"
)

type Topic string

const (
	DbEntityChanged     = Topic(types.DbEntityChanged)
	WalletEventsTopic   = Topic(types.WalletEventsTopic)
	ChainEventsTopic    = Topic(types.ChainEventsTopic)
	ContractEventsTopic = Topic(types.ContractEventsTopic)
)

type DbChange string

const (
	EventListenerCreated DbChange = "event-listener-created"
	EventListenerDeleted DbChange = "event-listener-deleted"
	ChainEnabled         DbChange = "chain-enabled"
	ChainDisabled        DbChange = "chain-disabled"
)

func toDbChange(data []byte) (DbChange, error) {
	strValue := string(data)
	switch strValue {
	case string(EventListenerCreated):
		return EventListenerCreated, nil
	case string(EventListenerDeleted):
		return EventListenerDeleted, nil
	case string(ChainEnabled):
		return ChainEnabled, nil
	case string(ChainDisabled):
		return ChainDisabled, nil
	default:
		return "", errors.New("invalid DbChange value")
	}
}

type KafkaInternal interface {
	ProduceDbChangeMsg(dbChange DbChange)
	ReadDbChanges(ctx context.Context, ch chan DbChange, subscribedChanges []DbChange)
	ProduceWalletEvents(msgs [][]byte)
	ProduceChainEvents(msgs [][]byte)
	ProduceContractEvents(msgs [][]byte)
}

type kafkaInternal struct {
	addresses []string
}

func NewKafkaInternal(addresses []string) KafkaInternal {
	return &kafkaInternal{
		addresses: addresses,
	}
}

func (k *kafkaInternal) reader(topic Topic) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   k.addresses,
		Topic:     string(topic),
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		MaxWait:   1 * time.Second,
	})
}

func (k *kafkaInternal) closeReader(r *kafka.Reader) {
	err := r.Close()
	if err != nil {
		log.Sugar.Fatal("failed to closeReader writer:", err)
	}
}

func (k *kafkaInternal) writer(topic Topic) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(k.addresses...),
		Topic:    string(topic),
		Balancer: &kafka.LeastBytes{},
	}
}

func (k *kafkaInternal) closeWriter(w *kafka.Writer) {
	err := w.Close()
	if err != nil {
		log.Sugar.Fatal("failed to close writer:", err)
	}
}

func (k *kafkaInternal) ProduceDbChangeMsg(dbChange DbChange) {
	w := k.writer(DbEntityChanged)
	defer k.closeWriter(w)

	err := w.WriteMessages(context.Background(), kafka.Message{Value: []byte(dbChange)})
	if err != nil {
		log.Sugar.Panicf(fmt.Sprintf("failed to write messages: %v", err))
	}
}

func (k *kafkaInternal) ReadDbChanges(ctx context.Context, ch chan DbChange, subscribedChanges []DbChange) {
	r := k.reader(DbEntityChanged)
	defer k.closeReader(r)
	err := r.SetOffsetAt(context.Background(), time.Now())
	if err != nil {
		log.Sugar.Panicf("failed to set offset: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Sugar.Errorf("failed to read message: %v", err)
			}
			dbChange, err := toDbChange(msg.Value)
			if err != nil {
				log.Sugar.Error(err)
			}
			if slices.Contains(subscribedChanges, dbChange) {
				ch <- dbChange
			}
		}
	}
}

func (k *kafkaInternal) produceEvents(topic Topic, msgs [][]byte) {
	w := k.writer(topic)
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

func (k *kafkaInternal) ProduceWalletEvents(msgs [][]byte) {
	k.produceEvents(WalletEventsTopic, msgs)
}

func (k *kafkaInternal) ProduceChainEvents(msgs [][]byte) {
	k.produceEvents(ChainEventsTopic, msgs)
}

func (k *kafkaInternal) ProduceContractEvents(msgs [][]byte) {
	k.produceEvents(ContractEventsTopic, msgs)
}
