// Code generated by ent, DO NOT EDIT.

package proposal

import (
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the proposal type in the database.
	Label = "proposal"
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
	// FieldVotingStartTime holds the string denoting the voting_start_time field in the database.
	FieldVotingStartTime = "voting_start_time"
	// FieldVotingEndTime holds the string denoting the voting_end_time field in the database.
	FieldVotingEndTime = "voting_end_time"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// EdgeChain holds the string denoting the chain edge name in mutations.
	EdgeChain = "chain"
	// Table holds the table name of the proposal in the database.
	Table = "proposals"
	// ChainTable is the table that holds the chain relation/edge.
	ChainTable = "proposals"
	// ChainInverseTable is the table name for the Chain entity.
	// It exists in this package in order to avoid circular dependency with the "chain" package.
	ChainInverseTable = "chains"
	// ChainColumn is the table column denoting the chain relation/edge.
	ChainColumn = "chain_proposals"
)

// Columns holds all SQL columns for proposal fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldProposalID,
	FieldTitle,
	FieldDescription,
	FieldVotingStartTime,
	FieldVotingEndTime,
	FieldStatus,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "proposals"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"chain_proposals",
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
	StatusPROPOSAL_STATUS_VOTING_PERIOD  Status = "PROPOSAL_STATUS_VOTING_PERIOD"
	StatusPROPOSAL_STATUS_PASSED         Status = "PROPOSAL_STATUS_PASSED"
	StatusPROPOSAL_STATUS_REJECTED       Status = "PROPOSAL_STATUS_REJECTED"
	StatusPROPOSAL_STATUS_FAILED         Status = "PROPOSAL_STATUS_FAILED"
	StatusPROPOSAL_STATUS_UNSPECIFIED    Status = "PROPOSAL_STATUS_UNSPECIFIED"
	StatusPROPOSAL_STATUS_DEPOSIT_PERIOD Status = "PROPOSAL_STATUS_DEPOSIT_PERIOD"
)

func (s Status) String() string {
	return string(s)
}

// StatusValidator is a validator for the "status" field enum values. It is called by the builders before save.
func StatusValidator(s Status) error {
	switch s {
	case StatusPROPOSAL_STATUS_VOTING_PERIOD, StatusPROPOSAL_STATUS_PASSED, StatusPROPOSAL_STATUS_REJECTED, StatusPROPOSAL_STATUS_FAILED, StatusPROPOSAL_STATUS_UNSPECIFIED, StatusPROPOSAL_STATUS_DEPOSIT_PERIOD:
		return nil
	default:
		return fmt.Errorf("proposal: invalid enum value for status field: %q", s)
	}
}

// OrderOption defines the ordering options for the Proposal queries.
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

// ByVotingStartTime orders the results by the voting_start_time field.
func ByVotingStartTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVotingStartTime, opts...).ToFunc()
}

// ByVotingEndTime orders the results by the voting_end_time field.
func ByVotingEndTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVotingEndTime, opts...).ToFunc()
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
