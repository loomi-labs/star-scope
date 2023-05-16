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
	"github.com/loomi-labs/star-scope/ent/chain"
	"github.com/loomi-labs/star-scope/ent/eventlistener"
	"github.com/loomi-labs/star-scope/ent/predicate"
)

// ChainUpdate is the builder for updating Chain entities.
type ChainUpdate struct {
	config
	hooks    []Hook
	mutation *ChainMutation
}

// Where appends a list predicates to the ChainUpdate builder.
func (cu *ChainUpdate) Where(ps ...predicate.Chain) *ChainUpdate {
	cu.mutation.Where(ps...)
	return cu
}

// SetUpdateTime sets the "update_time" field.
func (cu *ChainUpdate) SetUpdateTime(t time.Time) *ChainUpdate {
	cu.mutation.SetUpdateTime(t)
	return cu
}

// SetName sets the "name" field.
func (cu *ChainUpdate) SetName(s string) *ChainUpdate {
	cu.mutation.SetName(s)
	return cu
}

// SetImage sets the "image" field.
func (cu *ChainUpdate) SetImage(s string) *ChainUpdate {
	cu.mutation.SetImage(s)
	return cu
}

// SetIndexingHeight sets the "indexing_height" field.
func (cu *ChainUpdate) SetIndexingHeight(u uint64) *ChainUpdate {
	cu.mutation.ResetIndexingHeight()
	cu.mutation.SetIndexingHeight(u)
	return cu
}

// SetNillableIndexingHeight sets the "indexing_height" field if the given value is not nil.
func (cu *ChainUpdate) SetNillableIndexingHeight(u *uint64) *ChainUpdate {
	if u != nil {
		cu.SetIndexingHeight(*u)
	}
	return cu
}

// AddIndexingHeight adds u to the "indexing_height" field.
func (cu *ChainUpdate) AddIndexingHeight(u int64) *ChainUpdate {
	cu.mutation.AddIndexingHeight(u)
	return cu
}

// SetPath sets the "path" field.
func (cu *ChainUpdate) SetPath(s string) *ChainUpdate {
	cu.mutation.SetPath(s)
	return cu
}

// SetHasCustomIndexer sets the "has_custom_indexer" field.
func (cu *ChainUpdate) SetHasCustomIndexer(b bool) *ChainUpdate {
	cu.mutation.SetHasCustomIndexer(b)
	return cu
}

// SetNillableHasCustomIndexer sets the "has_custom_indexer" field if the given value is not nil.
func (cu *ChainUpdate) SetNillableHasCustomIndexer(b *bool) *ChainUpdate {
	if b != nil {
		cu.SetHasCustomIndexer(*b)
	}
	return cu
}

// SetUnhandledMessageTypes sets the "unhandled_message_types" field.
func (cu *ChainUpdate) SetUnhandledMessageTypes(s string) *ChainUpdate {
	cu.mutation.SetUnhandledMessageTypes(s)
	return cu
}

// SetNillableUnhandledMessageTypes sets the "unhandled_message_types" field if the given value is not nil.
func (cu *ChainUpdate) SetNillableUnhandledMessageTypes(s *string) *ChainUpdate {
	if s != nil {
		cu.SetUnhandledMessageTypes(*s)
	}
	return cu
}

// AddEventListenerIDs adds the "event_listeners" edge to the EventListener entity by IDs.
func (cu *ChainUpdate) AddEventListenerIDs(ids ...int) *ChainUpdate {
	cu.mutation.AddEventListenerIDs(ids...)
	return cu
}

// AddEventListeners adds the "event_listeners" edges to the EventListener entity.
func (cu *ChainUpdate) AddEventListeners(e ...*EventListener) *ChainUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return cu.AddEventListenerIDs(ids...)
}

// Mutation returns the ChainMutation object of the builder.
func (cu *ChainUpdate) Mutation() *ChainMutation {
	return cu.mutation
}

// ClearEventListeners clears all "event_listeners" edges to the EventListener entity.
func (cu *ChainUpdate) ClearEventListeners() *ChainUpdate {
	cu.mutation.ClearEventListeners()
	return cu
}

// RemoveEventListenerIDs removes the "event_listeners" edge to EventListener entities by IDs.
func (cu *ChainUpdate) RemoveEventListenerIDs(ids ...int) *ChainUpdate {
	cu.mutation.RemoveEventListenerIDs(ids...)
	return cu
}

// RemoveEventListeners removes "event_listeners" edges to EventListener entities.
func (cu *ChainUpdate) RemoveEventListeners(e ...*EventListener) *ChainUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return cu.RemoveEventListenerIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (cu *ChainUpdate) Save(ctx context.Context) (int, error) {
	cu.defaults()
	return withHooks(ctx, cu.sqlSave, cu.mutation, cu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cu *ChainUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *ChainUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *ChainUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cu *ChainUpdate) defaults() {
	if _, ok := cu.mutation.UpdateTime(); !ok {
		v := chain.UpdateDefaultUpdateTime()
		cu.mutation.SetUpdateTime(v)
	}
}

func (cu *ChainUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(chain.Table, chain.Columns, sqlgraph.NewFieldSpec(chain.FieldID, field.TypeInt))
	if ps := cu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cu.mutation.UpdateTime(); ok {
		_spec.SetField(chain.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := cu.mutation.Name(); ok {
		_spec.SetField(chain.FieldName, field.TypeString, value)
	}
	if value, ok := cu.mutation.Image(); ok {
		_spec.SetField(chain.FieldImage, field.TypeString, value)
	}
	if value, ok := cu.mutation.IndexingHeight(); ok {
		_spec.SetField(chain.FieldIndexingHeight, field.TypeUint64, value)
	}
	if value, ok := cu.mutation.AddedIndexingHeight(); ok {
		_spec.AddField(chain.FieldIndexingHeight, field.TypeUint64, value)
	}
	if value, ok := cu.mutation.Path(); ok {
		_spec.SetField(chain.FieldPath, field.TypeString, value)
	}
	if value, ok := cu.mutation.HasCustomIndexer(); ok {
		_spec.SetField(chain.FieldHasCustomIndexer, field.TypeBool, value)
	}
	if value, ok := cu.mutation.UnhandledMessageTypes(); ok {
		_spec.SetField(chain.FieldUnhandledMessageTypes, field.TypeString, value)
	}
	if cu.mutation.EventListenersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chain.EventListenersTable,
			Columns: []string{chain.EventListenersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(eventlistener.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.RemovedEventListenersIDs(); len(nodes) > 0 && !cu.mutation.EventListenersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chain.EventListenersTable,
			Columns: []string{chain.EventListenersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(eventlistener.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.EventListenersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chain.EventListenersTable,
			Columns: []string{chain.EventListenersColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, cu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{chain.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	cu.mutation.done = true
	return n, nil
}

// ChainUpdateOne is the builder for updating a single Chain entity.
type ChainUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ChainMutation
}

// SetUpdateTime sets the "update_time" field.
func (cuo *ChainUpdateOne) SetUpdateTime(t time.Time) *ChainUpdateOne {
	cuo.mutation.SetUpdateTime(t)
	return cuo
}

// SetName sets the "name" field.
func (cuo *ChainUpdateOne) SetName(s string) *ChainUpdateOne {
	cuo.mutation.SetName(s)
	return cuo
}

// SetImage sets the "image" field.
func (cuo *ChainUpdateOne) SetImage(s string) *ChainUpdateOne {
	cuo.mutation.SetImage(s)
	return cuo
}

// SetIndexingHeight sets the "indexing_height" field.
func (cuo *ChainUpdateOne) SetIndexingHeight(u uint64) *ChainUpdateOne {
	cuo.mutation.ResetIndexingHeight()
	cuo.mutation.SetIndexingHeight(u)
	return cuo
}

// SetNillableIndexingHeight sets the "indexing_height" field if the given value is not nil.
func (cuo *ChainUpdateOne) SetNillableIndexingHeight(u *uint64) *ChainUpdateOne {
	if u != nil {
		cuo.SetIndexingHeight(*u)
	}
	return cuo
}

// AddIndexingHeight adds u to the "indexing_height" field.
func (cuo *ChainUpdateOne) AddIndexingHeight(u int64) *ChainUpdateOne {
	cuo.mutation.AddIndexingHeight(u)
	return cuo
}

// SetPath sets the "path" field.
func (cuo *ChainUpdateOne) SetPath(s string) *ChainUpdateOne {
	cuo.mutation.SetPath(s)
	return cuo
}

// SetHasCustomIndexer sets the "has_custom_indexer" field.
func (cuo *ChainUpdateOne) SetHasCustomIndexer(b bool) *ChainUpdateOne {
	cuo.mutation.SetHasCustomIndexer(b)
	return cuo
}

// SetNillableHasCustomIndexer sets the "has_custom_indexer" field if the given value is not nil.
func (cuo *ChainUpdateOne) SetNillableHasCustomIndexer(b *bool) *ChainUpdateOne {
	if b != nil {
		cuo.SetHasCustomIndexer(*b)
	}
	return cuo
}

// SetUnhandledMessageTypes sets the "unhandled_message_types" field.
func (cuo *ChainUpdateOne) SetUnhandledMessageTypes(s string) *ChainUpdateOne {
	cuo.mutation.SetUnhandledMessageTypes(s)
	return cuo
}

// SetNillableUnhandledMessageTypes sets the "unhandled_message_types" field if the given value is not nil.
func (cuo *ChainUpdateOne) SetNillableUnhandledMessageTypes(s *string) *ChainUpdateOne {
	if s != nil {
		cuo.SetUnhandledMessageTypes(*s)
	}
	return cuo
}

// AddEventListenerIDs adds the "event_listeners" edge to the EventListener entity by IDs.
func (cuo *ChainUpdateOne) AddEventListenerIDs(ids ...int) *ChainUpdateOne {
	cuo.mutation.AddEventListenerIDs(ids...)
	return cuo
}

// AddEventListeners adds the "event_listeners" edges to the EventListener entity.
func (cuo *ChainUpdateOne) AddEventListeners(e ...*EventListener) *ChainUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return cuo.AddEventListenerIDs(ids...)
}

// Mutation returns the ChainMutation object of the builder.
func (cuo *ChainUpdateOne) Mutation() *ChainMutation {
	return cuo.mutation
}

// ClearEventListeners clears all "event_listeners" edges to the EventListener entity.
func (cuo *ChainUpdateOne) ClearEventListeners() *ChainUpdateOne {
	cuo.mutation.ClearEventListeners()
	return cuo
}

// RemoveEventListenerIDs removes the "event_listeners" edge to EventListener entities by IDs.
func (cuo *ChainUpdateOne) RemoveEventListenerIDs(ids ...int) *ChainUpdateOne {
	cuo.mutation.RemoveEventListenerIDs(ids...)
	return cuo
}

// RemoveEventListeners removes "event_listeners" edges to EventListener entities.
func (cuo *ChainUpdateOne) RemoveEventListeners(e ...*EventListener) *ChainUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return cuo.RemoveEventListenerIDs(ids...)
}

// Where appends a list predicates to the ChainUpdate builder.
func (cuo *ChainUpdateOne) Where(ps ...predicate.Chain) *ChainUpdateOne {
	cuo.mutation.Where(ps...)
	return cuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (cuo *ChainUpdateOne) Select(field string, fields ...string) *ChainUpdateOne {
	cuo.fields = append([]string{field}, fields...)
	return cuo
}

// Save executes the query and returns the updated Chain entity.
func (cuo *ChainUpdateOne) Save(ctx context.Context) (*Chain, error) {
	cuo.defaults()
	return withHooks(ctx, cuo.sqlSave, cuo.mutation, cuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *ChainUpdateOne) SaveX(ctx context.Context) *Chain {
	node, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (cuo *ChainUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *ChainUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cuo *ChainUpdateOne) defaults() {
	if _, ok := cuo.mutation.UpdateTime(); !ok {
		v := chain.UpdateDefaultUpdateTime()
		cuo.mutation.SetUpdateTime(v)
	}
}

func (cuo *ChainUpdateOne) sqlSave(ctx context.Context) (_node *Chain, err error) {
	_spec := sqlgraph.NewUpdateSpec(chain.Table, chain.Columns, sqlgraph.NewFieldSpec(chain.FieldID, field.TypeInt))
	id, ok := cuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Chain.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := cuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, chain.FieldID)
		for _, f := range fields {
			if !chain.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != chain.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := cuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cuo.mutation.UpdateTime(); ok {
		_spec.SetField(chain.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := cuo.mutation.Name(); ok {
		_spec.SetField(chain.FieldName, field.TypeString, value)
	}
	if value, ok := cuo.mutation.Image(); ok {
		_spec.SetField(chain.FieldImage, field.TypeString, value)
	}
	if value, ok := cuo.mutation.IndexingHeight(); ok {
		_spec.SetField(chain.FieldIndexingHeight, field.TypeUint64, value)
	}
	if value, ok := cuo.mutation.AddedIndexingHeight(); ok {
		_spec.AddField(chain.FieldIndexingHeight, field.TypeUint64, value)
	}
	if value, ok := cuo.mutation.Path(); ok {
		_spec.SetField(chain.FieldPath, field.TypeString, value)
	}
	if value, ok := cuo.mutation.HasCustomIndexer(); ok {
		_spec.SetField(chain.FieldHasCustomIndexer, field.TypeBool, value)
	}
	if value, ok := cuo.mutation.UnhandledMessageTypes(); ok {
		_spec.SetField(chain.FieldUnhandledMessageTypes, field.TypeString, value)
	}
	if cuo.mutation.EventListenersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chain.EventListenersTable,
			Columns: []string{chain.EventListenersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(eventlistener.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.RemovedEventListenersIDs(); len(nodes) > 0 && !cuo.mutation.EventListenersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chain.EventListenersTable,
			Columns: []string{chain.EventListenersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(eventlistener.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.EventListenersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chain.EventListenersTable,
			Columns: []string{chain.EventListenersColumn},
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
	_node = &Chain{config: cuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, cuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{chain.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	cuo.mutation.done = true
	return _node, nil
}
