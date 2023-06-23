// Code generated by ent, DO NOT EDIT.

package contractproposal

import (
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the contractproposal type in the database.
	Label = "contract_proposal"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time field in the database.
	FieldUpdateTime = "update_time"
	// FieldProposalID holds the string denoting the proposal_id field in the database.
	FieldProposalID = "proposal_id"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldFirstSeenTime holds the string denoting the first_seen_time field in the database.
	FieldFirstSeenTime = "first_seen_time"
	// FieldVotingEndTime holds the string denoting the voting_end_time field in the database.
	FieldVotingEndTime = "voting_end_time"
	// FieldContractAddress holds the string denoting the contract_address field in the database.
	FieldContractAddress = "contract_address"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// EdgeChain holds the string denoting the chain edge name in mutations.
	EdgeChain = "chain"
	// Table holds the table name of the contractproposal in the database.
	Table = "contract_proposals"
	// ChainTable is the table that holds the chain relation/edge.
	ChainTable = "contract_proposals"
	// ChainInverseTable is the table name for the Chain entity.
	// It exists in this package in order to avoid circular dependency with the "chain" package.
	ChainInverseTable = "chains"
	// ChainColumn is the table column denoting the chain relation/edge.
	ChainColumn = "chain_contract_proposals"
)

// Columns holds all SQL columns for contractproposal fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldProposalID,
	FieldTitle,
	FieldDescription,
	FieldFirstSeenTime,
	FieldVotingEndTime,
	FieldContractAddress,
	FieldStatus,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "contract_proposals"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"chain_contract_proposals",
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
)

// Status defines the type for the "status" enum field.
type Status string

// Status values.
const (
	StatusOPEN             Status = "OPEN"
	StatusREJECTED         Status = "REJECTED"
	StatusPASSED           Status = "PASSED"
	StatusEXECUTED         Status = "EXECUTED"
	StatusCLOSED           Status = "CLOSED"
	StatusEXECUTION_FAILED Status = "EXECUTION_FAILED"
)

func (s Status) String() string {
	return string(s)
}

// StatusValidator is a validator for the "status" field enum values. It is called by the builders before save.
func StatusValidator(s Status) error {
	switch s {
	case StatusOPEN, StatusREJECTED, StatusPASSED, StatusEXECUTED, StatusCLOSED, StatusEXECUTION_FAILED:
		return nil
	default:
		return fmt.Errorf("contractproposal: invalid enum value for status field: %q", s)
	}
}

// OrderOption defines the ordering options for the ContractProposal queries.
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

// ByProposalID orders the results by the proposal_id field.
func ByProposalID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldProposalID, opts...).ToFunc()
}

// ByTitle orders the results by the title field.
func ByTitle(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTitle, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByFirstSeenTime orders the results by the first_seen_time field.
func ByFirstSeenTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFirstSeenTime, opts...).ToFunc()
}

// ByVotingEndTime orders the results by the voting_end_time field.
func ByVotingEndTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVotingEndTime, opts...).ToFunc()
}

// ByContractAddress orders the results by the contract_address field.
func ByContractAddress(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldContractAddress, opts...).ToFunc()
}

// ByStatus orders the results by the status field.
func ByStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatus, opts...).ToFunc()
}

// ByChainField orders the results by chain field.
func ByChainField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newChainStep(), sql.OrderByField(field, opts...))
	}
}
func newChainStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ChainInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, ChainTable, ChainColumn),
	)
}
