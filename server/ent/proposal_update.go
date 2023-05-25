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
	"github.com/loomi-labs/star-scope/ent/predicate"
	"github.com/loomi-labs/star-scope/ent/proposal"
)

// ProposalUpdate is the builder for updating Proposal entities.
type ProposalUpdate struct {
	config
	hooks    []Hook
	mutation *ProposalMutation
}

// Where appends a list predicates to the ProposalUpdate builder.
func (pu *ProposalUpdate) Where(ps ...predicate.Proposal) *ProposalUpdate {
	pu.mutation.Where(ps...)
	return pu
}

// SetUpdateTime sets the "update_time" field.
func (pu *ProposalUpdate) SetUpdateTime(t time.Time) *ProposalUpdate {
	pu.mutation.SetUpdateTime(t)
	return pu
}

// SetProposalID sets the "proposal_id" field.
func (pu *ProposalUpdate) SetProposalID(u uint64) *ProposalUpdate {
	pu.mutation.ResetProposalID()
	pu.mutation.SetProposalID(u)
	return pu
}

// AddProposalID adds u to the "proposal_id" field.
func (pu *ProposalUpdate) AddProposalID(u int64) *ProposalUpdate {
	pu.mutation.AddProposalID(u)
	return pu
}

// SetTitle sets the "title" field.
func (pu *ProposalUpdate) SetTitle(s string) *ProposalUpdate {
	pu.mutation.SetTitle(s)
	return pu
}

// SetDescription sets the "description" field.
func (pu *ProposalUpdate) SetDescription(s string) *ProposalUpdate {
	pu.mutation.SetDescription(s)
	return pu
}

// SetVotingStartTime sets the "voting_start_time" field.
func (pu *ProposalUpdate) SetVotingStartTime(t time.Time) *ProposalUpdate {
	pu.mutation.SetVotingStartTime(t)
	return pu
}

// SetVotingEndTime sets the "voting_end_time" field.
func (pu *ProposalUpdate) SetVotingEndTime(t time.Time) *ProposalUpdate {
	pu.mutation.SetVotingEndTime(t)
	return pu
}

// SetStatus sets the "status" field.
func (pu *ProposalUpdate) SetStatus(pr proposal.Status) *ProposalUpdate {
	pu.mutation.SetStatus(pr)
	return pu
}

// SetChainID sets the "chain" edge to the Chain entity by ID.
func (pu *ProposalUpdate) SetChainID(id int) *ProposalUpdate {
	pu.mutation.SetChainID(id)
	return pu
}

// SetNillableChainID sets the "chain" edge to the Chain entity by ID if the given value is not nil.
func (pu *ProposalUpdate) SetNillableChainID(id *int) *ProposalUpdate {
	if id != nil {
		pu = pu.SetChainID(*id)
	}
	return pu
}

// SetChain sets the "chain" edge to the Chain entity.
func (pu *ProposalUpdate) SetChain(c *Chain) *ProposalUpdate {
	return pu.SetChainID(c.ID)
}

// Mutation returns the ProposalMutation object of the builder.
func (pu *ProposalUpdate) Mutation() *ProposalMutation {
	return pu.mutation
}

// ClearChain clears the "chain" edge to the Chain entity.
func (pu *ProposalUpdate) ClearChain() *ProposalUpdate {
	pu.mutation.ClearChain()
	return pu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pu *ProposalUpdate) Save(ctx context.Context) (int, error) {
	pu.defaults()
	return withHooks(ctx, pu.sqlSave, pu.mutation, pu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pu *ProposalUpdate) SaveX(ctx context.Context) int {
	affected, err := pu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pu *ProposalUpdate) Exec(ctx context.Context) error {
	_, err := pu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pu *ProposalUpdate) ExecX(ctx context.Context) {
	if err := pu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pu *ProposalUpdate) defaults() {
	if _, ok := pu.mutation.UpdateTime(); !ok {
		v := proposal.UpdateDefaultUpdateTime()
		pu.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pu *ProposalUpdate) check() error {
	if v, ok := pu.mutation.Status(); ok {
		if err := proposal.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Proposal.status": %w`, err)}
		}
	}
	return nil
}

func (pu *ProposalUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := pu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(proposal.Table, proposal.Columns, sqlgraph.NewFieldSpec(proposal.FieldID, field.TypeInt))
	if ps := pu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pu.mutation.UpdateTime(); ok {
		_spec.SetField(proposal.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := pu.mutation.ProposalID(); ok {
		_spec.SetField(proposal.FieldProposalID, field.TypeUint64, value)
	}
	if value, ok := pu.mutation.AddedProposalID(); ok {
		_spec.AddField(proposal.FieldProposalID, field.TypeUint64, value)
	}
	if value, ok := pu.mutation.Title(); ok {
		_spec.SetField(proposal.FieldTitle, field.TypeString, value)
	}
	if value, ok := pu.mutation.Description(); ok {
		_spec.SetField(proposal.FieldDescription, field.TypeString, value)
	}
	if value, ok := pu.mutation.VotingStartTime(); ok {
		_spec.SetField(proposal.FieldVotingStartTime, field.TypeTime, value)
	}
	if value, ok := pu.mutation.VotingEndTime(); ok {
		_spec.SetField(proposal.FieldVotingEndTime, field.TypeTime, value)
	}
	if value, ok := pu.mutation.Status(); ok {
		_spec.SetField(proposal.FieldStatus, field.TypeEnum, value)
	}
	if pu.mutation.ChainCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   proposal.ChainTable,
			Columns: []string{proposal.ChainColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chain.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.ChainIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   proposal.ChainTable,
			Columns: []string{proposal.ChainColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, pu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{proposal.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	pu.mutation.done = true
	return n, nil
}

// ProposalUpdateOne is the builder for updating a single Proposal entity.
type ProposalUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ProposalMutation
}

// SetUpdateTime sets the "update_time" field.
func (puo *ProposalUpdateOne) SetUpdateTime(t time.Time) *ProposalUpdateOne {
	puo.mutation.SetUpdateTime(t)
	return puo
}

// SetProposalID sets the "proposal_id" field.
func (puo *ProposalUpdateOne) SetProposalID(u uint64) *ProposalUpdateOne {
	puo.mutation.ResetProposalID()
	puo.mutation.SetProposalID(u)
	return puo
}

// AddProposalID adds u to the "proposal_id" field.
func (puo *ProposalUpdateOne) AddProposalID(u int64) *ProposalUpdateOne {
	puo.mutation.AddProposalID(u)
	return puo
}

// SetTitle sets the "title" field.
func (puo *ProposalUpdateOne) SetTitle(s string) *ProposalUpdateOne {
	puo.mutation.SetTitle(s)
	return puo
}

// SetDescription sets the "description" field.
func (puo *ProposalUpdateOne) SetDescription(s string) *ProposalUpdateOne {
	puo.mutation.SetDescription(s)
	return puo
}

// SetVotingStartTime sets the "voting_start_time" field.
func (puo *ProposalUpdateOne) SetVotingStartTime(t time.Time) *ProposalUpdateOne {
	puo.mutation.SetVotingStartTime(t)
	return puo
}

// SetVotingEndTime sets the "voting_end_time" field.
func (puo *ProposalUpdateOne) SetVotingEndTime(t time.Time) *ProposalUpdateOne {
	puo.mutation.SetVotingEndTime(t)
	return puo
}

// SetStatus sets the "status" field.
func (puo *ProposalUpdateOne) SetStatus(pr proposal.Status) *ProposalUpdateOne {
	puo.mutation.SetStatus(pr)
	return puo
}

// SetChainID sets the "chain" edge to the Chain entity by ID.
func (puo *ProposalUpdateOne) SetChainID(id int) *ProposalUpdateOne {
	puo.mutation.SetChainID(id)
	return puo
}

// SetNillableChainID sets the "chain" edge to the Chain entity by ID if the given value is not nil.
func (puo *ProposalUpdateOne) SetNillableChainID(id *int) *ProposalUpdateOne {
	if id != nil {
		puo = puo.SetChainID(*id)
	}
	return puo
}

// SetChain sets the "chain" edge to the Chain entity.
func (puo *ProposalUpdateOne) SetChain(c *Chain) *ProposalUpdateOne {
	return puo.SetChainID(c.ID)
}

// Mutation returns the ProposalMutation object of the builder.
func (puo *ProposalUpdateOne) Mutation() *ProposalMutation {
	return puo.mutation
}

// ClearChain clears the "chain" edge to the Chain entity.
func (puo *ProposalUpdateOne) ClearChain() *ProposalUpdateOne {
	puo.mutation.ClearChain()
	return puo
}

// Where appends a list predicates to the ProposalUpdate builder.
func (puo *ProposalUpdateOne) Where(ps ...predicate.Proposal) *ProposalUpdateOne {
	puo.mutation.Where(ps...)
	return puo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (puo *ProposalUpdateOne) Select(field string, fields ...string) *ProposalUpdateOne {
	puo.fields = append([]string{field}, fields...)
	return puo
}

// Save executes the query and returns the updated Proposal entity.
func (puo *ProposalUpdateOne) Save(ctx context.Context) (*Proposal, error) {
	puo.defaults()
	return withHooks(ctx, puo.sqlSave, puo.mutation, puo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (puo *ProposalUpdateOne) SaveX(ctx context.Context) *Proposal {
	node, err := puo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (puo *ProposalUpdateOne) Exec(ctx context.Context) error {
	_, err := puo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (puo *ProposalUpdateOne) ExecX(ctx context.Context) {
	if err := puo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (puo *ProposalUpdateOne) defaults() {
	if _, ok := puo.mutation.UpdateTime(); !ok {
		v := proposal.UpdateDefaultUpdateTime()
		puo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (puo *ProposalUpdateOne) check() error {
	if v, ok := puo.mutation.Status(); ok {
		if err := proposal.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Proposal.status": %w`, err)}
		}
	}
	return nil
}

func (puo *ProposalUpdateOne) sqlSave(ctx context.Context) (_node *Proposal, err error) {
	if err := puo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(proposal.Table, proposal.Columns, sqlgraph.NewFieldSpec(proposal.FieldID, field.TypeInt))
	id, ok := puo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Proposal.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := puo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, proposal.FieldID)
		for _, f := range fields {
			if !proposal.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != proposal.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := puo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := puo.mutation.UpdateTime(); ok {
		_spec.SetField(proposal.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := puo.mutation.ProposalID(); ok {
		_spec.SetField(proposal.FieldProposalID, field.TypeUint64, value)
	}
	if value, ok := puo.mutation.AddedProposalID(); ok {
		_spec.AddField(proposal.FieldProposalID, field.TypeUint64, value)
	}
	if value, ok := puo.mutation.Title(); ok {
		_spec.SetField(proposal.FieldTitle, field.TypeString, value)
	}
	if value, ok := puo.mutation.Description(); ok {
		_spec.SetField(proposal.FieldDescription, field.TypeString, value)
	}
	if value, ok := puo.mutation.VotingStartTime(); ok {
		_spec.SetField(proposal.FieldVotingStartTime, field.TypeTime, value)
	}
	if value, ok := puo.mutation.VotingEndTime(); ok {
		_spec.SetField(proposal.FieldVotingEndTime, field.TypeTime, value)
	}
	if value, ok := puo.mutation.Status(); ok {
		_spec.SetField(proposal.FieldStatus, field.TypeEnum, value)
	}
	if puo.mutation.ChainCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   proposal.ChainTable,
			Columns: []string{proposal.ChainColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chain.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.ChainIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   proposal.ChainTable,
			Columns: []string{proposal.ChainColumn},
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
	_node = &Proposal{config: puo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, puo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{proposal.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	puo.mutation.done = true
	return _node, nil
}