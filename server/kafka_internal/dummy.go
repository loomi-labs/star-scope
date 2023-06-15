package kafka_internal

import "context"

type kafkaInternalDummy struct {
}

func NewKafkaInternalDummy() KafkaInternal {
	return &kafkaInternalDummy{}
}

func (k kafkaInternalDummy) ProduceDbChangeMsg(_ DbChange) {
}

func (k kafkaInternalDummy) ReadDbChanges(_ context.Context, _ chan DbChange, _ []DbChange) {
}

func (k kafkaInternalDummy) ProduceWalletEvents(_ [][]byte) {
}

func (k kafkaInternalDummy) ProduceChainEvents(_ [][]byte) {
}

func (k kafkaInternalDummy) ProduceContractEvents(_ [][]byte) {
}
