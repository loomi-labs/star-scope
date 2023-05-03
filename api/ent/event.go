// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/shifty11/blocklog-backend/ent/event"
	"github.com/shifty11/blocklog-backend/ent/eventlistener"
)

// Event is the model entity for the Event schema.
type Event struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Type holds the value of the "type" field.
	Type event.Type `json:"type,omitempty"`
	// TxEvent holds the value of the "tx_event" field.
	TxEvent []byte `json:"tx_event,omitempty"`
	// NotifyTime holds the value of the "notify_time" field.
	NotifyTime time.Time `json:"notify_time,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the EventQuery when eager-loading is set.
	Edges                 EventEdges `json:"edges"`
	event_listener_events *int
	selectValues          sql.SelectValues
}

// EventEdges holds the relations/edges for other nodes in the graph.
type EventEdges struct {
	// EventListener holds the value of the event_listener edge.
	EventListener *EventListener `json:"event_listener,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// EventListenerOrErr returns the EventListener value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EventEdges) EventListenerOrErr() (*EventListener, error) {
	if e.loadedTypes[0] {
		if e.EventListener == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: eventlistener.Label}
		}
		return e.EventListener, nil
	}
	return nil, &NotLoadedError{edge: "event_listener"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Event) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case event.FieldTxEvent:
			values[i] = new([]byte)
		case event.FieldID:
			values[i] = new(sql.NullInt64)
		case event.FieldType:
			values[i] = new(sql.NullString)
		case event.FieldCreateTime, event.FieldUpdateTime, event.FieldNotifyTime:
			values[i] = new(sql.NullTime)
		case event.ForeignKeys[0]: // event_listener_events
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Event fields.
func (e *Event) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case event.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			e.ID = int(value.Int64)
		case event.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field create_time", values[i])
			} else if value.Valid {
				e.CreateTime = value.Time
			}
		case event.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field update_time", values[i])
			} else if value.Valid {
				e.UpdateTime = value.Time
			}
		case event.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				e.Type = event.Type(value.String)
			}
		case event.FieldTxEvent:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field tx_event", values[i])
			} else if value != nil {
				e.TxEvent = *value
			}
		case event.FieldNotifyTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field notify_time", values[i])
			} else if value.Valid {
				e.NotifyTime = value.Time
			}
		case event.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field event_listener_events", value)
			} else if value.Valid {
				e.event_listener_events = new(int)
				*e.event_listener_events = int(value.Int64)
			}
		default:
			e.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Event.
// This includes values selected through modifiers, order, etc.
func (e *Event) Value(name string) (ent.Value, error) {
	return e.selectValues.Get(name)
}

// QueryEventListener queries the "event_listener" edge of the Event entity.
func (e *Event) QueryEventListener() *EventListenerQuery {
	return NewEventClient(e.config).QueryEventListener(e)
}

// Update returns a builder for updating this Event.
// Note that you need to call Event.Unwrap() before calling this method if this Event
// was returned from a transaction, and the transaction was committed or rolled back.
func (e *Event) Update() *EventUpdateOne {
	return NewEventClient(e.config).UpdateOne(e)
}

// Unwrap unwraps the Event entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (e *Event) Unwrap() *Event {
	_tx, ok := e.config.driver.(*txDriver)
	if !ok {
		panic("ent: Event is not a transactional entity")
	}
	e.config.driver = _tx.drv
	return e
}

// String implements the fmt.Stringer.
func (e *Event) String() string {
	var builder strings.Builder
	builder.WriteString("Event(")
	builder.WriteString(fmt.Sprintf("id=%v, ", e.ID))
	builder.WriteString("create_time=")
	builder.WriteString(e.CreateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("update_time=")
	builder.WriteString(e.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(fmt.Sprintf("%v", e.Type))
	builder.WriteString(", ")
	builder.WriteString("tx_event=")
	builder.WriteString(fmt.Sprintf("%v", e.TxEvent))
	builder.WriteString(", ")
	builder.WriteString("notify_time=")
	builder.WriteString(e.NotifyTime.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// Events is a parsable slice of Event.
type Events []*Event
