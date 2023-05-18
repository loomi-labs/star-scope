// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/loomi-labs/star-scope/ent/event"
	"github.com/loomi-labs/star-scope/ent/eventlistener"
	"github.com/loomi-labs/star-scope/ent/predicate"
)

// EventUpdate is the builder for updating Event entities.
type EventUpdate struct {
	config
	hooks    []Hook
	mutation *EventMutation
}

// Where appends a list predicates to the EventUpdate builder.
func (eu *EventUpdate) Where(ps ...predicate.Event) *EventUpdate {
	eu.mutation.Where(ps...)
	return eu
}

// SetUpdateTime sets the "update_time" field.
func (eu *EventUpdate) SetUpdateTime(t time.Time) *EventUpdate {
	eu.mutation.SetUpdateTime(t)
	return eu
}

// SetType sets the "type" field.
func (eu *EventUpdate) SetType(e event.Type) *EventUpdate {
	eu.mutation.SetType(e)
	return eu
}

// SetTxEvent sets the "tx_event" field.
func (eu *EventUpdate) SetTxEvent(b []byte) *EventUpdate {
	eu.mutation.SetTxEvent(b)
	return eu
}

// SetNotifyTime sets the "notify_time" field.
func (eu *EventUpdate) SetNotifyTime(t time.Time) *EventUpdate {
	eu.mutation.SetNotifyTime(t)
	return eu
}

// SetNillableNotifyTime sets the "notify_time" field if the given value is not nil.
func (eu *EventUpdate) SetNillableNotifyTime(t *time.Time) *EventUpdate {
	if t != nil {
		eu.SetNotifyTime(*t)
	}
	return eu
}

// SetEventListenerID sets the "event_listener" edge to the EventListener entity by ID.
func (eu *EventUpdate) SetEventListenerID(id int) *EventUpdate {
	eu.mutation.SetEventListenerID(id)
	return eu
}

// SetNillableEventListenerID sets the "event_listener" edge to the EventListener entity by ID if the given value is not nil.
func (eu *EventUpdate) SetNillableEventListenerID(id *int) *EventUpdate {
	if id != nil {
		eu = eu.SetEventListenerID(*id)
	}
	return eu
}

// SetEventListener sets the "event_listener" edge to the EventListener entity.
func (eu *EventUpdate) SetEventListener(e *EventListener) *EventUpdate {
	return eu.SetEventListenerID(e.ID)
}

// Mutation returns the EventMutation object of the builder.
func (eu *EventUpdate) Mutation() *EventMutation {
	return eu.mutation
}

// ClearEventListener clears the "event_listener" edge to the EventListener entity.
func (eu *EventUpdate) ClearEventListener() *EventUpdate {
	eu.mutation.ClearEventListener()
	return eu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (eu *EventUpdate) Save(ctx context.Context) (int, error) {
	eu.defaults()
	return withHooks(ctx, eu.sqlSave, eu.mutation, eu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (eu *EventUpdate) SaveX(ctx context.Context) int {
	affected, err := eu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (eu *EventUpdate) Exec(ctx context.Context) error {
	_, err := eu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eu *EventUpdate) ExecX(ctx context.Context) {
	if err := eu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (eu *EventUpdate) defaults() {
	if _, ok := eu.mutation.UpdateTime(); !ok {
		v := event.UpdateDefaultUpdateTime()
		eu.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (eu *EventUpdate) check() error {
	if v, ok := eu.mutation.GetType(); ok {
		if err := event.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "Event.type": %w`, err)}
		}
	}
	return nil
}

func (eu *EventUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := eu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(event.Table, event.Columns, sqlgraph.NewFieldSpec(event.FieldID, field.TypeInt))
	if ps := eu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := eu.mutation.UpdateTime(); ok {
		_spec.SetField(event.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := eu.mutation.GetType(); ok {
		_spec.SetField(event.FieldType, field.TypeEnum, value)
	}
	if value, ok := eu.mutation.TxEvent(); ok {
		_spec.SetField(event.FieldTxEvent, field.TypeBytes, value)
	}
	if value, ok := eu.mutation.NotifyTime(); ok {
		_spec.SetField(event.FieldNotifyTime, field.TypeTime, value)
	}
	if eu.mutation.EventListenerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   event.EventListenerTable,
			Columns: []string{event.EventListenerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(eventlistener.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.EventListenerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   event.EventListenerTable,
			Columns: []string{event.EventListenerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(eventlistener.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, eu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{event.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	eu.mutation.done = true
	return n, nil
}

// EventUpdateOne is the builder for updating a single Event entity.
type EventUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *EventMutation
}

// SetUpdateTime sets the "update_time" field.
func (euo *EventUpdateOne) SetUpdateTime(t time.Time) *EventUpdateOne {
	euo.mutation.SetUpdateTime(t)
	return euo
}

// SetType sets the "type" field.
func (euo *EventUpdateOne) SetType(e event.Type) *EventUpdateOne {
	euo.mutation.SetType(e)
	return euo
}

// SetTxEvent sets the "tx_event" field.
func (euo *EventUpdateOne) SetTxEvent(b []byte) *EventUpdateOne {
	euo.mutation.SetTxEvent(b)
	return euo
}

// SetNotifyTime sets the "notify_time" field.
func (euo *EventUpdateOne) SetNotifyTime(t time.Time) *EventUpdateOne {
	euo.mutation.SetNotifyTime(t)
	return euo
}

// SetNillableNotifyTime sets the "notify_time" field if the given value is not nil.
func (euo *EventUpdateOne) SetNillableNotifyTime(t *time.Time) *EventUpdateOne {
	if t != nil {
		euo.SetNotifyTime(*t)
	}
	return euo
}

// SetEventListenerID sets the "event_listener" edge to the EventListener entity by ID.
func (euo *EventUpdateOne) SetEventListenerID(id int) *EventUpdateOne {
	euo.mutation.SetEventListenerID(id)
	return euo
}

// SetNillableEventListenerID sets the "event_listener" edge to the EventListener entity by ID if the given value is not nil.
func (euo *EventUpdateOne) SetNillableEventListenerID(id *int) *EventUpdateOne {
	if id != nil {
		euo = euo.SetEventListenerID(*id)
	}
	return euo
}

// SetEventListener sets the "event_listener" edge to the EventListener entity.
func (euo *EventUpdateOne) SetEventListener(e *EventListener) *EventUpdateOne {
	return euo.SetEventListenerID(e.ID)
}

// Mutation returns the EventMutation object of the builder.
func (euo *EventUpdateOne) Mutation() *EventMutation {
	return euo.mutation
}

// ClearEventListener clears the "event_listener" edge to the EventListener entity.
func (euo *EventUpdateOne) ClearEventListener() *EventUpdateOne {
	euo.mutation.ClearEventListener()
	return euo
}

// Where appends a list predicates to the EventUpdate builder.
func (euo *EventUpdateOne) Where(ps ...predicate.Event) *EventUpdateOne {
	euo.mutation.Where(ps...)
	return euo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (euo *EventUpdateOne) Select(field string, fields ...string) *EventUpdateOne {
	euo.fields = append([]string{field}, fields...)
	return euo
}

// Save executes the query and returns the updated Event entity.
func (euo *EventUpdateOne) Save(ctx context.Context) (*Event, error) {
	euo.defaults()
	return withHooks(ctx, euo.sqlSave, euo.mutation, euo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (euo *EventUpdateOne) SaveX(ctx context.Context) *Event {
	node, err := euo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (euo *EventUpdateOne) Exec(ctx context.Context) error {
	_, err := euo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (euo *EventUpdateOne) ExecX(ctx context.Context) {
	if err := euo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (euo *EventUpdateOne) defaults() {
	if _, ok := euo.mutation.UpdateTime(); !ok {
		v := event.UpdateDefaultUpdateTime()
		euo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (euo *EventUpdateOne) check() error {
	if v, ok := euo.mutation.GetType(); ok {
		if err := event.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "Event.type": %w`, err)}
		}
	}
	return nil
}

func (euo *EventUpdateOne) sqlSave(ctx context.Context) (_node *Event, err error) {
	if err := euo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(event.Table, event.Columns, sqlgraph.NewFieldSpec(event.FieldID, field.TypeInt))
	id, ok := euo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Event.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := euo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, event.FieldID)
		for _, f := range fields {
			if !event.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != event.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := euo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := euo.mutation.UpdateTime(); ok {
		_spec.SetField(event.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := euo.mutation.GetType(); ok {
		_spec.SetField(event.FieldType, field.TypeEnum, value)
	}
	if value, ok := euo.mutation.TxEvent(); ok {
		_spec.SetField(event.FieldTxEvent, field.TypeBytes, value)
	}
	if value, ok := euo.mutation.NotifyTime(); ok {
		_spec.SetField(event.FieldNotifyTime, field.TypeTime, value)
	}
	if euo.mutation.EventListenerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   event.EventListenerTable,
			Columns: []string{event.EventListenerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(eventlistener.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.EventListenerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   event.EventListenerTable,
			Columns: []string{event.EventListenerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(eventlistener.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Event{config: euo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, euo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{event.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	euo.mutation.done = true
	return _node, nil
}
