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
	"github.com/loomi-labs/star-scope/ent/contractproposal"
	"github.com/loomi-labs/star-scope/ent/predicate"
)

// ContractProposalUpdate is the builder for updating ContractProposal entities.
type ContractProposalUpdate struct {
	config
	hooks    []Hook
	mutation *ContractProposalMutation
}

// Where appends a list predicates to the ContractProposalUpdate builder.
func (cpu *ContractProposalUpdate) Where(ps ...predicate.ContractProposal) *ContractProposalUpdate {
	cpu.mutation.Where(ps...)
	return cpu
}

// SetUpdateTime sets the "update_time" field.
func (cpu *ContractProposalUpdate) SetUpdateTime(t time.Time) *ContractProposalUpdate {
	cpu.mutation.SetUpdateTime(t)
	return cpu
}

// SetProposalID sets the "proposal_id" field.
func (cpu *ContractProposalUpdate) SetProposalID(u uint64) *ContractProposalUpdate {
	cpu.mutation.ResetProposalID()
	cpu.mutation.SetProposalID(u)
	return cpu
}

// AddProposalID adds u to the "proposal_id" field.
func (cpu *ContractProposalUpdate) AddProposalID(u int64) *ContractProposalUpdate {
	cpu.mutation.AddProposalID(u)
	return cpu
}

// SetTitle sets the "title" field.
func (cpu *ContractProposalUpdate) SetTitle(s string) *ContractProposalUpdate {
	cpu.mutation.SetTitle(s)
	return cpu
}

// SetDescription sets the "description" field.
func (cpu *ContractProposalUpdate) SetDescription(s string) *ContractProposalUpdate {
	cpu.mutation.SetDescription(s)
	return cpu
}

// SetFirstSeenTime sets the "first_seen_time" field.
func (cpu *ContractProposalUpdate) SetFirstSeenTime(t time.Time) *ContractProposalUpdate {
	cpu.mutation.SetFirstSeenTime(t)
	return cpu
}

// SetVotingEndTime sets the "voting_end_time" field.
func (cpu *ContractProposalUpdate) SetVotingEndTime(t time.Time) *ContractProposalUpdate {
	cpu.mutation.SetVotingEndTime(t)
	return cpu
}

// SetContractAddress sets the "contract_address" field.
func (cpu *ContractProposalUpdate) SetContractAddress(s string) *ContractProposalUpdate {
	cpu.mutation.SetContractAddress(s)
	return cpu
}

// SetStatus sets the "status" field.
func (cpu *ContractProposalUpdate) SetStatus(c contractproposal.Status) *ContractProposalUpdate {
	cpu.mutation.SetStatus(c)
	return cpu
}

// SetChainID sets the "chain" edge to the Chain entity by ID.
func (cpu *ContractProposalUpdate) SetChainID(id int) *ContractProposalUpdate {
	cpu.mutation.SetChainID(id)
	return cpu
}

// SetNillableChainID sets the "chain" edge to the Chain entity by ID if the given value is not nil.
func (cpu *ContractProposalUpdate) SetNillableChainID(id *int) *ContractProposalUpdate {
	if id != nil {
		cpu = cpu.SetChainID(*id)
	}
	return cpu
}

// SetChain sets the "chain" edge to the Chain entity.
func (cpu *ContractProposalUpdate) SetChain(c *Chain) *ContractProposalUpdate {
	return cpu.SetChainID(c.ID)
}

// Mutation returns the ContractProposalMutation object of the builder.
func (cpu *ContractProposalUpdate) Mutation() *ContractProposalMutation {
	return cpu.mutation
}

// ClearChain clears the "chain" edge to the Chain entity.
func (cpu *ContractProposalUpdate) ClearChain() *ContractProposalUpdate {
	cpu.mutation.ClearChain()
	return cpu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (cpu *ContractProposalUpdate) Save(ctx context.Context) (int, error) {
	cpu.defaults()
	return withHooks(ctx, cpu.sqlSave, cpu.mutation, cpu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cpu *ContractProposalUpdate) SaveX(ctx context.Context) int {
	affected, err := cpu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cpu *ContractProposalUpdate) Exec(ctx context.Context) error {
	_, err := cpu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cpu *ContractProposalUpdate) ExecX(ctx context.Context) {
	if err := cpu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cpu *ContractProposalUpdate) defaults() {
	if _, ok := cpu.mutation.UpdateTime(); !ok {
		v := contractproposal.UpdateDefaultUpdateTime()
		cpu.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cpu *ContractProposalUpdate) check() error {
	if v, ok := cpu.mutation.Status(); ok {
		if err := contractproposal.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "ContractProposal.status": %w`, err)}
		}
	}
	return nil
}

func (cpu *ContractProposalUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := cpu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(contractproposal.Table, contractproposal.Columns, sqlgraph.NewFieldSpec(contractproposal.FieldID, field.TypeInt))
	if ps := cpu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cpu.mutation.UpdateTime(); ok {
		_spec.SetField(contractproposal.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := cpu.mutation.ProposalID(); ok {
		_spec.SetField(contractproposal.FieldProposalID, field.TypeUint64, value)
	}
	if value, ok := cpu.mutation.AddedProposalID(); ok {
		_spec.AddField(contractproposal.FieldProposalID, field.TypeUint64, value)
	}
	if value, ok := cpu.mutation.Title(); ok {
		_spec.SetField(contractproposal.FieldTitle, field.TypeString, value)
	}
	if value, ok := cpu.mutation.Description(); ok {
		_spec.SetField(contractproposal.FieldDescription, field.TypeString, value)
	}
	if value, ok := cpu.mutation.FirstSeenTime(); ok {
		_spec.SetField(contractproposal.FieldFirstSeenTime, field.TypeTime, value)
	}
	if value, ok := cpu.mutation.VotingEndTime(); ok {
		_spec.SetField(contractproposal.FieldVotingEndTime, field.TypeTime, value)
	}
	if value, ok := cpu.mutation.ContractAddress(); ok {
		_spec.SetField(contractproposal.FieldContractAddress, field.TypeString, value)
	}
	if value, ok := cpu.mutation.Status(); ok {
		_spec.SetField(contractproposal.FieldStatus, field.TypeEnum, value)
	}
	if cpu.mutation.ChainCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   contractproposal.ChainTable,
			Columns: []string{contractproposal.ChainColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chain.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cpu.mutation.ChainIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   contractproposal.ChainTable,
			Columns: []string{contractproposal.ChainColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, cpu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{contractproposal.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	cpu.mutation.done = true
	return n, nil
}

// ContractProposalUpdateOne is the builder for updating a single ContractProposal entity.
type ContractProposalUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ContractProposalMutation
}

// SetUpdateTime sets the "update_time" field.
func (cpuo *ContractProposalUpdateOne) SetUpdateTime(t time.Time) *ContractProposalUpdateOne {
	cpuo.mutation.SetUpdateTime(t)
	return cpuo
}

// SetProposalID sets the "proposal_id" field.
func (cpuo *ContractProposalUpdateOne) SetProposalID(u uint64) *ContractProposalUpdateOne {
	cpuo.mutation.ResetProposalID()
	cpuo.mutation.SetProposalID(u)
	return cpuo
}

// AddProposalID adds u to the "proposal_id" field.
func (cpuo *ContractProposalUpdateOne) AddProposalID(u int64) *ContractProposalUpdateOne {
	cpuo.mutation.AddProposalID(u)
	return cpuo
}

// SetTitle sets the "title" field.
func (cpuo *ContractProposalUpdateOne) SetTitle(s string) *ContractProposalUpdateOne {
	cpuo.mutation.SetTitle(s)
	return cpuo
}

// SetDescription sets the "description" field.
func (cpuo *ContractProposalUpdateOne) SetDescription(s string) *ContractProposalUpdateOne {
	cpuo.mutation.SetDescription(s)
	return cpuo
}

// SetFirstSeenTime sets the "first_seen_time" field.
func (cpuo *ContractProposalUpdateOne) SetFirstSeenTime(t time.Time) *ContractProposalUpdateOne {
	cpuo.mutation.SetFirstSeenTime(t)
	return cpuo
}

// SetVotingEndTime sets the "voting_end_time" field.
func (cpuo *ContractProposalUpdateOne) SetVotingEndTime(t time.Time) *ContractProposalUpdateOne {
	cpuo.mutation.SetVotingEndTime(t)
	return cpuo
}

// SetContractAddress sets the "contract_address" field.
func (cpuo *ContractProposalUpdateOne) SetContractAddress(s string) *ContractProposalUpdateOne {
	cpuo.mutation.SetContractAddress(s)
	return cpuo
}

// SetStatus sets the "status" field.
func (cpuo *ContractProposalUpdateOne) SetStatus(c contractproposal.Status) *ContractProposalUpdateOne {
	cpuo.mutation.SetStatus(c)
	return cpuo
}

// SetChainID sets the "chain" edge to the Chain entity by ID.
func (cpuo *ContractProposalUpdateOne) SetChainID(id int) *ContractProposalUpdateOne {
	cpuo.mutation.SetChainID(id)
	return cpuo
}

// SetNillableChainID sets the "chain" edge to the Chain entity by ID if the given value is not nil.
func (cpuo *ContractProposalUpdateOne) SetNillableChainID(id *int) *ContractProposalUpdateOne {
	if id != nil {
		cpuo = cpuo.SetChainID(*id)
	}
	return cpuo
}

// SetChain sets the "chain" edge to the Chain entity.
func (cpuo *ContractProposalUpdateOne) SetChain(c *Chain) *ContractProposalUpdateOne {
	return cpuo.SetChainID(c.ID)
}

// Mutation returns the ContractProposalMutation object of the builder.
func (cpuo *ContractProposalUpdateOne) Mutation() *ContractProposalMutation {
	return cpuo.mutation
}

// ClearChain clears the "chain" edge to the Chain entity.
func (cpuo *ContractProposalUpdateOne) ClearChain() *ContractProposalUpdateOne {
	cpuo.mutation.ClearChain()
	return cpuo
}

// Where appends a list predicates to the ContractProposalUpdate builder.
func (cpuo *ContractProposalUpdateOne) Where(ps ...predicate.ContractProposal) *ContractProposalUpdateOne {
	cpuo.mutation.Where(ps...)
	return cpuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (cpuo *ContractProposalUpdateOne) Select(field string, fields ...string) *ContractProposalUpdateOne {
	cpuo.fields = append([]string{field}, fields...)
	return cpuo
}

// Save executes the query and returns the updated ContractProposal entity.
func (cpuo *ContractProposalUpdateOne) Save(ctx context.Context) (*ContractProposal, error) {
	cpuo.defaults()
	return withHooks(ctx, cpuo.sqlSave, cpuo.mutation, cpuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cpuo *ContractProposalUpdateOne) SaveX(ctx context.Context) *ContractProposal {
	node, err := cpuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (cpuo *ContractProposalUpdateOne) Exec(ctx context.Context) error {
	_, err := cpuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cpuo *ContractProposalUpdateOne) ExecX(ctx context.Context) {
	if err := cpuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cpuo *ContractProposalUpdateOne) defaults() {
	if _, ok := cpuo.mutation.UpdateTime(); !ok {
		v := contractproposal.UpdateDefaultUpdateTime()
		cpuo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cpuo *ContractProposalUpdateOne) check() error {
	if v, ok := cpuo.mutation.Status(); ok {
		if err := contractproposal.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "ContractProposal.status": %w`, err)}
		}
	}
	return nil
}

func (cpuo *ContractProposalUpdateOne) sqlSave(ctx context.Context) (_node *ContractProposal, err error) {
	if err := cpuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(contractproposal.Table, contractproposal.Columns, sqlgraph.NewFieldSpec(contractproposal.FieldID, field.TypeInt))
	id, ok := cpuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "ContractProposal.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := cpuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, contractproposal.FieldID)
		for _, f := range fields {
			if !contractproposal.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != contractproposal.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := cpuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cpuo.mutation.UpdateTime(); ok {
		_spec.SetField(contractproposal.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := cpuo.mutation.ProposalID(); ok {
		_spec.SetField(contractproposal.FieldProposalID, field.TypeUint64, value)
	}
	if value, ok := cpuo.mutation.AddedProposalID(); ok {
		_spec.AddField(contractproposal.FieldProposalID, field.TypeUint64, value)
	}
	if value, ok := cpuo.mutation.Title(); ok {
		_spec.SetField(contractproposal.FieldTitle, field.TypeString, value)
	}
	if value, ok := cpuo.mutation.Description(); ok {
		_spec.SetField(contractproposal.FieldDescription, field.TypeString, value)
	}
	if value, ok := cpuo.mutation.FirstSeenTime(); ok {
		_spec.SetField(contractproposal.FieldFirstSeenTime, field.TypeTime, value)
	}
	if value, ok := cpuo.mutation.VotingEndTime(); ok {
		_spec.SetField(contractproposal.FieldVotingEndTime, field.TypeTime, value)
	}
	if value, ok := cpuo.mutation.ContractAddress(); ok {
		_spec.SetField(contractproposal.FieldContractAddress, field.TypeString, value)
	}
	if value, ok := cpuo.mutation.Status(); ok {
		_spec.SetField(contractproposal.FieldStatus, field.TypeEnum, value)
	}
	if cpuo.mutation.ChainCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   contractproposal.ChainTable,
			Columns: []string{contractproposal.ChainColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chain.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cpuo.mutation.ChainIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   contractproposal.ChainTable,
			Columns: []string{contractproposal.ChainColumn},
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
	_node = &ContractProposal{config: cpuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, cpuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{contractproposal.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	cpuo.mutation.done = true
	return _node, nil
}