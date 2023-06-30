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
	"github.com/loomi-labs/star-scope/ent/usersetup"
	"github.com/loomi-labs/star-scope/ent/validator"
)

// ValidatorCreate is the builder for creating a Validator entity.
type ValidatorCreate struct {
	config
	mutation *ValidatorMutation
	hooks    []Hook
}

// SetCreateTime sets the "create_time" field.
func (vc *ValidatorCreate) SetCreateTime(t time.Time) *ValidatorCreate {
	vc.mutation.SetCreateTime(t)
	return vc
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (vc *ValidatorCreate) SetNillableCreateTime(t *time.Time) *ValidatorCreate {
	if t != nil {
		vc.SetCreateTime(*t)
	}
	return vc
}

// SetUpdateTime sets the "update_time" field.
func (vc *ValidatorCreate) SetUpdateTime(t time.Time) *ValidatorCreate {
	vc.mutation.SetUpdateTime(t)
	return vc
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (vc *ValidatorCreate) SetNillableUpdateTime(t *time.Time) *ValidatorCreate {
	if t != nil {
		vc.SetUpdateTime(*t)
	}
	return vc
}

// SetOperatorAddress sets the "operator_address" field.
func (vc *ValidatorCreate) SetOperatorAddress(s string) *ValidatorCreate {
	vc.mutation.SetOperatorAddress(s)
	return vc
}

// SetAddress sets the "address" field.
func (vc *ValidatorCreate) SetAddress(s string) *ValidatorCreate {
	vc.mutation.SetAddress(s)
	return vc
}

// SetMoniker sets the "moniker" field.
func (vc *ValidatorCreate) SetMoniker(s string) *ValidatorCreate {
	vc.mutation.SetMoniker(s)
	return vc
}

// SetFirstInactiveTime sets the "first_inactive_time" field.
func (vc *ValidatorCreate) SetFirstInactiveTime(t time.Time) *ValidatorCreate {
	vc.mutation.SetFirstInactiveTime(t)
	return vc
}

// SetNillableFirstInactiveTime sets the "first_inactive_time" field if the given value is not nil.
func (vc *ValidatorCreate) SetNillableFirstInactiveTime(t *time.Time) *ValidatorCreate {
	if t != nil {
		vc.SetFirstInactiveTime(*t)
	}
	return vc
}

// SetLastSlashValidatorPeriod sets the "last_slash_validator_period" field.
func (vc *ValidatorCreate) SetLastSlashValidatorPeriod(u uint64) *ValidatorCreate {
	vc.mutation.SetLastSlashValidatorPeriod(u)
	return vc
}

// SetNillableLastSlashValidatorPeriod sets the "last_slash_validator_period" field if the given value is not nil.
func (vc *ValidatorCreate) SetNillableLastSlashValidatorPeriod(u *uint64) *ValidatorCreate {
	if u != nil {
		vc.SetLastSlashValidatorPeriod(*u)
	}
	return vc
}

// SetChainID sets the "chain" edge to the Chain entity by ID.
func (vc *ValidatorCreate) SetChainID(id int) *ValidatorCreate {
	vc.mutation.SetChainID(id)
	return vc
}

// SetChain sets the "chain" edge to the Chain entity.
func (vc *ValidatorCreate) SetChain(c *Chain) *ValidatorCreate {
	return vc.SetChainID(c.ID)
}

// AddSelectedBySetupIDs adds the "selected_by_setups" edge to the UserSetup entity by IDs.
func (vc *ValidatorCreate) AddSelectedBySetupIDs(ids ...int) *ValidatorCreate {
	vc.mutation.AddSelectedBySetupIDs(ids...)
	return vc
}

// AddSelectedBySetups adds the "selected_by_setups" edges to the UserSetup entity.
func (vc *ValidatorCreate) AddSelectedBySetups(u ...*UserSetup) *ValidatorCreate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return vc.AddSelectedBySetupIDs(ids...)
}

// Mutation returns the ValidatorMutation object of the builder.
func (vc *ValidatorCreate) Mutation() *ValidatorMutation {
	return vc.mutation
}

// Save creates the Validator in the database.
func (vc *ValidatorCreate) Save(ctx context.Context) (*Validator, error) {
	vc.defaults()
	return withHooks(ctx, vc.sqlSave, vc.mutation, vc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (vc *ValidatorCreate) SaveX(ctx context.Context) *Validator {
	v, err := vc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (vc *ValidatorCreate) Exec(ctx context.Context) error {
	_, err := vc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (vc *ValidatorCreate) ExecX(ctx context.Context) {
	if err := vc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (vc *ValidatorCreate) defaults() {
	if _, ok := vc.mutation.CreateTime(); !ok {
		v := validator.DefaultCreateTime()
		vc.mutation.SetCreateTime(v)
	}
	if _, ok := vc.mutation.UpdateTime(); !ok {
		v := validator.DefaultUpdateTime()
		vc.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (vc *ValidatorCreate) check() error {
	if _, ok := vc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "create_time", err: errors.New(`ent: missing required field "Validator.create_time"`)}
	}
	if _, ok := vc.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "update_time", err: errors.New(`ent: missing required field "Validator.update_time"`)}
	}
	if _, ok := vc.mutation.OperatorAddress(); !ok {
		return &ValidationError{Name: "operator_address", err: errors.New(`ent: missing required field "Validator.operator_address"`)}
	}
	if v, ok := vc.mutation.OperatorAddress(); ok {
		if err := validator.OperatorAddressValidator(v); err != nil {
			return &ValidationError{Name: "operator_address", err: fmt.Errorf(`ent: validator failed for field "Validator.operator_address": %w`, err)}
		}
	}
	if _, ok := vc.mutation.Address(); !ok {
		return &ValidationError{Name: "address", err: errors.New(`ent: missing required field "Validator.address"`)}
	}
	if v, ok := vc.mutation.Address(); ok {
		if err := validator.AddressValidator(v); err != nil {
			return &ValidationError{Name: "address", err: fmt.Errorf(`ent: validator failed for field "Validator.address": %w`, err)}
		}
	}
	if _, ok := vc.mutation.Moniker(); !ok {
		return &ValidationError{Name: "moniker", err: errors.New(`ent: missing required field "Validator.moniker"`)}
	}
	if _, ok := vc.mutation.ChainID(); !ok {
		return &ValidationError{Name: "chain", err: errors.New(`ent: missing required edge "Validator.chain"`)}
	}
	return nil
}

func (vc *ValidatorCreate) sqlSave(ctx context.Context) (*Validator, error) {
	if err := vc.check(); err != nil {
		return nil, err
	}
	_node, _spec := vc.createSpec()
	if err := sqlgraph.CreateNode(ctx, vc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	vc.mutation.id = &_node.ID
	vc.mutation.done = true
	return _node, nil
}

func (vc *ValidatorCreate) createSpec() (*Validator, *sqlgraph.CreateSpec) {
	var (
		_node = &Validator{config: vc.config}
		_spec = sqlgraph.NewCreateSpec(validator.Table, sqlgraph.NewFieldSpec(validator.FieldID, field.TypeInt))
	)
	if value, ok := vc.mutation.CreateTime(); ok {
		_spec.SetField(validator.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = value
	}
	if value, ok := vc.mutation.UpdateTime(); ok {
		_spec.SetField(validator.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = value
	}
	if value, ok := vc.mutation.OperatorAddress(); ok {
		_spec.SetField(validator.FieldOperatorAddress, field.TypeString, value)
		_node.OperatorAddress = value
	}
	if value, ok := vc.mutation.Address(); ok {
		_spec.SetField(validator.FieldAddress, field.TypeString, value)
		_node.Address = value
	}
	if value, ok := vc.mutation.Moniker(); ok {
		_spec.SetField(validator.FieldMoniker, field.TypeString, value)
		_node.Moniker = value
	}
	if value, ok := vc.mutation.FirstInactiveTime(); ok {
		_spec.SetField(validator.FieldFirstInactiveTime, field.TypeTime, value)
		_node.FirstInactiveTime = &value
	}
	if value, ok := vc.mutation.LastSlashValidatorPeriod(); ok {
		_spec.SetField(validator.FieldLastSlashValidatorPeriod, field.TypeUint64, value)
		_node.LastSlashValidatorPeriod = &value
	}
	if nodes := vc.mutation.ChainIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   validator.ChainTable,
			Columns: []string{validator.ChainColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chain.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.chain_validators = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := vc.mutation.SelectedBySetupsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   validator.SelectedBySetupsTable,
			Columns: validator.SelectedBySetupsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(usersetup.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// ValidatorCreateBulk is the builder for creating many Validator entities in bulk.
type ValidatorCreateBulk struct {
	config
	builders []*ValidatorCreate
}

// Save creates the Validator entities in the database.
func (vcb *ValidatorCreateBulk) Save(ctx context.Context) ([]*Validator, error) {
	specs := make([]*sqlgraph.CreateSpec, len(vcb.builders))
	nodes := make([]*Validator, len(vcb.builders))
	mutators := make([]Mutator, len(vcb.builders))
	for i := range vcb.builders {
		func(i int, root context.Context) {
			builder := vcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ValidatorMutation)
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
					_, err = mutators[i+1].Mutate(root, vcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, vcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, vcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (vcb *ValidatorCreateBulk) SaveX(ctx context.Context) []*Validator {
	v, err := vcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (vcb *ValidatorCreateBulk) Exec(ctx context.Context) error {
	_, err := vcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (vcb *ValidatorCreateBulk) ExecX(ctx context.Context) {
	if err := vcb.Exec(ctx); err != nil {
		panic(err)
	}
}
