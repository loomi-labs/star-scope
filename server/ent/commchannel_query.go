// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/loomi-labs/star-scope/ent/commchannel"
	"github.com/loomi-labs/star-scope/ent/eventlistener"
	"github.com/loomi-labs/star-scope/ent/predicate"
	"github.com/loomi-labs/star-scope/ent/user"
)

// CommChannelQuery is the builder for querying CommChannel entities.
type CommChannelQuery struct {
	config
	ctx                *QueryContext
	order              []commchannel.OrderOption
	inters             []Interceptor
	predicates         []predicate.CommChannel
	withEventListeners *EventListenerQuery
	withUsers          *UserQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the CommChannelQuery builder.
func (ccq *CommChannelQuery) Where(ps ...predicate.CommChannel) *CommChannelQuery {
	ccq.predicates = append(ccq.predicates, ps...)
	return ccq
}

// Limit the number of records to be returned by this query.
func (ccq *CommChannelQuery) Limit(limit int) *CommChannelQuery {
	ccq.ctx.Limit = &limit
	return ccq
}

// Offset to start from.
func (ccq *CommChannelQuery) Offset(offset int) *CommChannelQuery {
	ccq.ctx.Offset = &offset
	return ccq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (ccq *CommChannelQuery) Unique(unique bool) *CommChannelQuery {
	ccq.ctx.Unique = &unique
	return ccq
}

// Order specifies how the records should be ordered.
func (ccq *CommChannelQuery) Order(o ...commchannel.OrderOption) *CommChannelQuery {
	ccq.order = append(ccq.order, o...)
	return ccq
}

// QueryEventListeners chains the current query on the "event_listeners" edge.
func (ccq *CommChannelQuery) QueryEventListeners() *EventListenerQuery {
	query := (&EventListenerClient{config: ccq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ccq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := ccq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(commchannel.Table, commchannel.FieldID, selector),
			sqlgraph.To(eventlistener.Table, eventlistener.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, commchannel.EventListenersTable, commchannel.EventListenersPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(ccq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryUsers chains the current query on the "users" edge.
func (ccq *CommChannelQuery) QueryUsers() *UserQuery {
	query := (&UserClient{config: ccq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ccq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := ccq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(commchannel.Table, commchannel.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, commchannel.UsersTable, commchannel.UsersPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(ccq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first CommChannel entity from the query.
// Returns a *NotFoundError when no CommChannel was found.
func (ccq *CommChannelQuery) First(ctx context.Context) (*CommChannel, error) {
	nodes, err := ccq.Limit(1).All(setContextOp(ctx, ccq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{commchannel.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (ccq *CommChannelQuery) FirstX(ctx context.Context) *CommChannel {
	node, err := ccq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first CommChannel ID from the query.
// Returns a *NotFoundError when no CommChannel ID was found.
func (ccq *CommChannelQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ccq.Limit(1).IDs(setContextOp(ctx, ccq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{commchannel.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (ccq *CommChannelQuery) FirstIDX(ctx context.Context) int {
	id, err := ccq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single CommChannel entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one CommChannel entity is found.
// Returns a *NotFoundError when no CommChannel entities are found.
func (ccq *CommChannelQuery) Only(ctx context.Context) (*CommChannel, error) {
	nodes, err := ccq.Limit(2).All(setContextOp(ctx, ccq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{commchannel.Label}
	default:
		return nil, &NotSingularError{commchannel.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (ccq *CommChannelQuery) OnlyX(ctx context.Context) *CommChannel {
	node, err := ccq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only CommChannel ID in the query.
// Returns a *NotSingularError when more than one CommChannel ID is found.
// Returns a *NotFoundError when no entities are found.
func (ccq *CommChannelQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ccq.Limit(2).IDs(setContextOp(ctx, ccq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{commchannel.Label}
	default:
		err = &NotSingularError{commchannel.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (ccq *CommChannelQuery) OnlyIDX(ctx context.Context) int {
	id, err := ccq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of CommChannels.
func (ccq *CommChannelQuery) All(ctx context.Context) ([]*CommChannel, error) {
	ctx = setContextOp(ctx, ccq.ctx, "All")
	if err := ccq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*CommChannel, *CommChannelQuery]()
	return withInterceptors[[]*CommChannel](ctx, ccq, qr, ccq.inters)
}

// AllX is like All, but panics if an error occurs.
func (ccq *CommChannelQuery) AllX(ctx context.Context) []*CommChannel {
	nodes, err := ccq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of CommChannel IDs.
func (ccq *CommChannelQuery) IDs(ctx context.Context) (ids []int, err error) {
	if ccq.ctx.Unique == nil && ccq.path != nil {
		ccq.Unique(true)
	}
	ctx = setContextOp(ctx, ccq.ctx, "IDs")
	if err = ccq.Select(commchannel.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (ccq *CommChannelQuery) IDsX(ctx context.Context) []int {
	ids, err := ccq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (ccq *CommChannelQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, ccq.ctx, "Count")
	if err := ccq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, ccq, querierCount[*CommChannelQuery](), ccq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (ccq *CommChannelQuery) CountX(ctx context.Context) int {
	count, err := ccq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (ccq *CommChannelQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, ccq.ctx, "Exist")
	switch _, err := ccq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (ccq *CommChannelQuery) ExistX(ctx context.Context) bool {
	exist, err := ccq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the CommChannelQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (ccq *CommChannelQuery) Clone() *CommChannelQuery {
	if ccq == nil {
		return nil
	}
	return &CommChannelQuery{
		config:             ccq.config,
		ctx:                ccq.ctx.Clone(),
		order:              append([]commchannel.OrderOption{}, ccq.order...),
		inters:             append([]Interceptor{}, ccq.inters...),
		predicates:         append([]predicate.CommChannel{}, ccq.predicates...),
		withEventListeners: ccq.withEventListeners.Clone(),
		withUsers:          ccq.withUsers.Clone(),
		// clone intermediate query.
		sql:  ccq.sql.Clone(),
		path: ccq.path,
	}
}

// WithEventListeners tells the query-builder to eager-load the nodes that are connected to
// the "event_listeners" edge. The optional arguments are used to configure the query builder of the edge.
func (ccq *CommChannelQuery) WithEventListeners(opts ...func(*EventListenerQuery)) *CommChannelQuery {
	query := (&EventListenerClient{config: ccq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	ccq.withEventListeners = query
	return ccq
}

// WithUsers tells the query-builder to eager-load the nodes that are connected to
// the "users" edge. The optional arguments are used to configure the query builder of the edge.
func (ccq *CommChannelQuery) WithUsers(opts ...func(*UserQuery)) *CommChannelQuery {
	query := (&UserClient{config: ccq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	ccq.withUsers = query
	return ccq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		CreateTime time.Time `json:"create_time,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.CommChannel.Query().
//		GroupBy(commchannel.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (ccq *CommChannelQuery) GroupBy(field string, fields ...string) *CommChannelGroupBy {
	ccq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &CommChannelGroupBy{build: ccq}
	grbuild.flds = &ccq.ctx.Fields
	grbuild.label = commchannel.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		CreateTime time.Time `json:"create_time,omitempty"`
//	}
//
//	client.CommChannel.Query().
//		Select(commchannel.FieldCreateTime).
//		Scan(ctx, &v)
func (ccq *CommChannelQuery) Select(fields ...string) *CommChannelSelect {
	ccq.ctx.Fields = append(ccq.ctx.Fields, fields...)
	sbuild := &CommChannelSelect{CommChannelQuery: ccq}
	sbuild.label = commchannel.Label
	sbuild.flds, sbuild.scan = &ccq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a CommChannelSelect configured with the given aggregations.
func (ccq *CommChannelQuery) Aggregate(fns ...AggregateFunc) *CommChannelSelect {
	return ccq.Select().Aggregate(fns...)
}

func (ccq *CommChannelQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range ccq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, ccq); err != nil {
				return err
			}
		}
	}
	for _, f := range ccq.ctx.Fields {
		if !commchannel.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if ccq.path != nil {
		prev, err := ccq.path(ctx)
		if err != nil {
			return err
		}
		ccq.sql = prev
	}
	return nil
}

func (ccq *CommChannelQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*CommChannel, error) {
	var (
		nodes       = []*CommChannel{}
		_spec       = ccq.querySpec()
		loadedTypes = [2]bool{
			ccq.withEventListeners != nil,
			ccq.withUsers != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*CommChannel).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &CommChannel{config: ccq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, ccq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := ccq.withEventListeners; query != nil {
		if err := ccq.loadEventListeners(ctx, query, nodes,
			func(n *CommChannel) { n.Edges.EventListeners = []*EventListener{} },
			func(n *CommChannel, e *EventListener) { n.Edges.EventListeners = append(n.Edges.EventListeners, e) }); err != nil {
			return nil, err
		}
	}
	if query := ccq.withUsers; query != nil {
		if err := ccq.loadUsers(ctx, query, nodes,
			func(n *CommChannel) { n.Edges.Users = []*User{} },
			func(n *CommChannel, e *User) { n.Edges.Users = append(n.Edges.Users, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (ccq *CommChannelQuery) loadEventListeners(ctx context.Context, query *EventListenerQuery, nodes []*CommChannel, init func(*CommChannel), assign func(*CommChannel, *EventListener)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[int]*CommChannel)
	nids := make(map[int]map[*CommChannel]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(commchannel.EventListenersTable)
		s.Join(joinT).On(s.C(eventlistener.FieldID), joinT.C(commchannel.EventListenersPrimaryKey[1]))
		s.Where(sql.InValues(joinT.C(commchannel.EventListenersPrimaryKey[0]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(commchannel.EventListenersPrimaryKey[0]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	qr := QuerierFunc(func(ctx context.Context, q Query) (Value, error) {
		return query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
			assign := spec.Assign
			values := spec.ScanValues
			spec.ScanValues = func(columns []string) ([]any, error) {
				values, err := values(columns[1:])
				if err != nil {
					return nil, err
				}
				return append([]any{new(sql.NullInt64)}, values...), nil
			}
			spec.Assign = func(columns []string, values []any) error {
				outValue := int(values[0].(*sql.NullInt64).Int64)
				inValue := int(values[1].(*sql.NullInt64).Int64)
				if nids[inValue] == nil {
					nids[inValue] = map[*CommChannel]struct{}{byID[outValue]: {}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byID[outValue]] = struct{}{}
				return nil
			}
		})
	})
	neighbors, err := withInterceptors[[]*EventListener](ctx, query, qr, query.inters)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "event_listeners" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}
func (ccq *CommChannelQuery) loadUsers(ctx context.Context, query *UserQuery, nodes []*CommChannel, init func(*CommChannel), assign func(*CommChannel, *User)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[int]*CommChannel)
	nids := make(map[int]map[*CommChannel]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(commchannel.UsersTable)
		s.Join(joinT).On(s.C(user.FieldID), joinT.C(commchannel.UsersPrimaryKey[0]))
		s.Where(sql.InValues(joinT.C(commchannel.UsersPrimaryKey[1]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(commchannel.UsersPrimaryKey[1]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	qr := QuerierFunc(func(ctx context.Context, q Query) (Value, error) {
		return query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
			assign := spec.Assign
			values := spec.ScanValues
			spec.ScanValues = func(columns []string) ([]any, error) {
				values, err := values(columns[1:])
				if err != nil {
					return nil, err
				}
				return append([]any{new(sql.NullInt64)}, values...), nil
			}
			spec.Assign = func(columns []string, values []any) error {
				outValue := int(values[0].(*sql.NullInt64).Int64)
				inValue := int(values[1].(*sql.NullInt64).Int64)
				if nids[inValue] == nil {
					nids[inValue] = map[*CommChannel]struct{}{byID[outValue]: {}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byID[outValue]] = struct{}{}
				return nil
			}
		})
	})
	neighbors, err := withInterceptors[[]*User](ctx, query, qr, query.inters)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "users" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}

func (ccq *CommChannelQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := ccq.querySpec()
	_spec.Node.Columns = ccq.ctx.Fields
	if len(ccq.ctx.Fields) > 0 {
		_spec.Unique = ccq.ctx.Unique != nil && *ccq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, ccq.driver, _spec)
}

func (ccq *CommChannelQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(commchannel.Table, commchannel.Columns, sqlgraph.NewFieldSpec(commchannel.FieldID, field.TypeInt))
	_spec.From = ccq.sql
	if unique := ccq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if ccq.path != nil {
		_spec.Unique = true
	}
	if fields := ccq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, commchannel.FieldID)
		for i := range fields {
			if fields[i] != commchannel.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := ccq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := ccq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := ccq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := ccq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (ccq *CommChannelQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(ccq.driver.Dialect())
	t1 := builder.Table(commchannel.Table)
	columns := ccq.ctx.Fields
	if len(columns) == 0 {
		columns = commchannel.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if ccq.sql != nil {
		selector = ccq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if ccq.ctx.Unique != nil && *ccq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range ccq.predicates {
		p(selector)
	}
	for _, p := range ccq.order {
		p(selector)
	}
	if offset := ccq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := ccq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// CommChannelGroupBy is the group-by builder for CommChannel entities.
type CommChannelGroupBy struct {
	selector
	build *CommChannelQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (ccgb *CommChannelGroupBy) Aggregate(fns ...AggregateFunc) *CommChannelGroupBy {
	ccgb.fns = append(ccgb.fns, fns...)
	return ccgb
}

// Scan applies the selector query and scans the result into the given value.
func (ccgb *CommChannelGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ccgb.build.ctx, "GroupBy")
	if err := ccgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CommChannelQuery, *CommChannelGroupBy](ctx, ccgb.build, ccgb, ccgb.build.inters, v)
}

func (ccgb *CommChannelGroupBy) sqlScan(ctx context.Context, root *CommChannelQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(ccgb.fns))
	for _, fn := range ccgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*ccgb.flds)+len(ccgb.fns))
		for _, f := range *ccgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*ccgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ccgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// CommChannelSelect is the builder for selecting fields of CommChannel entities.
type CommChannelSelect struct {
	*CommChannelQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ccs *CommChannelSelect) Aggregate(fns ...AggregateFunc) *CommChannelSelect {
	ccs.fns = append(ccs.fns, fns...)
	return ccs
}

// Scan applies the selector query and scans the result into the given value.
func (ccs *CommChannelSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ccs.ctx, "Select")
	if err := ccs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CommChannelQuery, *CommChannelSelect](ctx, ccs.CommChannelQuery, ccs, ccs.inters, v)
}

func (ccs *CommChannelSelect) sqlScan(ctx context.Context, root *CommChannelQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ccs.fns))
	for _, fn := range ccs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ccs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ccs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}