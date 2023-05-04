// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/loomi-labs/star-scope/ent/channel"
	"github.com/loomi-labs/star-scope/ent/project"
)

// Channel is the model entity for the Channel schema.
type Channel struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Type holds the value of the "type" field.
	Type channel.Type `json:"type,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ChannelQuery when eager-loading is set.
	Edges            ChannelEdges `json:"edges"`
	project_channels *int
	selectValues     sql.SelectValues
}

// ChannelEdges holds the relations/edges for other nodes in the graph.
type ChannelEdges struct {
	// Project holds the value of the project edge.
	Project *Project `json:"project,omitempty"`
	// EventListeners holds the value of the event_listeners edge.
	EventListeners []*EventListener `json:"event_listeners,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// ProjectOrErr returns the Project value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ChannelEdges) ProjectOrErr() (*Project, error) {
	if e.loadedTypes[0] {
		if e.Project == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: project.Label}
		}
		return e.Project, nil
	}
	return nil, &NotLoadedError{edge: "project"}
}

// EventListenersOrErr returns the EventListeners value or an error if the edge
// was not loaded in eager-loading.
func (e ChannelEdges) EventListenersOrErr() ([]*EventListener, error) {
	if e.loadedTypes[1] {
		return e.EventListeners, nil
	}
	return nil, &NotLoadedError{edge: "event_listeners"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Channel) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case channel.FieldID:
			values[i] = new(sql.NullInt64)
		case channel.FieldName, channel.FieldType:
			values[i] = new(sql.NullString)
		case channel.FieldCreateTime, channel.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		case channel.ForeignKeys[0]: // project_channels
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Channel fields.
func (c *Channel) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case channel.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			c.ID = int(value.Int64)
		case channel.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field create_time", values[i])
			} else if value.Valid {
				c.CreateTime = value.Time
			}
		case channel.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field update_time", values[i])
			} else if value.Valid {
				c.UpdateTime = value.Time
			}
		case channel.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				c.Name = value.String
			}
		case channel.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				c.Type = channel.Type(value.String)
			}
		case channel.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field project_channels", value)
			} else if value.Valid {
				c.project_channels = new(int)
				*c.project_channels = int(value.Int64)
			}
		default:
			c.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Channel.
// This includes values selected through modifiers, order, etc.
func (c *Channel) Value(name string) (ent.Value, error) {
	return c.selectValues.Get(name)
}

// QueryProject queries the "project" edge of the Channel entity.
func (c *Channel) QueryProject() *ProjectQuery {
	return NewChannelClient(c.config).QueryProject(c)
}

// QueryEventListeners queries the "event_listeners" edge of the Channel entity.
func (c *Channel) QueryEventListeners() *EventListenerQuery {
	return NewChannelClient(c.config).QueryEventListeners(c)
}

// Update returns a builder for updating this Channel.
// Note that you need to call Channel.Unwrap() before calling this method if this Channel
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Channel) Update() *ChannelUpdateOne {
	return NewChannelClient(c.config).UpdateOne(c)
}

// Unwrap unwraps the Channel entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Channel) Unwrap() *Channel {
	_tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Channel is not a transactional entity")
	}
	c.config.driver = _tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Channel) String() string {
	var builder strings.Builder
	builder.WriteString("Channel(")
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
	builder.WriteString("type=")
	builder.WriteString(fmt.Sprintf("%v", c.Type))
	builder.WriteByte(')')
	return builder.String()
}

// Channels is a parsable slice of Channel.
type Channels []*Channel
