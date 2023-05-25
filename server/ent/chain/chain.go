// Code generated by ent, DO NOT EDIT.

package chain

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the chain type in the database.
	Label = "chain"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time field in the database.
	FieldUpdateTime = "update_time"
	// FieldChainID holds the string denoting the chain_id field in the database.
	FieldChainID = "chain_id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldPrettyName holds the string denoting the pretty_name field in the database.
	FieldPrettyName = "pretty_name"
	// FieldPath holds the string denoting the path field in the database.
	FieldPath = "path"
	// FieldImage holds the string denoting the image field in the database.
	FieldImage = "image"
	// FieldBech32Prefix holds the string denoting the bech32_prefix field in the database.
	FieldBech32Prefix = "bech32_prefix"
	// FieldIndexingHeight holds the string denoting the indexing_height field in the database.
	FieldIndexingHeight = "indexing_height"
	// FieldHasCustomIndexer holds the string denoting the has_custom_indexer field in the database.
	FieldHasCustomIndexer = "has_custom_indexer"
	// FieldHandledMessageTypes holds the string denoting the handled_message_types field in the database.
	FieldHandledMessageTypes = "handled_message_types"
	// FieldUnhandledMessageTypes holds the string denoting the unhandled_message_types field in the database.
	FieldUnhandledMessageTypes = "unhandled_message_types"
	// FieldIsEnabled holds the string denoting the is_enabled field in the database.
	FieldIsEnabled = "is_enabled"
	// EdgeEventListeners holds the string denoting the event_listeners edge name in mutations.
	EdgeEventListeners = "event_listeners"
	// EdgeProposals holds the string denoting the proposals edge name in mutations.
	EdgeProposals = "proposals"
	// EdgeContractProposals holds the string denoting the contract_proposals edge name in mutations.
	EdgeContractProposals = "contract_proposals"
	// Table holds the table name of the chain in the database.
	Table = "chains"
	// EventListenersTable is the table that holds the event_listeners relation/edge.
	EventListenersTable = "event_listeners"
	// EventListenersInverseTable is the table name for the EventListener entity.
	// It exists in this package in order to avoid circular dependency with the "eventlistener" package.
	EventListenersInverseTable = "event_listeners"
	// EventListenersColumn is the table column denoting the event_listeners relation/edge.
	EventListenersColumn = "chain_event_listeners"
	// ProposalsTable is the table that holds the proposals relation/edge.
	ProposalsTable = "proposals"
	// ProposalsInverseTable is the table name for the Proposal entity.
	// It exists in this package in order to avoid circular dependency with the "proposal" package.
	ProposalsInverseTable = "proposals"
	// ProposalsColumn is the table column denoting the proposals relation/edge.
	ProposalsColumn = "chain_proposals"
	// ContractProposalsTable is the table that holds the contract_proposals relation/edge.
	ContractProposalsTable = "contract_proposals"
	// ContractProposalsInverseTable is the table name for the ContractProposal entity.
	// It exists in this package in order to avoid circular dependency with the "contractproposal" package.
	ContractProposalsInverseTable = "contract_proposals"
	// ContractProposalsColumn is the table column denoting the contract_proposals relation/edge.
	ContractProposalsColumn = "chain_contract_proposals"
)

// Columns holds all SQL columns for chain fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldChainID,
	FieldName,
	FieldPrettyName,
	FieldPath,
	FieldImage,
	FieldBech32Prefix,
	FieldIndexingHeight,
	FieldHasCustomIndexer,
	FieldHandledMessageTypes,
	FieldUnhandledMessageTypes,
	FieldIsEnabled,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
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
	// DefaultIndexingHeight holds the default value on creation for the "indexing_height" field.
	DefaultIndexingHeight uint64
	// DefaultHasCustomIndexer holds the default value on creation for the "has_custom_indexer" field.
	DefaultHasCustomIndexer bool
	// DefaultHandledMessageTypes holds the default value on creation for the "handled_message_types" field.
	DefaultHandledMessageTypes string
	// DefaultUnhandledMessageTypes holds the default value on creation for the "unhandled_message_types" field.
	DefaultUnhandledMessageTypes string
	// DefaultIsEnabled holds the default value on creation for the "is_enabled" field.
	DefaultIsEnabled bool
)

// OrderOption defines the ordering options for the Chain queries.
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

// ByChainID orders the results by the chain_id field.
func ByChainID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldChainID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByPrettyName orders the results by the pretty_name field.
func ByPrettyName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPrettyName, opts...).ToFunc()
}

// ByPath orders the results by the path field.
func ByPath(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPath, opts...).ToFunc()
}

// ByImage orders the results by the image field.
func ByImage(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldImage, opts...).ToFunc()
}

// ByBech32Prefix orders the results by the bech32_prefix field.
func ByBech32Prefix(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBech32Prefix, opts...).ToFunc()
}

// ByIndexingHeight orders the results by the indexing_height field.
func ByIndexingHeight(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIndexingHeight, opts...).ToFunc()
}

// ByHasCustomIndexer orders the results by the has_custom_indexer field.
func ByHasCustomIndexer(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHasCustomIndexer, opts...).ToFunc()
}

// ByHandledMessageTypes orders the results by the handled_message_types field.
func ByHandledMessageTypes(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHandledMessageTypes, opts...).ToFunc()
}

// ByUnhandledMessageTypes orders the results by the unhandled_message_types field.
func ByUnhandledMessageTypes(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUnhandledMessageTypes, opts...).ToFunc()
}

// ByIsEnabled orders the results by the is_enabled field.
func ByIsEnabled(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsEnabled, opts...).ToFunc()
}

// ByEventListenersCount orders the results by event_listeners count.
func ByEventListenersCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newEventListenersStep(), opts...)
	}
}

// ByEventListeners orders the results by event_listeners terms.
func ByEventListeners(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newEventListenersStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByProposalsCount orders the results by proposals count.
func ByProposalsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newProposalsStep(), opts...)
	}
}

// ByProposals orders the results by proposals terms.
func ByProposals(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newProposalsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByContractProposalsCount orders the results by contract_proposals count.
func ByContractProposalsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newContractProposalsStep(), opts...)
	}
}

// ByContractProposals orders the results by contract_proposals terms.
func ByContractProposals(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newContractProposalsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newEventListenersStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(EventListenersInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, EventListenersTable, EventListenersColumn),
	)
}
func newProposalsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ProposalsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ProposalsTable, ProposalsColumn),
	)
}
func newContractProposalsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ContractProposalsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ContractProposalsTable, ContractProposalsColumn),
	)
}
