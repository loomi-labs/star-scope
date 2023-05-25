// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/loomi-labs/star-scope/ent/chain"
)

// Chain is the model entity for the Chain schema.
type Chain struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// ChainID holds the value of the "chain_id" field.
	ChainID string `json:"chain_id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// PrettyName holds the value of the "pretty_name" field.
	PrettyName string `json:"pretty_name,omitempty"`
	// Path holds the value of the "path" field.
	Path string `json:"path,omitempty"`
	// Image holds the value of the "image" field.
	Image string `json:"image,omitempty"`
	// Bech32Prefix holds the value of the "bech32_prefix" field.
	Bech32Prefix string `json:"bech32_prefix,omitempty"`
	// IndexingHeight holds the value of the "indexing_height" field.
	IndexingHeight uint64 `json:"indexing_height,omitempty"`
	// HasCustomIndexer holds the value of the "has_custom_indexer" field.
	HasCustomIndexer bool `json:"has_custom_indexer,omitempty"`
	// HandledMessageTypes holds the value of the "handled_message_types" field.
	HandledMessageTypes string `json:"handled_message_types,omitempty"`
	// UnhandledMessageTypes holds the value of the "unhandled_message_types" field.
	UnhandledMessageTypes string `json:"unhandled_message_types,omitempty"`
	// IsEnabled holds the value of the "is_enabled" field.
	IsEnabled bool `json:"is_enabled,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ChainQuery when eager-loading is set.
	Edges        ChainEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ChainEdges holds the relations/edges for other nodes in the graph.
type ChainEdges struct {
	// EventListeners holds the value of the event_listeners edge.
	EventListeners []*EventListener `json:"event_listeners,omitempty"`
	// Proposals holds the value of the proposals edge.
	Proposals []*Proposal `json:"proposals,omitempty"`
	// ContractProposals holds the value of the contract_proposals edge.
	ContractProposals []*ContractProposal `json:"contract_proposals,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// EventListenersOrErr returns the EventListeners value or an error if the edge
// was not loaded in eager-loading.
func (e ChainEdges) EventListenersOrErr() ([]*EventListener, error) {
	if e.loadedTypes[0] {
		return e.EventListeners, nil
	}
	return nil, &NotLoadedError{edge: "event_listeners"}
}

// ProposalsOrErr returns the Proposals value or an error if the edge
// was not loaded in eager-loading.
func (e ChainEdges) ProposalsOrErr() ([]*Proposal, error) {
	if e.loadedTypes[1] {
		return e.Proposals, nil
	}
	return nil, &NotLoadedError{edge: "proposals"}
}

// ContractProposalsOrErr returns the ContractProposals value or an error if the edge
// was not loaded in eager-loading.
func (e ChainEdges) ContractProposalsOrErr() ([]*ContractProposal, error) {
	if e.loadedTypes[2] {
		return e.ContractProposals, nil
	}
	return nil, &NotLoadedError{edge: "contract_proposals"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Chain) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case chain.FieldHasCustomIndexer, chain.FieldIsEnabled:
			values[i] = new(sql.NullBool)
		case chain.FieldID, chain.FieldIndexingHeight:
			values[i] = new(sql.NullInt64)
		case chain.FieldChainID, chain.FieldName, chain.FieldPrettyName, chain.FieldPath, chain.FieldImage, chain.FieldBech32Prefix, chain.FieldHandledMessageTypes, chain.FieldUnhandledMessageTypes:
			values[i] = new(sql.NullString)
		case chain.FieldCreateTime, chain.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Chain fields.
func (c *Chain) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case chain.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			c.ID = int(value.Int64)
		case chain.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field create_time", values[i])
			} else if value.Valid {
				c.CreateTime = value.Time
			}
		case chain.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field update_time", values[i])
			} else if value.Valid {
				c.UpdateTime = value.Time
			}
		case chain.FieldChainID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field chain_id", values[i])
			} else if value.Valid {
				c.ChainID = value.String
			}
		case chain.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				c.Name = value.String
			}
		case chain.FieldPrettyName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field pretty_name", values[i])
			} else if value.Valid {
				c.PrettyName = value.String
			}
		case chain.FieldPath:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field path", values[i])
			} else if value.Valid {
				c.Path = value.String
			}
		case chain.FieldImage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field image", values[i])
			} else if value.Valid {
				c.Image = value.String
			}
		case chain.FieldBech32Prefix:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field bech32_prefix", values[i])
			} else if value.Valid {
				c.Bech32Prefix = value.String
			}
		case chain.FieldIndexingHeight:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field indexing_height", values[i])
			} else if value.Valid {
				c.IndexingHeight = uint64(value.Int64)
			}
		case chain.FieldHasCustomIndexer:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field has_custom_indexer", values[i])
			} else if value.Valid {
				c.HasCustomIndexer = value.Bool
			}
		case chain.FieldHandledMessageTypes:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field handled_message_types", values[i])
			} else if value.Valid {
				c.HandledMessageTypes = value.String
			}
		case chain.FieldUnhandledMessageTypes:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field unhandled_message_types", values[i])
			} else if value.Valid {
				c.UnhandledMessageTypes = value.String
			}
		case chain.FieldIsEnabled:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_enabled", values[i])
			} else if value.Valid {
				c.IsEnabled = value.Bool
			}
		default:
			c.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Chain.
// This includes values selected through modifiers, order, etc.
func (c *Chain) Value(name string) (ent.Value, error) {
	return c.selectValues.Get(name)
}

// QueryEventListeners queries the "event_listeners" edge of the Chain entity.
func (c *Chain) QueryEventListeners() *EventListenerQuery {
	return NewChainClient(c.config).QueryEventListeners(c)
}

// QueryProposals queries the "proposals" edge of the Chain entity.
func (c *Chain) QueryProposals() *ProposalQuery {
	return NewChainClient(c.config).QueryProposals(c)
}

// QueryContractProposals queries the "contract_proposals" edge of the Chain entity.
func (c *Chain) QueryContractProposals() *ContractProposalQuery {
	return NewChainClient(c.config).QueryContractProposals(c)
}

// Update returns a builder for updating this Chain.
// Note that you need to call Chain.Unwrap() before calling this method if this Chain
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Chain) Update() *ChainUpdateOne {
	return NewChainClient(c.config).UpdateOne(c)
}

// Unwrap unwraps the Chain entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Chain) Unwrap() *Chain {
	_tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Chain is not a transactional entity")
	}
	c.config.driver = _tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Chain) String() string {
	var builder strings.Builder
	builder.WriteString("Chain(")
	builder.WriteString(fmt.Sprintf("id=%v, ", c.ID))
	builder.WriteString("create_time=")
	builder.WriteString(c.CreateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("update_time=")
	builder.WriteString(c.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("chain_id=")
	builder.WriteString(c.ChainID)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(c.Name)
	builder.WriteString(", ")
	builder.WriteString("pretty_name=")
	builder.WriteString(c.PrettyName)
	builder.WriteString(", ")
	builder.WriteString("path=")
	builder.WriteString(c.Path)
	builder.WriteString(", ")
	builder.WriteString("image=")
	builder.WriteString(c.Image)
	builder.WriteString(", ")
	builder.WriteString("bech32_prefix=")
	builder.WriteString(c.Bech32Prefix)
	builder.WriteString(", ")
	builder.WriteString("indexing_height=")
	builder.WriteString(fmt.Sprintf("%v", c.IndexingHeight))
	builder.WriteString(", ")
	builder.WriteString("has_custom_indexer=")
	builder.WriteString(fmt.Sprintf("%v", c.HasCustomIndexer))
	builder.WriteString(", ")
	builder.WriteString("handled_message_types=")
	builder.WriteString(c.HandledMessageTypes)
	builder.WriteString(", ")
	builder.WriteString("unhandled_message_types=")
	builder.WriteString(c.UnhandledMessageTypes)
	builder.WriteString(", ")
	builder.WriteString("is_enabled=")
	builder.WriteString(fmt.Sprintf("%v", c.IsEnabled))
	builder.WriteByte(')')
	return builder.String()
}

// Chains is a parsable slice of Chain.
type Chains []*Chain
