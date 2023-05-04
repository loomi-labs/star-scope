// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/loomi-labs/star-scope/ent/event"
	"github.com/loomi-labs/star-scope/ent/eventlistener"
)

// EventCreate is the builder for creating a Event entity.
type EventCreate struct {
	config
	mutation *EventMutation
	hooks    []Hook
}

// SetCreateTime sets the "create_time" field.
func (ec *EventCreate) SetCreateTime(t time.Time) *EventCreate {
	ec.mutation.SetCreateTime(t)
	return ec
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (ec *EventCreate) SetNillableCreateTime(t *time.Time) *EventCreate {
	if t != nil {
		ec.SetCreateTime(*t)
	}
	return ec
}

// SetUpdateTime sets the "update_time" field.
func (ec *EventCreate) SetUpdateTime(t time.Time) *EventCreate {
	ec.mutation.SetUpdateTime(t)
	return ec
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (ec *EventCreate) SetNillableUpdateTime(t *time.Time) *EventCreate {
	if t != nil {
		ec.SetUpdateTime(*t)
	}
	return ec
}

// SetType sets the "type" field.
func (ec *EventCreate) SetType(e event.Type) *EventCreate {
	ec.mutation.SetType(e)
	return ec
}

// SetTxEvent sets the "tx_event" field.
func (ec *EventCreate) SetTxEvent(b []byte) *EventCreate {
	ec.mutation.SetTxEvent(b)
	return ec
}

// SetNotifyTime sets the "notify_time" field.
func (ec *EventCreate) SetNotifyTime(t time.Time) *EventCreate {
	ec.mutation.SetNotifyTime(t)
	return ec
}

// SetNillableNotifyTime sets the "notify_time" field if the given value is not nil.
func (ec *EventCreate) SetNillableNotifyTime(t *time.Time) *EventCreate {
	if t != nil {
		ec.SetNotifyTime(*t)
	}
	return ec
}

// SetEventListenerID sets the "event_listener" edge to the EventListener entity by ID.
func (ec *EventCreate) SetEventListenerID(id int) *EventCreate {
	ec.mutation.SetEventListenerID(id)
	return ec
}

// SetNillableEventListenerID sets the "event_listener" edge to the EventListener entity by ID if the given value is not nil.
func (ec *EventCreate) SetNillableEventListenerID(id *int) *EventCreate {
	if id != nil {
		ec = ec.SetEventListenerID(*id)
	}
	return ec
}

// SetEventListener sets the "event_listener" edge to the EventListener entity.
func (ec *EventCreate) SetEventListener(e *EventListener) *EventCreate {
	return ec.SetEventListenerID(e.ID)
}

// Mutation returns the EventMutation object of the builder.
func (ec *EventCreate) Mutation() *EventMutation {
	return ec.mutation
}

// Save creates the Event in the database.
func (ec *EventCreate) Save(ctx context.Context) (*Event, error) {
	ec.defaults()
	return withHooks[*Event, EventMutation](ctx, ec.sqlSave, ec.mutation, ec.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ec *EventCreate) SaveX(ctx context.Context) *Event {
	v, err := ec.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ec *EventCreate) Exec(ctx context.Context) error {
	_, err := ec.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ec *EventCreate) ExecX(ctx context.Context) {
	if err := ec.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ec *EventCreate) defaults() {
	if _, ok := ec.mutation.CreateTime(); !ok {
		v := event.DefaultCreateTime()
		ec.mutation.SetCreateTime(v)
	}
	if _, ok := ec.mutation.UpdateTime(); !ok {
		v := event.DefaultUpdateTime()
		ec.mutation.SetUpdateTime(v)
	}
	if _, ok := ec.mutation.NotifyTime(); !ok {
		v := event.DefaultNotifyTime
		ec.mutation.SetNotifyTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ec *EventCreate) check() error {
	if _, ok := ec.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "create_time", err: errors.New(`ent: missing required field "Event.create_time"`)}
	}
	if _, ok := ec.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "update_time", err: errors.New(`ent: missing required field "Event.update_time"`)}
	}
	if _, ok := ec.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New(`ent: missing required field "Event.type"`)}
	}
	if v, ok := ec.mutation.GetType(); ok {
		if err := event.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "Event.type": %w`, err)}
		}
	}
	if _, ok := ec.mutation.TxEvent(); !ok {
		return &ValidationError{Name: "tx_event", err: errors.New(`ent: missing required field "Event.tx_event"`)}
	}
	if _, ok := ec.mutation.NotifyTime(); !ok {
		return &ValidationError{Name: "notify_time", err: errors.New(`ent: missing required field "Event.notify_time"`)}
	}
	return nil
}

func (ec *EventCreate) sqlSave(ctx context.Context) (*Event, error) {
	if err := ec.check(); err != nil {
		return nil, err
	}
	_node, _spec := ec.createSpec()
	if err := sqlgraph.CreateNode(ctx, ec.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	ec.mutation.id = &_node.ID
	ec.mutation.done = true
	return _node, nil
}

func (ec *EventCreate) createSpec() (*Event, *sqlgraph.CreateSpec) {
	var (
		_node = &Event{config: ec.config}
		_spec = sqlgraph.NewCreateSpec(event.Table, sqlgraph.NewFieldSpec(event.FieldID, field.TypeInt))
	)
	if value, ok := ec.mutation.CreateTime(); ok {
		_spec.SetField(event.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = value
	}
	if value, ok := ec.mutation.UpdateTime(); ok {
		_spec.SetField(event.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = value
	}
	if value, ok := ec.mutation.GetType(); ok {
		_spec.SetField(event.FieldType, field.TypeEnum, value)
		_node.Type = value
	}
	if value, ok := ec.mutation.TxEvent(); ok {
		_spec.SetField(event.FieldTxEvent, field.TypeBytes, value)
		_node.TxEvent = value
	}
	if value, ok := ec.mutation.NotifyTime(); ok {
		_spec.SetField(event.FieldNotifyTime, field.TypeTime, value)
		_node.NotifyTime = value
	}
	if nodes := ec.mutation.EventListenerIDs(); len(nodes) > 0 {
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
		_node.event_listener_events = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// EventCreateBulk is the builder for creating many Event entities in bulk.
type EventCreateBulk struct {
	config
	builders []*EventCreate
}

// Save creates the Event entities in the database.
func (ecb *EventCreateBulk) Save(ctx context.Context) ([]*Event, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ecb.builders))
	nodes := make([]*Event, len(ecb.builders))
	mutators := make([]Mutator, len(ecb.builders))
	for i := range ecb.builders {
		func(i int, root context.Context) {
			builder := ecb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*EventMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ecb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ecb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ecb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ecb *EventCreateBulk) SaveX(ctx context.Context) []*Event {
	v, err := ecb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ecb *EventCreateBulk) Exec(ctx context.Context) error {
	_, err := ecb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ecb *EventCreateBulk) ExecX(ctx context.Context) {
	if err := ecb.Exec(ctx); err != nil {
		panic(err)
	}
}
