// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/shifty11/blocklog-backend/ent/chain"
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
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Image holds the value of the "image" field.
	Image string `json:"image,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ChainQuery when eager-loading is set.
	Edges        ChainEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ChainEdges holds the relations/edges for other nodes in the graph.
type ChainEdges struct {
	// EventListeners holds the value of the event_listeners edge.
	EventListeners []*EventListener `json:"event_listeners,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// EventListenersOrErr returns the EventListeners value or an error if the edge
// was not loaded in eager-loading.
func (e ChainEdges) EventListenersOrErr() ([]*EventListener, error) {
	if e.loadedTypes[0] {
		return e.EventListeners, nil
	}
	return nil, &NotLoadedError{edge: "event_listeners"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Chain) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case chain.FieldID:
			values[i] = new(sql.NullInt64)
		case chain.FieldName, chain.FieldImage:
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
		case chain.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				c.Name = value.String
			}
		case chain.FieldImage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field image", values[i])
			} else if value.Valid {
				c.Image = value.String
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
	builder.WriteString("name=")
	builder.WriteString(c.Name)
	builder.WriteString(", ")
	builder.WriteString("image=")
	builder.WriteString(c.Image)
	builder.WriteByte(')')
	return builder.String()
}

// Chains is a parsable slice of Chain.
type Chains []*Chain
