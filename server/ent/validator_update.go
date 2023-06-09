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
	"github.com/loomi-labs/star-scope/ent/usersetup"
	"github.com/loomi-labs/star-scope/ent/validator"
)

// ValidatorUpdate is the builder for updating Validator entities.
type ValidatorUpdate struct {
	config
	hooks    []Hook
	mutation *ValidatorMutation
}

// Where appends a list predicates to the ValidatorUpdate builder.
func (vu *ValidatorUpdate) Where(ps ...predicate.Validator) *ValidatorUpdate {
	vu.mutation.Where(ps...)
	return vu
}

// SetUpdateTime sets the "update_time" field.
func (vu *ValidatorUpdate) SetUpdateTime(t time.Time) *ValidatorUpdate {
	vu.mutation.SetUpdateTime(t)
	return vu
}

// SetMoniker sets the "moniker" field.
func (vu *ValidatorUpdate) SetMoniker(s string) *ValidatorUpdate {
	vu.mutation.SetMoniker(s)
	return vu
}

// SetFirstInactiveTime sets the "first_inactive_time" field.
func (vu *ValidatorUpdate) SetFirstInactiveTime(t time.Time) *ValidatorUpdate {
	vu.mutation.SetFirstInactiveTime(t)
	return vu
}

// SetNillableFirstInactiveTime sets the "first_inactive_time" field if the given value is not nil.
func (vu *ValidatorUpdate) SetNillableFirstInactiveTime(t *time.Time) *ValidatorUpdate {
	if t != nil {
		vu.SetFirstInactiveTime(*t)
	}
	return vu
}

// ClearFirstInactiveTime clears the value of the "first_inactive_time" field.
func (vu *ValidatorUpdate) ClearFirstInactiveTime() *ValidatorUpdate {
	vu.mutation.ClearFirstInactiveTime()
	return vu
}

// SetLastSlashValidatorPeriod sets the "last_slash_validator_period" field.
func (vu *ValidatorUpdate) SetLastSlashValidatorPeriod(u uint64) *ValidatorUpdate {
	vu.mutation.ResetLastSlashValidatorPeriod()
	vu.mutation.SetLastSlashValidatorPeriod(u)
	return vu
}

// SetNillableLastSlashValidatorPeriod sets the "last_slash_validator_period" field if the given value is not nil.
func (vu *ValidatorUpdate) SetNillableLastSlashValidatorPeriod(u *uint64) *ValidatorUpdate {
	if u != nil {
		vu.SetLastSlashValidatorPeriod(*u)
	}
	return vu
}

// AddLastSlashValidatorPeriod adds u to the "last_slash_validator_period" field.
func (vu *ValidatorUpdate) AddLastSlashValidatorPeriod(u int64) *ValidatorUpdate {
	vu.mutation.AddLastSlashValidatorPeriod(u)
	return vu
}

// ClearLastSlashValidatorPeriod clears the value of the "last_slash_validator_period" field.
func (vu *ValidatorUpdate) ClearLastSlashValidatorPeriod() *ValidatorUpdate {
	vu.mutation.ClearLastSlashValidatorPeriod()
	return vu
}

// SetChainID sets the "chain" edge to the Chain entity by ID.
func (vu *ValidatorUpdate) SetChainID(id int) *ValidatorUpdate {
	vu.mutation.SetChainID(id)
	return vu
}

// SetChain sets the "chain" edge to the Chain entity.
func (vu *ValidatorUpdate) SetChain(c *Chain) *ValidatorUpdate {
	return vu.SetChainID(c.ID)
}

// AddSelectedBySetupIDs adds the "selected_by_setups" edge to the UserSetup entity by IDs.
func (vu *ValidatorUpdate) AddSelectedBySetupIDs(ids ...int) *ValidatorUpdate {
	vu.mutation.AddSelectedBySetupIDs(ids...)
	return vu
}

// AddSelectedBySetups adds the "selected_by_setups" edges to the UserSetup entity.
func (vu *ValidatorUpdate) AddSelectedBySetups(u ...*UserSetup) *ValidatorUpdate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return vu.AddSelectedBySetupIDs(ids...)
}

// Mutation returns the ValidatorMutation object of the builder.
func (vu *ValidatorUpdate) Mutation() *ValidatorMutation {
	return vu.mutation
}

// ClearChain clears the "chain" edge to the Chain entity.
func (vu *ValidatorUpdate) ClearChain() *ValidatorUpdate {
	vu.mutation.ClearChain()
	return vu
}

// ClearSelectedBySetups clears all "selected_by_setups" edges to the UserSetup entity.
func (vu *ValidatorUpdate) ClearSelectedBySetups() *ValidatorUpdate {
	vu.mutation.ClearSelectedBySetups()
	return vu
}

// RemoveSelectedBySetupIDs removes the "selected_by_setups" edge to UserSetup entities by IDs.
func (vu *ValidatorUpdate) RemoveSelectedBySetupIDs(ids ...int) *ValidatorUpdate {
	vu.mutation.RemoveSelectedBySetupIDs(ids...)
	return vu
}

// RemoveSelectedBySetups removes "selected_by_setups" edges to UserSetup entities.
func (vu *ValidatorUpdate) RemoveSelectedBySetups(u ...*UserSetup) *ValidatorUpdate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return vu.RemoveSelectedBySetupIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (vu *ValidatorUpdate) Save(ctx context.Context) (int, error) {
	vu.defaults()
	return withHooks(ctx, vu.sqlSave, vu.mutation, vu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (vu *ValidatorUpdate) SaveX(ctx context.Context) int {
	affected, err := vu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (vu *ValidatorUpdate) Exec(ctx context.Context) error {
	_, err := vu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (vu *ValidatorUpdate) ExecX(ctx context.Context) {
	if err := vu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (vu *ValidatorUpdate) defaults() {
	if _, ok := vu.mutation.UpdateTime(); !ok {
		v := validator.UpdateDefaultUpdateTime()
		vu.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (vu *ValidatorUpdate) check() error {
	if _, ok := vu.mutation.ChainID(); vu.mutation.ChainCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Validator.chain"`)
	}
	return nil
}

func (vu *ValidatorUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := vu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(validator.Table, validator.Columns, sqlgraph.NewFieldSpec(validator.FieldID, field.TypeInt))
	if ps := vu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := vu.mutation.UpdateTime(); ok {
		_spec.SetField(validator.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := vu.mutation.Moniker(); ok {
		_spec.SetField(validator.FieldMoniker, field.TypeString, value)
	}
	if value, ok := vu.mutation.FirstInactiveTime(); ok {
		_spec.SetField(validator.FieldFirstInactiveTime, field.TypeTime, value)
	}
	if vu.mutation.FirstInactiveTimeCleared() {
		_spec.ClearField(validator.FieldFirstInactiveTime, field.TypeTime)
	}
	if value, ok := vu.mutation.LastSlashValidatorPeriod(); ok {
		_spec.SetField(validator.FieldLastSlashValidatorPeriod, field.TypeUint64, value)
	}
	if value, ok := vu.mutation.AddedLastSlashValidatorPeriod(); ok {
		_spec.AddField(validator.FieldLastSlashValidatorPeriod, field.TypeUint64, value)
	}
	if vu.mutation.LastSlashValidatorPeriodCleared() {
		_spec.ClearField(validator.FieldLastSlashValidatorPeriod, field.TypeUint64)
	}
	if vu.mutation.ChainCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := vu.mutation.ChainIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if vu.mutation.SelectedBySetupsCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := vu.mutation.RemovedSelectedBySetupsIDs(); len(nodes) > 0 && !vu.mutation.SelectedBySetupsCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := vu.mutation.SelectedBySetupsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, vu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{validator.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	vu.mutation.done = true
	return n, nil
}

// ValidatorUpdateOne is the builder for updating a single Validator entity.
type ValidatorUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ValidatorMutation
}

// SetUpdateTime sets the "update_time" field.
func (vuo *ValidatorUpdateOne) SetUpdateTime(t time.Time) *ValidatorUpdateOne {
	vuo.mutation.SetUpdateTime(t)
	return vuo
}

// SetMoniker sets the "moniker" field.
func (vuo *ValidatorUpdateOne) SetMoniker(s string) *ValidatorUpdateOne {
	vuo.mutation.SetMoniker(s)
	return vuo
}

// SetFirstInactiveTime sets the "first_inactive_time" field.
func (vuo *ValidatorUpdateOne) SetFirstInactiveTime(t time.Time) *ValidatorUpdateOne {
	vuo.mutation.SetFirstInactiveTime(t)
	return vuo
}

// SetNillableFirstInactiveTime sets the "first_inactive_time" field if the given value is not nil.
func (vuo *ValidatorUpdateOne) SetNillableFirstInactiveTime(t *time.Time) *ValidatorUpdateOne {
	if t != nil {
		vuo.SetFirstInactiveTime(*t)
	}
	return vuo
}

// ClearFirstInactiveTime clears the value of the "first_inactive_time" field.
func (vuo *ValidatorUpdateOne) ClearFirstInactiveTime() *ValidatorUpdateOne {
	vuo.mutation.ClearFirstInactiveTime()
	return vuo
}

// SetLastSlashValidatorPeriod sets the "last_slash_validator_period" field.
func (vuo *ValidatorUpdateOne) SetLastSlashValidatorPeriod(u uint64) *ValidatorUpdateOne {
	vuo.mutation.ResetLastSlashValidatorPeriod()
	vuo.mutation.SetLastSlashValidatorPeriod(u)
	return vuo
}

// SetNillableLastSlashValidatorPeriod sets the "last_slash_validator_period" field if the given value is not nil.
func (vuo *ValidatorUpdateOne) SetNillableLastSlashValidatorPeriod(u *uint64) *ValidatorUpdateOne {
	if u != nil {
		vuo.SetLastSlashValidatorPeriod(*u)
	}
	return vuo
}

// AddLastSlashValidatorPeriod adds u to the "last_slash_validator_period" field.
func (vuo *ValidatorUpdateOne) AddLastSlashValidatorPeriod(u int64) *ValidatorUpdateOne {
	vuo.mutation.AddLastSlashValidatorPeriod(u)
	return vuo
}

// ClearLastSlashValidatorPeriod clears the value of the "last_slash_validator_period" field.
func (vuo *ValidatorUpdateOne) ClearLastSlashValidatorPeriod() *ValidatorUpdateOne {
	vuo.mutation.ClearLastSlashValidatorPeriod()
	return vuo
}

// SetChainID sets the "chain" edge to the Chain entity by ID.
func (vuo *ValidatorUpdateOne) SetChainID(id int) *ValidatorUpdateOne {
	vuo.mutation.SetChainID(id)
	return vuo
}

// SetChain sets the "chain" edge to the Chain entity.
func (vuo *ValidatorUpdateOne) SetChain(c *Chain) *ValidatorUpdateOne {
	return vuo.SetChainID(c.ID)
}

// AddSelectedBySetupIDs adds the "selected_by_setups" edge to the UserSetup entity by IDs.
func (vuo *ValidatorUpdateOne) AddSelectedBySetupIDs(ids ...int) *ValidatorUpdateOne {
	vuo.mutation.AddSelectedBySetupIDs(ids...)
	return vuo
}

// AddSelectedBySetups adds the "selected_by_setups" edges to the UserSetup entity.
func (vuo *ValidatorUpdateOne) AddSelectedBySetups(u ...*UserSetup) *ValidatorUpdateOne {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return vuo.AddSelectedBySetupIDs(ids...)
}

// Mutation returns the ValidatorMutation object of the builder.
func (vuo *ValidatorUpdateOne) Mutation() *ValidatorMutation {
	return vuo.mutation
}

// ClearChain clears the "chain" edge to the Chain entity.
func (vuo *ValidatorUpdateOne) ClearChain() *ValidatorUpdateOne {
	vuo.mutation.ClearChain()
	return vuo
}

// ClearSelectedBySetups clears all "selected_by_setups" edges to the UserSetup entity.
func (vuo *ValidatorUpdateOne) ClearSelectedBySetups() *ValidatorUpdateOne {
	vuo.mutation.ClearSelectedBySetups()
	return vuo
}

// RemoveSelectedBySetupIDs removes the "selected_by_setups" edge to UserSetup entities by IDs.
func (vuo *ValidatorUpdateOne) RemoveSelectedBySetupIDs(ids ...int) *ValidatorUpdateOne {
	vuo.mutation.RemoveSelectedBySetupIDs(ids...)
	return vuo
}

// RemoveSelectedBySetups removes "selected_by_setups" edges to UserSetup entities.
func (vuo *ValidatorUpdateOne) RemoveSelectedBySetups(u ...*UserSetup) *ValidatorUpdateOne {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return vuo.RemoveSelectedBySetupIDs(ids...)
}

// Where appends a list predicates to the ValidatorUpdate builder.
func (vuo *ValidatorUpdateOne) Where(ps ...predicate.Validator) *ValidatorUpdateOne {
	vuo.mutation.Where(ps...)
	return vuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (vuo *ValidatorUpdateOne) Select(field string, fields ...string) *ValidatorUpdateOne {
	vuo.fields = append([]string{field}, fields...)
	return vuo
}

// Save executes the query and returns the updated Validator entity.
func (vuo *ValidatorUpdateOne) Save(ctx context.Context) (*Validator, error) {
	vuo.defaults()
	return withHooks(ctx, vuo.sqlSave, vuo.mutation, vuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (vuo *ValidatorUpdateOne) SaveX(ctx context.Context) *Validator {
	node, err := vuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (vuo *ValidatorUpdateOne) Exec(ctx context.Context) error {
	_, err := vuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (vuo *ValidatorUpdateOne) ExecX(ctx context.Context) {
	if err := vuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (vuo *ValidatorUpdateOne) defaults() {
	if _, ok := vuo.mutation.UpdateTime(); !ok {
		v := validator.UpdateDefaultUpdateTime()
		vuo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (vuo *ValidatorUpdateOne) check() error {
	if _, ok := vuo.mutation.ChainID(); vuo.mutation.ChainCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Validator.chain"`)
	}
	return nil
}

func (vuo *ValidatorUpdateOne) sqlSave(ctx context.Context) (_node *Validator, err error) {
	if err := vuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(validator.Table, validator.Columns, sqlgraph.NewFieldSpec(validator.FieldID, field.TypeInt))
	id, ok := vuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Validator.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := vuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, validator.FieldID)
		for _, f := range fields {
			if !validator.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != validator.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := vuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := vuo.mutation.UpdateTime(); ok {
		_spec.SetField(validator.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := vuo.mutation.Moniker(); ok {
		_spec.SetField(validator.FieldMoniker, field.TypeString, value)
	}
	if value, ok := vuo.mutation.FirstInactiveTime(); ok {
		_spec.SetField(validator.FieldFirstInactiveTime, field.TypeTime, value)
	}
	if vuo.mutation.FirstInactiveTimeCleared() {
		_spec.ClearField(validator.FieldFirstInactiveTime, field.TypeTime)
	}
	if value, ok := vuo.mutation.LastSlashValidatorPeriod(); ok {
		_spec.SetField(validator.FieldLastSlashValidatorPeriod, field.TypeUint64, value)
	}
	if value, ok := vuo.mutation.AddedLastSlashValidatorPeriod(); ok {
		_spec.AddField(validator.FieldLastSlashValidatorPeriod, field.TypeUint64, value)
	}
	if vuo.mutation.LastSlashValidatorPeriodCleared() {
		_spec.ClearField(validator.FieldLastSlashValidatorPeriod, field.TypeUint64)
	}
	if vuo.mutation.ChainCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := vuo.mutation.ChainIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if vuo.mutation.SelectedBySetupsCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := vuo.mutation.RemovedSelectedBySetupsIDs(); len(nodes) > 0 && !vuo.mutation.SelectedBySetupsCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := vuo.mutation.SelectedBySetupsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Validator{config: vuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, vuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{validator.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	vuo.mutation.done = true
	return _node, nil
}
