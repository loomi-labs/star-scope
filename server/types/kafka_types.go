package types

type Topic string

const (
	ChainEventsTopic     Topic = "chain-events"
	ContractEventsTopic  Topic = "contract-events"
	WalletEventsTopic    Topic = "wallet-events"
	ProcessedEventsTopic Topic = "processed-events"
	DbEntityChanged      Topic = "db-entity-changed"
)
