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
	"github.com/shifty11/blocklog-backend/ent/chain"
	"github.com/shifty11/blocklog-backend/ent/channel"
	"github.com/shifty11/blocklog-backend/ent/event"
	"github.com/shifty11/blocklog-backend/ent/eventlistener"
	"github.com/shifty11/blocklog-backend/ent/predicate"
)

// EventListenerQuery is the builder for querying EventListener entities.
type EventListenerQuery struct {
	config
	ctx         *QueryContext
	order       []eventlistener.OrderOption
	inters      []Interceptor
	predicates  []predicate.EventListener
	withChannel *ChannelQuery
	withChain   *ChainQuery
	withEvents  *EventQuery
	withFKs     bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the EventListenerQuery builder.
func (elq *EventListenerQuery) Where(ps ...predicate.EventListener) *EventListenerQuery {
	elq.predicates = append(elq.predicates, ps...)
	return elq
}

// Limit the number of records to be returned by this query.
func (elq *EventListenerQuery) Limit(limit int) *EventListenerQuery {
	elq.ctx.Limit = &limit
	return elq
}

// Offset to start from.
func (elq *EventListenerQuery) Offset(offset int) *EventListenerQuery {
	elq.ctx.Offset = &offset
	return elq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (elq *EventListenerQuery) Unique(unique bool) *EventListenerQuery {
	elq.ctx.Unique = &unique
	return elq
}

// Order specifies how the records should be ordered.
func (elq *EventListenerQuery) Order(o ...eventlistener.OrderOption) *EventListenerQuery {
	elq.order = append(elq.order, o...)
	return elq
}

// QueryChannel chains the current query on the "channel" edge.
func (elq *EventListenerQuery) QueryChannel() *ChannelQuery {
	query := (&ChannelClient{config: elq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := elq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := elq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(eventlistener.Table, eventlistener.FieldID, selector),
			sqlgraph.To(channel.Table, channel.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, eventlistener.ChannelTable, eventlistener.ChannelColumn),
		)
		fromU = sqlgraph.SetNeighbors(elq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryChain chains the current query on the "chain" edge.
func (elq *EventListenerQuery) QueryChain() *ChainQuery {
	query := (&ChainClient{config: elq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := elq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := elq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(eventlistener.Table, eventlistener.FieldID, selector),
			sqlgraph.To(chain.Table, chain.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, eventlistener.ChainTable, eventlistener.ChainColumn),
		)
		fromU = sqlgraph.SetNeighbors(elq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryEvents chains the current query on the "events" edge.
func (elq *EventListenerQuery) QueryEvents() *EventQuery {
	query := (&EventClient{config: elq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := elq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := elq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(eventlistener.Table, eventlistener.FieldID, selector),
			sqlgraph.To(event.Table, event.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, eventlistener.EventsTable, eventlistener.EventsColumn),
		)
		fromU = sqlgraph.SetNeighbors(elq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first EventListener entity from the query.
// Returns a *NotFoundError when no EventListener was found.
func (elq *EventListenerQuery) First(ctx context.Context) (*EventListener, error) {
	nodes, err := elq.Limit(1).All(setContextOp(ctx, elq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{eventlistener.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (elq *EventListenerQuery) FirstX(ctx context.Context) *EventListener {
	node, err := elq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first EventListener ID from the query.
// Returns a *NotFoundError when no EventListener ID was found.
func (elq *EventListenerQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = elq.Limit(1).IDs(setContextOp(ctx, elq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{eventlistener.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (elq *EventListenerQuery) FirstIDX(ctx context.Context) int {
	id, err := elq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single EventListener entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one EventListener entity is found.
// Returns a *NotFoundError when no EventListener entities are found.
func (elq *EventListenerQuery) Only(ctx context.Context) (*EventListener, error) {
	nodes, err := elq.Limit(2).All(setContextOp(ctx, elq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{eventlistener.Label}
	default:
		return nil, &NotSingularError{eventlistener.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (elq *EventListenerQuery) OnlyX(ctx context.Context) *EventListener {
	node, err := elq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only EventListener ID in the query.
// Returns a *NotSingularError when more than one EventListener ID is found.
// Returns a *NotFoundError when no entities are found.
func (elq *EventListenerQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = elq.Limit(2).IDs(setContextOp(ctx, elq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{eventlistener.Label}
	default:
		err = &NotSingularError{eventlistener.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (elq *EventListenerQuery) OnlyIDX(ctx context.Context) int {
	id, err := elq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of EventListeners.
func (elq *EventListenerQuery) All(ctx context.Context) ([]*EventListener, error) {
	ctx = setContextOp(ctx, elq.ctx, "All")
	if err := elq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*EventListener, *EventListenerQuery]()
	return withInterceptors[[]*EventListener](ctx, elq, qr, elq.inters)
}

// AllX is like All, but panics if an error occurs.
func (elq *EventListenerQuery) AllX(ctx context.Context) []*EventListener {
	nodes, err := elq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of EventListener IDs.
func (elq *EventListenerQuery) IDs(ctx context.Context) (ids []int, err error) {
	if elq.ctx.Unique == nil && elq.path != nil {
		elq.Unique(true)
	}
	ctx = setContextOp(ctx, elq.ctx, "IDs")
	if err = elq.Select(eventlistener.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (elq *EventListenerQuery) IDsX(ctx context.Context) []int {
	ids, err := elq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (elq *EventListenerQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, elq.ctx, "Count")
	if err := elq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, elq, querierCount[*EventListenerQuery](), elq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (elq *EventListenerQuery) CountX(ctx context.Context) int {
	count, err := elq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (elq *EventListenerQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, elq.ctx, "Exist")
	switch _, err := elq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (elq *EventListenerQuery) ExistX(ctx context.Context) bool {
	exist, err := elq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the EventListenerQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (elq *EventListenerQuery) Clone() *EventListenerQuery {
	if elq == nil {
		return nil
	}
	return &EventListenerQuery{
		config:      elq.config,
		ctx:         elq.ctx.Clone(),
		order:       append([]eventlistener.OrderOption{}, elq.order...),
		inters:      append([]Interceptor{}, elq.inters...),
		predicates:  append([]predicate.EventListener{}, elq.predicates...),
		withChannel: elq.withChannel.Clone(),
		withChain:   elq.withChain.Clone(),
		withEvents:  elq.withEvents.Clone(),
		// clone intermediate query.
		sql:  elq.sql.Clone(),
		path: elq.path,
	}
}

// WithChannel tells the query-builder to eager-load the nodes that are connected to
// the "channel" edge. The optional arguments are used to configure the query builder of the edge.
func (elq *EventListenerQuery) WithChannel(opts ...func(*ChannelQuery)) *EventListenerQuery {
	query := (&ChannelClient{config: elq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	elq.withChannel = query
	return elq
}

// WithChain tells the query-builder to eager-load the nodes that are connected to
// the "chain" edge. The optional arguments are used to configure the query builder of the edge.
func (elq *EventListenerQuery) WithChain(opts ...func(*ChainQuery)) *EventListenerQuery {
	query := (&ChainClient{config: elq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	elq.withChain = query
	return elq
}

// WithEvents tells the query-builder to eager-load the nodes that are connected to
// the "events" edge. The optional arguments are used to configure the query builder of the edge.
func (elq *EventListenerQuery) WithEvents(opts ...func(*EventQuery)) *EventListenerQuery {
	query := (&EventClient{config: elq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	elq.withEvents = query
	return elq
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
//	client.EventListener.Query().
//		GroupBy(eventlistener.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (elq *EventListenerQuery) GroupBy(field string, fields ...string) *EventListenerGroupBy {
	elq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &EventListenerGroupBy{build: elq}
	grbuild.flds = &elq.ctx.Fields
	grbuild.label = eventlistener.Label
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
//	client.EventListener.Query().
//		Select(eventlistener.FieldCreateTime).
//		Scan(ctx, &v)
func (elq *EventListenerQuery) Select(fields ...string) *EventListenerSelect {
	elq.ctx.Fields = append(elq.ctx.Fields, fields...)
	sbuild := &EventListenerSelect{EventListenerQuery: elq}
	sbuild.label = eventlistener.Label
	sbuild.flds, sbuild.scan = &elq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a EventListenerSelect configured with the given aggregations.
func (elq *EventListenerQuery) Aggregate(fns ...AggregateFunc) *EventListenerSelect {
	return elq.Select().Aggregate(fns...)
}

func (elq *EventListenerQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range elq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, elq); err != nil {
				return err
			}
		}
	}
	for _, f := range elq.ctx.Fields {
		if !eventlistener.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if elq.path != nil {
		prev, err := elq.path(ctx)
		if err != nil {
			return err
		}
		elq.sql = prev
	}
	return nil
}

func (elq *EventListenerQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*EventListener, error) {
	var (
		nodes       = []*EventListener{}
		withFKs     = elq.withFKs
		_spec       = elq.querySpec()
		loadedTypes = [3]bool{
			elq.withChannel != nil,
			elq.withChain != nil,
			elq.withEvents != nil,
		}
	)
	if elq.withChannel != nil || elq.withChain != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, eventlistener.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*EventListener).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &EventListener{config: elq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, elq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := elq.withChannel; query != nil {
		if err := elq.loadChannel(ctx, query, nodes, nil,
			func(n *EventListener, e *Channel) { n.Edges.Channel = e }); err != nil {
			return nil, err
		}
	}
	if query := elq.withChain; query != nil {
		if err := elq.loadChain(ctx, query, nodes, nil,
			func(n *EventListener, e *Chain) { n.Edges.Chain = e }); err != nil {
			return nil, err
		}
	}
	if query := elq.withEvents; query != nil {
		if err := elq.loadEvents(ctx, query, nodes,
			func(n *EventListener) { n.Edges.Events = []*Event{} },
			func(n *EventListener, e *Event) { n.Edges.Events = append(n.Edges.Events, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (elq *EventListenerQuery) loadChannel(ctx context.Context, query *ChannelQuery, nodes []*EventListener, init func(*EventListener), assign func(*EventListener, *Channel)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*EventListener)
	for i := range nodes {
		if nodes[i].channel_event_listeners == nil {
			continue
		}
		fk := *nodes[i].channel_event_listeners
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(channel.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "channel_event_listeners" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (elq *EventListenerQuery) loadChain(ctx context.Context, query *ChainQuery, nodes []*EventListener, init func(*EventListener), assign func(*EventListener, *Chain)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*EventListener)
	for i := range nodes {
		if nodes[i].chain_event_listeners == nil {
			continue
		}
		fk := *nodes[i].chain_event_listeners
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(chain.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "chain_event_listeners" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (elq *EventListenerQuery) loadEvents(ctx context.Context, query *EventQuery, nodes []*EventListener, init func(*EventListener), assign func(*EventListener, *Event)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*EventListener)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.Event(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(eventlistener.EventsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.event_listener_events
		if fk == nil {
			return fmt.Errorf(`foreign-key "event_listener_events" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "event_listener_events" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (elq *EventListenerQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := elq.querySpec()
	_spec.Node.Columns = elq.ctx.Fields
	if len(elq.ctx.Fields) > 0 {
		_spec.Unique = elq.ctx.Unique != nil && *elq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, elq.driver, _spec)
}

func (elq *EventListenerQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(eventlistener.Table, eventlistener.Columns, sqlgraph.NewFieldSpec(eventlistener.FieldID, field.TypeInt))
	_spec.From = elq.sql
	if unique := elq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if elq.path != nil {
		_spec.Unique = true
	}
	if fields := elq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, eventlistener.FieldID)
		for i := range fields {
			if fields[i] != eventlistener.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := elq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := elq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := elq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := elq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (elq *EventListenerQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(elq.driver.Dialect())
	t1 := builder.Table(eventlistener.Table)
	columns := elq.ctx.Fields
	if len(columns) == 0 {
		columns = eventlistener.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if elq.sql != nil {
		selector = elq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if elq.ctx.Unique != nil && *elq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range elq.predicates {
		p(selector)
	}
	for _, p := range elq.order {
		p(selector)
	}
	if offset := elq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := elq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// EventListenerGroupBy is the group-by builder for EventListener entities.
type EventListenerGroupBy struct {
	selector
	build *EventListenerQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (elgb *EventListenerGroupBy) Aggregate(fns ...AggregateFunc) *EventListenerGroupBy {
	elgb.fns = append(elgb.fns, fns...)
	return elgb
}

// Scan applies the selector query and scans the result into the given value.
func (elgb *EventListenerGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, elgb.build.ctx, "GroupBy")
	if err := elgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*EventListenerQuery, *EventListenerGroupBy](ctx, elgb.build, elgb, elgb.build.inters, v)
}

func (elgb *EventListenerGroupBy) sqlScan(ctx context.Context, root *EventListenerQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(elgb.fns))
	for _, fn := range elgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*elgb.flds)+len(elgb.fns))
		for _, f := range *elgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*elgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := elgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// EventListenerSelect is the builder for selecting fields of EventListener entities.
type EventListenerSelect struct {
	*EventListenerQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (els *EventListenerSelect) Aggregate(fns ...AggregateFunc) *EventListenerSelect {
	els.fns = append(els.fns, fns...)
	return els
}

// Scan applies the selector query and scans the result into the given value.
func (els *EventListenerSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, els.ctx, "Select")
	if err := els.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*EventListenerQuery, *EventListenerSelect](ctx, els.EventListenerQuery, els, els.inters, v)
}

func (els *EventListenerSelect) sqlScan(ctx context.Context, root *EventListenerQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(els.fns))
	for _, fn := range els.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*els.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := els.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
