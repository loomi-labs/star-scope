// Code generated by ent, DO NOT EDIT.

package event

import (
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the event type in the database.
	Label = "event"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time field in the database.
	FieldUpdateTime = "update_time"
	// FieldEventType holds the string denoting the event_type field in the database.
	FieldEventType = "event_type"
	// FieldChainEvent holds the string denoting the chain_event field in the database.
	FieldChainEvent = "chain_event"
	// FieldContractEvent holds the string denoting the contract_event field in the database.
	FieldContractEvent = "contract_event"
	// FieldWalletEvent holds the string denoting the wallet_event field in the database.
	FieldWalletEvent = "wallet_event"
	// FieldDataType holds the string denoting the data_type field in the database.
	FieldDataType = "data_type"
	// FieldNotifyTime holds the string denoting the notify_time field in the database.
	FieldNotifyTime = "notify_time"
	// FieldIsRead holds the string denoting the is_read field in the database.
	FieldIsRead = "is_read"
	// FieldIsBackground holds the string denoting the is_background field in the database.
	FieldIsBackground = "is_background"
	// EdgeEventListener holds the string denoting the event_listener edge name in mutations.
	EdgeEventListener = "event_listener"
	// Table holds the table name of the event in the database.
	Table = "events"
	// EventListenerTable is the table that holds the event_listener relation/edge.
	EventListenerTable = "events"
	// EventListenerInverseTable is the table name for the EventListener entity.
	// It exists in this package in order to avoid circular dependency with the "eventlistener" package.
	EventListenerInverseTable = "event_listeners"
	// EventListenerColumn is the table column denoting the event_listener relation/edge.
	EventListenerColumn = "event_listener_events"
)

// Columns holds all SQL columns for event fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldEventType,
	FieldChainEvent,
	FieldContractEvent,
	FieldWalletEvent,
	FieldDataType,
	FieldNotifyTime,
	FieldIsRead,
	FieldIsBackground,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "events"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"event_listener_events",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreateTime holds the default value on creation for the "create_time" field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the "update_time" field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the "update_time" field.
	UpdateDefaultUpdateTime func() time.Time
	// DefaultNotifyTime holds the default value on creation for the "notify_time" field.
	DefaultNotifyTime time.Time
	// DefaultIsRead holds the default value on creation for the "is_read" field.
	DefaultIsRead bool
	// DefaultIsBackground holds the default value on creation for the "is_background" field.
	DefaultIsBackground bool
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// EventType defines the type for the "event_type" enum field.
type EventType string

// EventType values.
const (
	EventTypeGOVERNANCE EventType = "GOVERNANCE"
	EventTypeFUNDING    EventType = "FUNDING"
	EventTypeSTAKING    EventType = "STAKING"
	EventTypeDEX        EventType = "DEX"
)

func (et EventType) String() string {
	return string(et)
}

// EventTypeValidator is a validator for the "event_type" field enum values. It is called by the builders before save.
func EventTypeValidator(et EventType) error {
	switch et {
	case EventTypeGOVERNANCE, EventTypeFUNDING, EventTypeSTAKING, EventTypeDEX:
		return nil
	default:
		return fmt.Errorf("event: invalid enum value for event_type field: %q", et)
	}
}

// DataType defines the type for the "data_type" enum field.
type DataType string

// DataType values.
const (
	DataTypeWalletEvent_CoinReceived                          DataType = "WalletEvent_CoinReceived"
	DataTypeWalletEvent_OsmosisPoolUnlock                     DataType = "WalletEvent_OsmosisPoolUnlock"
	DataTypeWalletEvent_Unstake                               DataType = "WalletEvent_Unstake"
	DataTypeWalletEvent_NeutronTokenVesting                   DataType = "WalletEvent_NeutronTokenVesting"
	DataTypeWalletEvent_Voted                                 DataType = "WalletEvent_Voted"
	DataTypeWalletEvent_VoteReminder                          DataType = "WalletEvent_VoteReminder"
	DataTypeChainEvent_ValidatorOutOfActiveSet                DataType = "ChainEvent_ValidatorOutOfActiveSet"
	DataTypeChainEvent_ValidatorSlash                         DataType = "ChainEvent_ValidatorSlash"
	DataTypeChainEvent_GovernanceProposal_Ongoing             DataType = "ChainEvent_GovernanceProposal_Ongoing"
	DataTypeChainEvent_GovernanceProposal_Finished            DataType = "ChainEvent_GovernanceProposal_Finished"
	DataTypeContractEvent_ContractGovernanceProposal_Ongoing  DataType = "ContractEvent_ContractGovernanceProposal_Ongoing"
	DataTypeContractEvent_ContractGovernanceProposal_Finished DataType = "ContractEvent_ContractGovernanceProposal_Finished"
)

func (dt DataType) String() string {
	return string(dt)
}

// DataTypeValidator is a validator for the "data_type" field enum values. It is called by the builders before save.
func DataTypeValidator(dt DataType) error {
	switch dt {
	case DataTypeWalletEvent_CoinReceived, DataTypeWalletEvent_OsmosisPoolUnlock, DataTypeWalletEvent_Unstake, DataTypeWalletEvent_NeutronTokenVesting, DataTypeWalletEvent_Voted, DataTypeWalletEvent_VoteReminder, DataTypeChainEvent_ValidatorOutOfActiveSet, DataTypeChainEvent_ValidatorSlash, DataTypeChainEvent_GovernanceProposal_Ongoing, DataTypeChainEvent_GovernanceProposal_Finished, DataTypeContractEvent_ContractGovernanceProposal_Ongoing, DataTypeContractEvent_ContractGovernanceProposal_Finished:
		return nil
	default:
		return fmt.Errorf("event: invalid enum value for data_type field: %q", dt)
	}
}

// OrderOption defines the ordering options for the Event queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCreateTime orders the results by the create_time field.
func ByCreateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreateTime, opts...).ToFunc()
}

// ByUpdateTime orders the results by the update_time field.
func ByUpdateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdateTime, opts...).ToFunc()
}

// ByEventType orders the results by the event_type field.
func ByEventType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEventType, opts...).ToFunc()
}

// ByDataType orders the results by the data_type field.
func ByDataType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDataType, opts...).ToFunc()
}

// ByNotifyTime orders the results by the notify_time field.
func ByNotifyTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldNotifyTime, opts...).ToFunc()
}

// ByIsRead orders the results by the is_read field.
func ByIsRead(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsRead, opts...).ToFunc()
}

// ByIsBackground orders the results by the is_background field.
func ByIsBackground(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsBackground, opts...).ToFunc()
}

// ByEventListenerField orders the results by event_listener field.
func ByEventListenerField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newEventListenerStep(), sql.OrderByField(field, opts...))
	}
}
func newEventListenerStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(EventListenerInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, EventListenerTable, EventListenerColumn),
	)
}
