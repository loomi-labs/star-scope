// Code generated by ent, DO NOT EDIT.

package eventlistener

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/loomi-labs/star-scope/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.EventListener {
	return predicate.EventListener(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.EventListener {
	return predicate.EventListener(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.EventListener {
	return predicate.EventListener(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.EventListener {
	return predicate.EventListener(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.EventListener {
	return predicate.EventListener(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.EventListener {
	return predicate.EventListener(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.EventListener {
	return predicate.EventListener(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.EventListener {
	return predicate.EventListener(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.EventListener {
	return predicate.EventListener(sql.FieldLTE(FieldID, id))
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.EventListener {
	return predicate.EventListener(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.EventListener {
	return predicate.EventListener(sql.FieldEQ(FieldUpdateTime, v))
}

// WalletAddress applies equality check predicate on the "wallet_address" field. It's identical to WalletAddressEQ.
func WalletAddress(v string) predicate.EventListener {
	return predicate.EventListener(sql.FieldEQ(FieldWalletAddress, v))
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.EventListener {
	return predicate.EventListener(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.EventListener {
	return predicate.EventListener(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.EventListener {
	return predicate.EventListener(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "create_time" field.
func CreateTimeNotIn(vs ...time.Time) predicate.EventListener {
	return predicate.EventListener(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "create_time" field.
func CreateTimeGT(v time.Time) predicate.EventListener {
	return predicate.EventListener(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.EventListener {
	return predicate.EventListener(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.EventListener {
	return predicate.EventListener(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.EventListener {
	return predicate.EventListener(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.EventListener {
	return predicate.EventListener(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.EventListener {
	return predicate.EventListener(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.EventListener {
	return predicate.EventListener(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "update_time" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.EventListener {
	return predicate.EventListener(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "update_time" field.
func UpdateTimeGT(v time.Time) predicate.EventListener {
	return predicate.EventListener(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.EventListener {
	return predicate.EventListener(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.EventListener {
	return predicate.EventListener(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.EventListener {
	return predicate.EventListener(sql.FieldLTE(FieldUpdateTime, v))
}

// WalletAddressEQ applies the EQ predicate on the "wallet_address" field.
func WalletAddressEQ(v string) predicate.EventListener {
	return predicate.EventListener(sql.FieldEQ(FieldWalletAddress, v))
}

// WalletAddressNEQ applies the NEQ predicate on the "wallet_address" field.
func WalletAddressNEQ(v string) predicate.EventListener {
	return predicate.EventListener(sql.FieldNEQ(FieldWalletAddress, v))
}

// WalletAddressIn applies the In predicate on the "wallet_address" field.
func WalletAddressIn(vs ...string) predicate.EventListener {
	return predicate.EventListener(sql.FieldIn(FieldWalletAddress, vs...))
}

// WalletAddressNotIn applies the NotIn predicate on the "wallet_address" field.
func WalletAddressNotIn(vs ...string) predicate.EventListener {
	return predicate.EventListener(sql.FieldNotIn(FieldWalletAddress, vs...))
}

// WalletAddressGT applies the GT predicate on the "wallet_address" field.
func WalletAddressGT(v string) predicate.EventListener {
	return predicate.EventListener(sql.FieldGT(FieldWalletAddress, v))
}

// WalletAddressGTE applies the GTE predicate on the "wallet_address" field.
func WalletAddressGTE(v string) predicate.EventListener {
	return predicate.EventListener(sql.FieldGTE(FieldWalletAddress, v))
}

// WalletAddressLT applies the LT predicate on the "wallet_address" field.
func WalletAddressLT(v string) predicate.EventListener {
	return predicate.EventListener(sql.FieldLT(FieldWalletAddress, v))
}

// WalletAddressLTE applies the LTE predicate on the "wallet_address" field.
func WalletAddressLTE(v string) predicate.EventListener {
	return predicate.EventListener(sql.FieldLTE(FieldWalletAddress, v))
}

// WalletAddressContains applies the Contains predicate on the "wallet_address" field.
func WalletAddressContains(v string) predicate.EventListener {
	return predicate.EventListener(sql.FieldContains(FieldWalletAddress, v))
}

// WalletAddressHasPrefix applies the HasPrefix predicate on the "wallet_address" field.
func WalletAddressHasPrefix(v string) predicate.EventListener {
	return predicate.EventListener(sql.FieldHasPrefix(FieldWalletAddress, v))
}

// WalletAddressHasSuffix applies the HasSuffix predicate on the "wallet_address" field.
func WalletAddressHasSuffix(v string) predicate.EventListener {
	return predicate.EventListener(sql.FieldHasSuffix(FieldWalletAddress, v))
}

// WalletAddressEqualFold applies the EqualFold predicate on the "wallet_address" field.
func WalletAddressEqualFold(v string) predicate.EventListener {
	return predicate.EventListener(sql.FieldEqualFold(FieldWalletAddress, v))
}

// WalletAddressContainsFold applies the ContainsFold predicate on the "wallet_address" field.
func WalletAddressContainsFold(v string) predicate.EventListener {
	return predicate.EventListener(sql.FieldContainsFold(FieldWalletAddress, v))
}

// HasChannel applies the HasEdge predicate on the "channel" edge.
func HasChannel() predicate.EventListener {
	return predicate.EventListener(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ChannelTable, ChannelColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasChannelWith applies the HasEdge predicate on the "channel" edge with a given conditions (other predicates).
func HasChannelWith(preds ...predicate.Channel) predicate.EventListener {
	return predicate.EventListener(func(s *sql.Selector) {
		step := newChannelStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasChain applies the HasEdge predicate on the "chain" edge.
func HasChain() predicate.EventListener {
	return predicate.EventListener(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ChainTable, ChainColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasChainWith applies the HasEdge predicate on the "chain" edge with a given conditions (other predicates).
func HasChainWith(preds ...predicate.Chain) predicate.EventListener {
	return predicate.EventListener(func(s *sql.Selector) {
		step := newChainStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasEvents applies the HasEdge predicate on the "events" edge.
func HasEvents() predicate.EventListener {
	return predicate.EventListener(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, EventsTable, EventsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEventsWith applies the HasEdge predicate on the "events" edge with a given conditions (other predicates).
func HasEventsWith(preds ...predicate.Event) predicate.EventListener {
	return predicate.EventListener(func(s *sql.Selector) {
		step := newEventsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.EventListener) predicate.EventListener {
	return predicate.EventListener(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.EventListener) predicate.EventListener {
	return predicate.EventListener(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.EventListener) predicate.EventListener {
	return predicate.EventListener(func(s *sql.Selector) {
		p(s.Not())
	})
}
