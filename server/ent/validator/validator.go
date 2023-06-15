// Code generated by ent, DO NOT EDIT.

package validator

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the validator type in the database.
	Label = "validator"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time field in the database.
	FieldUpdateTime = "update_time"
	// FieldOperatorAddress holds the string denoting the operator_address field in the database.
	FieldOperatorAddress = "operator_address"
	// FieldAddress holds the string denoting the address field in the database.
	FieldAddress = "address"
	// FieldMoniker holds the string denoting the moniker field in the database.
	FieldMoniker = "moniker"
	// FieldFirstInactiveTime holds the string denoting the first_inactive_time field in the database.
	FieldFirstInactiveTime = "first_inactive_time"
	// FieldLastSlashValidatorPeriod holds the string denoting the last_slash_validator_period field in the database.
	FieldLastSlashValidatorPeriod = "last_slash_validator_period"
	// EdgeChain holds the string denoting the chain edge name in mutations.
	EdgeChain = "chain"
	// Table holds the table name of the validator in the database.
	Table = "validators"
	// ChainTable is the table that holds the chain relation/edge.
	ChainTable = "validators"
	// ChainInverseTable is the table name for the Chain entity.
	// It exists in this package in order to avoid circular dependency with the "chain" package.
	ChainInverseTable = "chains"
	// ChainColumn is the table column denoting the chain relation/edge.
	ChainColumn = "chain_validators"
)

// Columns holds all SQL columns for validator fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldOperatorAddress,
	FieldAddress,
	FieldMoniker,
	FieldFirstInactiveTime,
	FieldLastSlashValidatorPeriod,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "validators"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"chain_validators",
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
	// OperatorAddressValidator is a validator for the "operator_address" field. It is called by the builders before save.
	OperatorAddressValidator func(string) error
	// AddressValidator is a validator for the "address" field. It is called by the builders before save.
	AddressValidator func(string) error
)

// OrderOption defines the ordering options for the Validator queries.
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

// ByOperatorAddress orders the results by the operator_address field.
func ByOperatorAddress(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOperatorAddress, opts...).ToFunc()
}

// ByAddress orders the results by the address field.
func ByAddress(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAddress, opts...).ToFunc()
}

// ByMoniker orders the results by the moniker field.
func ByMoniker(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMoniker, opts...).ToFunc()
}

// ByFirstInactiveTime orders the results by the first_inactive_time field.
func ByFirstInactiveTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFirstInactiveTime, opts...).ToFunc()
}

// ByLastSlashValidatorPeriod orders the results by the last_slash_validator_period field.
func ByLastSlashValidatorPeriod(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLastSlashValidatorPeriod, opts...).ToFunc()
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
