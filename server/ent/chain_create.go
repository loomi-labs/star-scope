// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/loomi-labs/star-scope/ent/chain"
	"github.com/loomi-labs/star-scope/ent/contractproposal"
	"github.com/loomi-labs/star-scope/ent/eventlistener"
	"github.com/loomi-labs/star-scope/ent/proposal"
)

// ChainCreate is the builder for creating a Chain entity.
type ChainCreate struct {
	config
	mutation *ChainMutation
	hooks    []Hook
}

// SetCreateTime sets the "create_time" field.
func (cc *ChainCreate) SetCreateTime(t time.Time) *ChainCreate {
	cc.mutation.SetCreateTime(t)
	return cc
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (cc *ChainCreate) SetNillableCreateTime(t *time.Time) *ChainCreate {
	if t != nil {
		cc.SetCreateTime(*t)
	}
	return cc
}

// SetUpdateTime sets the "update_time" field.
func (cc *ChainCreate) SetUpdateTime(t time.Time) *ChainCreate {
	cc.mutation.SetUpdateTime(t)
	return cc
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (cc *ChainCreate) SetNillableUpdateTime(t *time.Time) *ChainCreate {
	if t != nil {
		cc.SetUpdateTime(*t)
	}
	return cc
}

// SetChainID sets the "chain_id" field.
func (cc *ChainCreate) SetChainID(s string) *ChainCreate {
	cc.mutation.SetChainID(s)
	return cc
}

// SetName sets the "name" field.
func (cc *ChainCreate) SetName(s string) *ChainCreate {
	cc.mutation.SetName(s)
	return cc
}

// SetPrettyName sets the "pretty_name" field.
func (cc *ChainCreate) SetPrettyName(s string) *ChainCreate {
	cc.mutation.SetPrettyName(s)
	return cc
}

// SetPath sets the "path" field.
func (cc *ChainCreate) SetPath(s string) *ChainCreate {
	cc.mutation.SetPath(s)
	return cc
}

// SetImage sets the "image" field.
func (cc *ChainCreate) SetImage(s string) *ChainCreate {
	cc.mutation.SetImage(s)
	return cc
}

// SetBech32Prefix sets the "bech32_prefix" field.
func (cc *ChainCreate) SetBech32Prefix(s string) *ChainCreate {
	cc.mutation.SetBech32Prefix(s)
	return cc
}

// SetIndexingHeight sets the "indexing_height" field.
func (cc *ChainCreate) SetIndexingHeight(u uint64) *ChainCreate {
	cc.mutation.SetIndexingHeight(u)
	return cc
}

// SetNillableIndexingHeight sets the "indexing_height" field if the given value is not nil.
func (cc *ChainCreate) SetNillableIndexingHeight(u *uint64) *ChainCreate {
	if u != nil {
		cc.SetIndexingHeight(*u)
	}
	return cc
}

// SetHasCustomIndexer sets the "has_custom_indexer" field.
func (cc *ChainCreate) SetHasCustomIndexer(b bool) *ChainCreate {
	cc.mutation.SetHasCustomIndexer(b)
	return cc
}

// SetNillableHasCustomIndexer sets the "has_custom_indexer" field if the given value is not nil.
func (cc *ChainCreate) SetNillableHasCustomIndexer(b *bool) *ChainCreate {
	if b != nil {
		cc.SetHasCustomIndexer(*b)
	}
	return cc
}

// SetHandledMessageTypes sets the "handled_message_types" field.
func (cc *ChainCreate) SetHandledMessageTypes(s string) *ChainCreate {
	cc.mutation.SetHandledMessageTypes(s)
	return cc
}

// SetNillableHandledMessageTypes sets the "handled_message_types" field if the given value is not nil.
func (cc *ChainCreate) SetNillableHandledMessageTypes(s *string) *ChainCreate {
	if s != nil {
		cc.SetHandledMessageTypes(*s)
	}
	return cc
}

// SetUnhandledMessageTypes sets the "unhandled_message_types" field.
func (cc *ChainCreate) SetUnhandledMessageTypes(s string) *ChainCreate {
	cc.mutation.SetUnhandledMessageTypes(s)
	return cc
}

// SetNillableUnhandledMessageTypes sets the "unhandled_message_types" field if the given value is not nil.
func (cc *ChainCreate) SetNillableUnhandledMessageTypes(s *string) *ChainCreate {
	if s != nil {
		cc.SetUnhandledMessageTypes(*s)
	}
	return cc
}

// SetIsEnabled sets the "is_enabled" field.
func (cc *ChainCreate) SetIsEnabled(b bool) *ChainCreate {
	cc.mutation.SetIsEnabled(b)
	return cc
}

// SetNillableIsEnabled sets the "is_enabled" field if the given value is not nil.
func (cc *ChainCreate) SetNillableIsEnabled(b *bool) *ChainCreate {
	if b != nil {
		cc.SetIsEnabled(*b)
	}
	return cc
}

// AddEventListenerIDs adds the "event_listeners" edge to the EventListener entity by IDs.
func (cc *ChainCreate) AddEventListenerIDs(ids ...int) *ChainCreate {
	cc.mutation.AddEventListenerIDs(ids...)
	return cc
}

// AddEventListeners adds the "event_listeners" edges to the EventListener entity.
func (cc *ChainCreate) AddEventListeners(e ...*EventListener) *ChainCreate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return cc.AddEventListenerIDs(ids...)
}

// AddProposalIDs adds the "proposals" edge to the Proposal entity by IDs.
func (cc *ChainCreate) AddProposalIDs(ids ...int) *ChainCreate {
	cc.mutation.AddProposalIDs(ids...)
	return cc
}

// AddProposals adds the "proposals" edges to the Proposal entity.
func (cc *ChainCreate) AddProposals(p ...*Proposal) *ChainCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return cc.AddProposalIDs(ids...)
}

// AddContractProposalIDs adds the "contract_proposals" edge to the ContractProposal entity by IDs.
func (cc *ChainCreate) AddContractProposalIDs(ids ...int) *ChainCreate {
	cc.mutation.AddContractProposalIDs(ids...)
	return cc
}

// AddContractProposals adds the "contract_proposals" edges to the ContractProposal entity.
func (cc *ChainCreate) AddContractProposals(c ...*ContractProposal) *ChainCreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cc.AddContractProposalIDs(ids...)
}

// Mutation returns the ChainMutation object of the builder.
func (cc *ChainCreate) Mutation() *ChainMutation {
	return cc.mutation
}

// Save creates the Chain in the database.
func (cc *ChainCreate) Save(ctx context.Context) (*Chain, error) {
	cc.defaults()
	return withHooks(ctx, cc.sqlSave, cc.mutation, cc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (cc *ChainCreate) SaveX(ctx context.Context) *Chain {
	v, err := cc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cc *ChainCreate) Exec(ctx context.Context) error {
	_, err := cc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cc *ChainCreate) ExecX(ctx context.Context) {
	if err := cc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cc *ChainCreate) defaults() {
	if _, ok := cc.mutation.CreateTime(); !ok {
		v := chain.DefaultCreateTime()
		cc.mutation.SetCreateTime(v)
	}
	if _, ok := cc.mutation.UpdateTime(); !ok {
		v := chain.DefaultUpdateTime()
		cc.mutation.SetUpdateTime(v)
	}
	if _, ok := cc.mutation.IndexingHeight(); !ok {
		v := chain.DefaultIndexingHeight
		cc.mutation.SetIndexingHeight(v)
	}
	if _, ok := cc.mutation.HasCustomIndexer(); !ok {
		v := chain.DefaultHasCustomIndexer
		cc.mutation.SetHasCustomIndexer(v)
	}
	if _, ok := cc.mutation.HandledMessageTypes(); !ok {
		v := chain.DefaultHandledMessageTypes
		cc.mutation.SetHandledMessageTypes(v)
	}
	if _, ok := cc.mutation.UnhandledMessageTypes(); !ok {
		v := chain.DefaultUnhandledMessageTypes
		cc.mutation.SetUnhandledMessageTypes(v)
	}
	if _, ok := cc.mutation.IsEnabled(); !ok {
		v := chain.DefaultIsEnabled
		cc.mutation.SetIsEnabled(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cc *ChainCreate) check() error {
	if _, ok := cc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "create_time", err: errors.New(`ent: missing required field "Chain.create_time"`)}
	}
	if _, ok := cc.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "update_time", err: errors.New(`ent: missing required field "Chain.update_time"`)}
	}
	if _, ok := cc.mutation.ChainID(); !ok {
		return &ValidationError{Name: "chain_id", err: errors.New(`ent: missing required field "Chain.chain_id"`)}
	}
	if _, ok := cc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Chain.name"`)}
	}
	if _, ok := cc.mutation.PrettyName(); !ok {
		return &ValidationError{Name: "pretty_name", err: errors.New(`ent: missing required field "Chain.pretty_name"`)}
	}
	if _, ok := cc.mutation.Path(); !ok {
		return &ValidationError{Name: "path", err: errors.New(`ent: missing required field "Chain.path"`)}
	}
	if _, ok := cc.mutation.Image(); !ok {
		return &ValidationError{Name: "image", err: errors.New(`ent: missing required field "Chain.image"`)}
	}
	if _, ok := cc.mutation.Bech32Prefix(); !ok {
		return &ValidationError{Name: "bech32_prefix", err: errors.New(`ent: missing required field "Chain.bech32_prefix"`)}
	}
	if _, ok := cc.mutation.IndexingHeight(); !ok {
		return &ValidationError{Name: "indexing_height", err: errors.New(`ent: missing required field "Chain.indexing_height"`)}
	}
	if _, ok := cc.mutation.HasCustomIndexer(); !ok {
		return &ValidationError{Name: "has_custom_indexer", err: errors.New(`ent: missing required field "Chain.has_custom_indexer"`)}
	}
	if _, ok := cc.mutation.HandledMessageTypes(); !ok {
		return &ValidationError{Name: "handled_message_types", err: errors.New(`ent: missing required field "Chain.handled_message_types"`)}
	}
	if _, ok := cc.mutation.UnhandledMessageTypes(); !ok {
		return &ValidationError{Name: "unhandled_message_types", err: errors.New(`ent: missing required field "Chain.unhandled_message_types"`)}
	}
	if _, ok := cc.mutation.IsEnabled(); !ok {
		return &ValidationError{Name: "is_enabled", err: errors.New(`ent: missing required field "Chain.is_enabled"`)}
	}
	return nil
}

func (cc *ChainCreate) sqlSave(ctx context.Context) (*Chain, error) {
	if err := cc.check(); err != nil {
		return nil, err
	}
	_node, _spec := cc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	cc.mutation.id = &_node.ID
	cc.mutation.done = true
	return _node, nil
}

func (cc *ChainCreate) createSpec() (*Chain, *sqlgraph.CreateSpec) {
	var (
		_node = &Chain{config: cc.config}
		_spec = sqlgraph.NewCreateSpec(chain.Table, sqlgraph.NewFieldSpec(chain.FieldID, field.TypeInt))
	)
	if value, ok := cc.mutation.CreateTime(); ok {
		_spec.SetField(chain.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = value
	}
	if value, ok := cc.mutation.UpdateTime(); ok {
		_spec.SetField(chain.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = value
	}
	if value, ok := cc.mutation.ChainID(); ok {
		_spec.SetField(chain.FieldChainID, field.TypeString, value)
		_node.ChainID = value
	}
	if value, ok := cc.mutation.Name(); ok {
		_spec.SetField(chain.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := cc.mutation.PrettyName(); ok {
		_spec.SetField(chain.FieldPrettyName, field.TypeString, value)
		_node.PrettyName = value
	}
	if value, ok := cc.mutation.Path(); ok {
		_spec.SetField(chain.FieldPath, field.TypeString, value)
		_node.Path = value
	}
	if value, ok := cc.mutation.Image(); ok {
		_spec.SetField(chain.FieldImage, field.TypeString, value)
		_node.Image = value
	}
	if value, ok := cc.mutation.Bech32Prefix(); ok {
		_spec.SetField(chain.FieldBech32Prefix, field.TypeString, value)
		_node.Bech32Prefix = value
	}
	if value, ok := cc.mutation.IndexingHeight(); ok {
		_spec.SetField(chain.FieldIndexingHeight, field.TypeUint64, value)
		_node.IndexingHeight = value
	}
	if value, ok := cc.mutation.HasCustomIndexer(); ok {
		_spec.SetField(chain.FieldHasCustomIndexer, field.TypeBool, value)
		_node.HasCustomIndexer = value
	}
	if value, ok := cc.mutation.HandledMessageTypes(); ok {
		_spec.SetField(chain.FieldHandledMessageTypes, field.TypeString, value)
		_node.HandledMessageTypes = value
	}
	if value, ok := cc.mutation.UnhandledMessageTypes(); ok {
		_spec.SetField(chain.FieldUnhandledMessageTypes, field.TypeString, value)
		_node.UnhandledMessageTypes = value
	}
	if value, ok := cc.mutation.IsEnabled(); ok {
		_spec.SetField(chain.FieldIsEnabled, field.TypeBool, value)
		_node.IsEnabled = value
	}
	if nodes := cc.mutation.EventListenersIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cc.mutation.ProposalsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chain.ProposalsTable,
			Columns: []string{chain.ProposalsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(proposal.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cc.mutation.ContractProposalsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chain.ContractProposalsTable,
			Columns: []string{chain.ContractProposalsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(contractproposal.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// ChainCreateBulk is the builder for creating many Chain entities in bulk.
type ChainCreateBulk struct {
	config
	builders []*ChainCreate
}

// Save creates the Chain entities in the database.
func (ccb *ChainCreateBulk) Save(ctx context.Context) ([]*Chain, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ccb.builders))
	nodes := make([]*Chain, len(ccb.builders))
	mutators := make([]Mutator, len(ccb.builders))
	for i := range ccb.builders {
		func(i int, root context.Context) {
			builder := ccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ChainMutation)
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
					_, err = mutators[i+1].Mutate(root, ccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ccb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, ccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ccb *ChainCreateBulk) SaveX(ctx context.Context) []*Chain {
	v, err := ccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ccb *ChainCreateBulk) Exec(ctx context.Context) error {
	_, err := ccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ccb *ChainCreateBulk) ExecX(ctx context.Context) {
	if err := ccb.Exec(ctx); err != nil {
		panic(err)
	}
}
