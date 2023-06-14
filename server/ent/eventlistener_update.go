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
	"github.com/google/uuid"
	"github.com/loomi-labs/star-scope/ent/chain"
	"github.com/loomi-labs/star-scope/ent/event"
	"github.com/loomi-labs/star-scope/ent/eventlistener"
	"github.com/loomi-labs/star-scope/ent/predicate"
	"github.com/loomi-labs/star-scope/ent/user"
)

// EventListenerUpdate is the builder for updating EventListener entities.
type EventListenerUpdate struct {
	config
	hooks    []Hook
	mutation *EventListenerMutation
}

// Where appends a list predicates to the EventListenerUpdate builder.
func (elu *EventListenerUpdate) Where(ps ...predicate.EventListener) *EventListenerUpdate {
	elu.mutation.Where(ps...)
	return elu
}

// SetUpdateTime sets the "update_time" field.
func (elu *EventListenerUpdate) SetUpdateTime(t time.Time) *EventListenerUpdate {
	elu.mutation.SetUpdateTime(t)
	return elu
}

// SetWalletAddress sets the "wallet_address" field.
func (elu *EventListenerUpdate) SetWalletAddress(s string) *EventListenerUpdate {
	elu.mutation.SetWalletAddress(s)
	return elu
}

// SetNillableWalletAddress sets the "wallet_address" field if the given value is not nil.
func (elu *EventListenerUpdate) SetNillableWalletAddress(s *string) *EventListenerUpdate {
	if s != nil {
		elu.SetWalletAddress(*s)
	}
	return elu
}

// ClearWalletAddress clears the value of the "wallet_address" field.
func (elu *EventListenerUpdate) ClearWalletAddress() *EventListenerUpdate {
	elu.mutation.ClearWalletAddress()
	return elu
}

// SetDataType sets the "data_type" field.
func (elu *EventListenerUpdate) SetDataType(et eventlistener.DataType) *EventListenerUpdate {
	elu.mutation.SetDataType(et)
	return elu
}

// SetUserID sets the "user" edge to the User entity by ID.
func (elu *EventListenerUpdate) SetUserID(id int) *EventListenerUpdate {
	elu.mutation.SetUserID(id)
	return elu
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (elu *EventListenerUpdate) SetNillableUserID(id *int) *EventListenerUpdate {
	if id != nil {
		elu = elu.SetUserID(*id)
	}
	return elu
}

// SetUser sets the "user" edge to the User entity.
func (elu *EventListenerUpdate) SetUser(u *User) *EventListenerUpdate {
	return elu.SetUserID(u.ID)
}

// SetChainID sets the "chain" edge to the Chain entity by ID.
func (elu *EventListenerUpdate) SetChainID(id int) *EventListenerUpdate {
	elu.mutation.SetChainID(id)
	return elu
}

// SetNillableChainID sets the "chain" edge to the Chain entity by ID if the given value is not nil.
func (elu *EventListenerUpdate) SetNillableChainID(id *int) *EventListenerUpdate {
	if id != nil {
		elu = elu.SetChainID(*id)
	}
	return elu
}

// SetChain sets the "chain" edge to the Chain entity.
func (elu *EventListenerUpdate) SetChain(c *Chain) *EventListenerUpdate {
	return elu.SetChainID(c.ID)
}

// AddEventIDs adds the "events" edge to the Event entity by IDs.
func (elu *EventListenerUpdate) AddEventIDs(ids ...uuid.UUID) *EventListenerUpdate {
	elu.mutation.AddEventIDs(ids...)
	return elu
}

// AddEvents adds the "events" edges to the Event entity.
func (elu *EventListenerUpdate) AddEvents(e ...*Event) *EventListenerUpdate {
	ids := make([]uuid.UUID, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return elu.AddEventIDs(ids...)
}

// Mutation returns the EventListenerMutation object of the builder.
func (elu *EventListenerUpdate) Mutation() *EventListenerMutation {
	return elu.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (elu *EventListenerUpdate) ClearUser() *EventListenerUpdate {
	elu.mutation.ClearUser()
	return elu
}

// ClearChain clears the "chain" edge to the Chain entity.
func (elu *EventListenerUpdate) ClearChain() *EventListenerUpdate {
	elu.mutation.ClearChain()
	return elu
}

// ClearEvents clears all "events" edges to the Event entity.
func (elu *EventListenerUpdate) ClearEvents() *EventListenerUpdate {
	elu.mutation.ClearEvents()
	return elu
}

// RemoveEventIDs removes the "events" edge to Event entities by IDs.
func (elu *EventListenerUpdate) RemoveEventIDs(ids ...uuid.UUID) *EventListenerUpdate {
	elu.mutation.RemoveEventIDs(ids...)
	return elu
}

// RemoveEvents removes "events" edges to Event entities.
func (elu *EventListenerUpdate) RemoveEvents(e ...*Event) *EventListenerUpdate {
	ids := make([]uuid.UUID, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return elu.RemoveEventIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (elu *EventListenerUpdate) Save(ctx context.Context) (int, error) {
	elu.defaults()
	return withHooks(ctx, elu.sqlSave, elu.mutation, elu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (elu *EventListenerUpdate) SaveX(ctx context.Context) int {
	affected, err := elu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (elu *EventListenerUpdate) Exec(ctx context.Context) error {
	_, err := elu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (elu *EventListenerUpdate) ExecX(ctx context.Context) {
	if err := elu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (elu *EventListenerUpdate) defaults() {
	if _, ok := elu.mutation.UpdateTime(); !ok {
		v := eventlistener.UpdateDefaultUpdateTime()
		elu.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (elu *EventListenerUpdate) check() error {
	if v, ok := elu.mutation.DataType(); ok {
		if err := eventlistener.DataTypeValidator(v); err != nil {
			return &ValidationError{Name: "data_type", err: fmt.Errorf(`ent: validator failed for field "EventListener.data_type": %w`, err)}
		}
	}
	return nil
}

func (elu *EventListenerUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := elu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(eventlistener.Table, eventlistener.Columns, sqlgraph.NewFieldSpec(eventlistener.FieldID, field.TypeInt))
	if ps := elu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := elu.mutation.UpdateTime(); ok {
		_spec.SetField(eventlistener.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := elu.mutation.WalletAddress(); ok {
		_spec.SetField(eventlistener.FieldWalletAddress, field.TypeString, value)
	}
	if elu.mutation.WalletAddressCleared() {
		_spec.ClearField(eventlistener.FieldWalletAddress, field.TypeString)
	}
	if value, ok := elu.mutation.DataType(); ok {
		_spec.SetField(eventlistener.FieldDataType, field.TypeEnum, value)
	}
	if elu.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   eventlistener.UserTable,
			Columns: []string{eventlistener.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := elu.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   eventlistener.UserTable,
			Columns: []string{eventlistener.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if elu.mutation.ChainCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   eventlistener.ChainTable,
			Columns: []string{eventlistener.ChainColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chain.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := elu.mutation.ChainIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   eventlistener.ChainTable,
			Columns: []string{eventlistener.ChainColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chain.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if elu.mutation.EventsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   eventlistener.EventsTable,
			Columns: []string{eventlistener.EventsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(event.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := elu.mutation.RemovedEventsIDs(); len(nodes) > 0 && !elu.mutation.EventsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   eventlistener.EventsTable,
			Columns: []string{eventlistener.EventsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(event.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := elu.mutation.EventsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   eventlistener.EventsTable,
			Columns: []string{eventlistener.EventsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(event.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, elu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{eventlistener.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	elu.mutation.done = true
	return n, nil
}

// EventListenerUpdateOne is the builder for updating a single EventListener entity.
type EventListenerUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *EventListenerMutation
}

// SetUpdateTime sets the "update_time" field.
func (eluo *EventListenerUpdateOne) SetUpdateTime(t time.Time) *EventListenerUpdateOne {
	eluo.mutation.SetUpdateTime(t)
	return eluo
}

// SetWalletAddress sets the "wallet_address" field.
func (eluo *EventListenerUpdateOne) SetWalletAddress(s string) *EventListenerUpdateOne {
	eluo.mutation.SetWalletAddress(s)
	return eluo
}

// SetNillableWalletAddress sets the "wallet_address" field if the given value is not nil.
func (eluo *EventListenerUpdateOne) SetNillableWalletAddress(s *string) *EventListenerUpdateOne {
	if s != nil {
		eluo.SetWalletAddress(*s)
	}
	return eluo
}

// ClearWalletAddress clears the value of the "wallet_address" field.
func (eluo *EventListenerUpdateOne) ClearWalletAddress() *EventListenerUpdateOne {
	eluo.mutation.ClearWalletAddress()
	return eluo
}

// SetDataType sets the "data_type" field.
func (eluo *EventListenerUpdateOne) SetDataType(et eventlistener.DataType) *EventListenerUpdateOne {
	eluo.mutation.SetDataType(et)
	return eluo
}

// SetUserID sets the "user" edge to the User entity by ID.
func (eluo *EventListenerUpdateOne) SetUserID(id int) *EventListenerUpdateOne {
	eluo.mutation.SetUserID(id)
	return eluo
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (eluo *EventListenerUpdateOne) SetNillableUserID(id *int) *EventListenerUpdateOne {
	if id != nil {
		eluo = eluo.SetUserID(*id)
	}
	return eluo
}

// SetUser sets the "user" edge to the User entity.
func (eluo *EventListenerUpdateOne) SetUser(u *User) *EventListenerUpdateOne {
	return eluo.SetUserID(u.ID)
}

// SetChainID sets the "chain" edge to the Chain entity by ID.
func (eluo *EventListenerUpdateOne) SetChainID(id int) *EventListenerUpdateOne {
	eluo.mutation.SetChainID(id)
	return eluo
}

// SetNillableChainID sets the "chain" edge to the Chain entity by ID if the given value is not nil.
func (eluo *EventListenerUpdateOne) SetNillableChainID(id *int) *EventListenerUpdateOne {
	if id != nil {
		eluo = eluo.SetChainID(*id)
	}
	return eluo
}

// SetChain sets the "chain" edge to the Chain entity.
func (eluo *EventListenerUpdateOne) SetChain(c *Chain) *EventListenerUpdateOne {
	return eluo.SetChainID(c.ID)
}

// AddEventIDs adds the "events" edge to the Event entity by IDs.
func (eluo *EventListenerUpdateOne) AddEventIDs(ids ...uuid.UUID) *EventListenerUpdateOne {
	eluo.mutation.AddEventIDs(ids...)
	return eluo
}

// AddEvents adds the "events" edges to the Event entity.
func (eluo *EventListenerUpdateOne) AddEvents(e ...*Event) *EventListenerUpdateOne {
	ids := make([]uuid.UUID, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eluo.AddEventIDs(ids...)
}

// Mutation returns the EventListenerMutation object of the builder.
func (eluo *EventListenerUpdateOne) Mutation() *EventListenerMutation {
	return eluo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (eluo *EventListenerUpdateOne) ClearUser() *EventListenerUpdateOne {
	eluo.mutation.ClearUser()
	return eluo
}

// ClearChain clears the "chain" edge to the Chain entity.
func (eluo *EventListenerUpdateOne) ClearChain() *EventListenerUpdateOne {
	eluo.mutation.ClearChain()
	return eluo
}

// ClearEvents clears all "events" edges to the Event entity.
func (eluo *EventListenerUpdateOne) ClearEvents() *EventListenerUpdateOne {
	eluo.mutation.ClearEvents()
	return eluo
}

// RemoveEventIDs removes the "events" edge to Event entities by IDs.
func (eluo *EventListenerUpdateOne) RemoveEventIDs(ids ...uuid.UUID) *EventListenerUpdateOne {
	eluo.mutation.RemoveEventIDs(ids...)
	return eluo
}

// RemoveEvents removes "events" edges to Event entities.
func (eluo *EventListenerUpdateOne) RemoveEvents(e ...*Event) *EventListenerUpdateOne {
	ids := make([]uuid.UUID, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eluo.RemoveEventIDs(ids...)
}

// Where appends a list predicates to the EventListenerUpdate builder.
func (eluo *EventListenerUpdateOne) Where(ps ...predicate.EventListener) *EventListenerUpdateOne {
	eluo.mutation.Where(ps...)
	return eluo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (eluo *EventListenerUpdateOne) Select(field string, fields ...string) *EventListenerUpdateOne {
	eluo.fields = append([]string{field}, fields...)
	return eluo
}

// Save executes the query and returns the updated EventListener entity.
func (eluo *EventListenerUpdateOne) Save(ctx context.Context) (*EventListener, error) {
	eluo.defaults()
	return withHooks(ctx, eluo.sqlSave, eluo.mutation, eluo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (eluo *EventListenerUpdateOne) SaveX(ctx context.Context) *EventListener {
	node, err := eluo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (eluo *EventListenerUpdateOne) Exec(ctx context.Context) error {
	_, err := eluo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eluo *EventListenerUpdateOne) ExecX(ctx context.Context) {
	if err := eluo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (eluo *EventListenerUpdateOne) defaults() {
	if _, ok := eluo.mutation.UpdateTime(); !ok {
		v := eventlistener.UpdateDefaultUpdateTime()
		eluo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (eluo *EventListenerUpdateOne) check() error {
	if v, ok := eluo.mutation.DataType(); ok {
		if err := eventlistener.DataTypeValidator(v); err != nil {
			return &ValidationError{Name: "data_type", err: fmt.Errorf(`ent: validator failed for field "EventListener.data_type": %w`, err)}
		}
	}
	return nil
}

func (eluo *EventListenerUpdateOne) sqlSave(ctx context.Context) (_node *EventListener, err error) {
	if err := eluo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(eventlistener.Table, eventlistener.Columns, sqlgraph.NewFieldSpec(eventlistener.FieldID, field.TypeInt))
	id, ok := eluo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "EventListener.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := eluo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, eventlistener.FieldID)
		for _, f := range fields {
			if !eventlistener.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != eventlistener.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := eluo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := eluo.mutation.UpdateTime(); ok {
		_spec.SetField(eventlistener.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := eluo.mutation.WalletAddress(); ok {
		_spec.SetField(eventlistener.FieldWalletAddress, field.TypeString, value)
	}
	if eluo.mutation.WalletAddressCleared() {
		_spec.ClearField(eventlistener.FieldWalletAddress, field.TypeString)
	}
	if value, ok := eluo.mutation.DataType(); ok {
		_spec.SetField(eventlistener.FieldDataType, field.TypeEnum, value)
	}
	if eluo.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   eventlistener.UserTable,
			Columns: []string{eventlistener.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eluo.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   eventlistener.UserTable,
			Columns: []string{eventlistener.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if eluo.mutation.ChainCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   eventlistener.ChainTable,
			Columns: []string{eventlistener.ChainColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chain.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eluo.mutation.ChainIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   eventlistener.ChainTable,
			Columns: []string{eventlistener.ChainColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chain.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if eluo.mutation.EventsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   eventlistener.EventsTable,
			Columns: []string{eventlistener.EventsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(event.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eluo.mutation.RemovedEventsIDs(); len(nodes) > 0 && !eluo.mutation.EventsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   eventlistener.EventsTable,
			Columns: []string{eventlistener.EventsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(event.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eluo.mutation.EventsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   eventlistener.EventsTable,
			Columns: []string{eventlistener.EventsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(event.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &EventListener{config: eluo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, eluo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{eventlistener.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	eluo.mutation.done = true
	return _node, nil
}
