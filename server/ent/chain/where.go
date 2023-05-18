// Code generated by ent, DO NOT EDIT.

package chain

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/loomi-labs/star-scope/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Chain {
	return predicate.Chain(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Chain {
	return predicate.Chain(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Chain {
	return predicate.Chain(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Chain {
	return predicate.Chain(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Chain {
	return predicate.Chain(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Chain {
	return predicate.Chain(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Chain {
	return predicate.Chain(sql.FieldLTE(FieldID, id))
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldUpdateTime, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldName, v))
}

// Image applies equality check predicate on the "image" field. It's identical to ImageEQ.
func Image(v string) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldImage, v))
}

// IndexingHeight applies equality check predicate on the "indexing_height" field. It's identical to IndexingHeightEQ.
func IndexingHeight(v uint64) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldIndexingHeight, v))
}

// Path applies equality check predicate on the "path" field. It's identical to PathEQ.
func Path(v string) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldPath, v))
}

// HasCustomIndexer applies equality check predicate on the "has_custom_indexer" field. It's identical to HasCustomIndexerEQ.
func HasCustomIndexer(v bool) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldHasCustomIndexer, v))
}

// HandledMessageTypes applies equality check predicate on the "handled_message_types" field. It's identical to HandledMessageTypesEQ.
func HandledMessageTypes(v string) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldHandledMessageTypes, v))
}

// UnhandledMessageTypes applies equality check predicate on the "unhandled_message_types" field. It's identical to UnhandledMessageTypesEQ.
func UnhandledMessageTypes(v string) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldUnhandledMessageTypes, v))
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.Chain {
	return predicate.Chain(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.Chain {
	return predicate.Chain(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "create_time" field.
func CreateTimeNotIn(vs ...time.Time) predicate.Chain {
	return predicate.Chain(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "create_time" field.
func CreateTimeGT(v time.Time) predicate.Chain {
	return predicate.Chain(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.Chain {
	return predicate.Chain(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.Chain {
	return predicate.Chain(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.Chain {
	return predicate.Chain(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.Chain {
	return predicate.Chain(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.Chain {
	return predicate.Chain(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "update_time" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.Chain {
	return predicate.Chain(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "update_time" field.
func UpdateTimeGT(v time.Time) predicate.Chain {
	return predicate.Chain(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.Chain {
	return predicate.Chain(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.Chain {
	return predicate.Chain(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.Chain {
	return predicate.Chain(sql.FieldLTE(FieldUpdateTime, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Chain {
	return predicate.Chain(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Chain {
	return predicate.Chain(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Chain {
	return predicate.Chain(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Chain {
	return predicate.Chain(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Chain {
	return predicate.Chain(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Chain {
	return predicate.Chain(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Chain {
	return predicate.Chain(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Chain {
	return predicate.Chain(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Chain {
	return predicate.Chain(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Chain {
	return predicate.Chain(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Chain {
	return predicate.Chain(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Chain {
	return predicate.Chain(sql.FieldContainsFold(FieldName, v))
}

// ImageEQ applies the EQ predicate on the "image" field.
func ImageEQ(v string) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldImage, v))
}

// ImageNEQ applies the NEQ predicate on the "image" field.
func ImageNEQ(v string) predicate.Chain {
	return predicate.Chain(sql.FieldNEQ(FieldImage, v))
}

// ImageIn applies the In predicate on the "image" field.
func ImageIn(vs ...string) predicate.Chain {
	return predicate.Chain(sql.FieldIn(FieldImage, vs...))
}

// ImageNotIn applies the NotIn predicate on the "image" field.
func ImageNotIn(vs ...string) predicate.Chain {
	return predicate.Chain(sql.FieldNotIn(FieldImage, vs...))
}

// ImageGT applies the GT predicate on the "image" field.
func ImageGT(v string) predicate.Chain {
	return predicate.Chain(sql.FieldGT(FieldImage, v))
}

// ImageGTE applies the GTE predicate on the "image" field.
func ImageGTE(v string) predicate.Chain {
	return predicate.Chain(sql.FieldGTE(FieldImage, v))
}

// ImageLT applies the LT predicate on the "image" field.
func ImageLT(v string) predicate.Chain {
	return predicate.Chain(sql.FieldLT(FieldImage, v))
}

// ImageLTE applies the LTE predicate on the "image" field.
func ImageLTE(v string) predicate.Chain {
	return predicate.Chain(sql.FieldLTE(FieldImage, v))
}

// ImageContains applies the Contains predicate on the "image" field.
func ImageContains(v string) predicate.Chain {
	return predicate.Chain(sql.FieldContains(FieldImage, v))
}

// ImageHasPrefix applies the HasPrefix predicate on the "image" field.
func ImageHasPrefix(v string) predicate.Chain {
	return predicate.Chain(sql.FieldHasPrefix(FieldImage, v))
}

// ImageHasSuffix applies the HasSuffix predicate on the "image" field.
func ImageHasSuffix(v string) predicate.Chain {
	return predicate.Chain(sql.FieldHasSuffix(FieldImage, v))
}

// ImageEqualFold applies the EqualFold predicate on the "image" field.
func ImageEqualFold(v string) predicate.Chain {
	return predicate.Chain(sql.FieldEqualFold(FieldImage, v))
}

// ImageContainsFold applies the ContainsFold predicate on the "image" field.
func ImageContainsFold(v string) predicate.Chain {
	return predicate.Chain(sql.FieldContainsFold(FieldImage, v))
}

// IndexingHeightEQ applies the EQ predicate on the "indexing_height" field.
func IndexingHeightEQ(v uint64) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldIndexingHeight, v))
}

// IndexingHeightNEQ applies the NEQ predicate on the "indexing_height" field.
func IndexingHeightNEQ(v uint64) predicate.Chain {
	return predicate.Chain(sql.FieldNEQ(FieldIndexingHeight, v))
}

// IndexingHeightIn applies the In predicate on the "indexing_height" field.
func IndexingHeightIn(vs ...uint64) predicate.Chain {
	return predicate.Chain(sql.FieldIn(FieldIndexingHeight, vs...))
}

// IndexingHeightNotIn applies the NotIn predicate on the "indexing_height" field.
func IndexingHeightNotIn(vs ...uint64) predicate.Chain {
	return predicate.Chain(sql.FieldNotIn(FieldIndexingHeight, vs...))
}

// IndexingHeightGT applies the GT predicate on the "indexing_height" field.
func IndexingHeightGT(v uint64) predicate.Chain {
	return predicate.Chain(sql.FieldGT(FieldIndexingHeight, v))
}

// IndexingHeightGTE applies the GTE predicate on the "indexing_height" field.
func IndexingHeightGTE(v uint64) predicate.Chain {
	return predicate.Chain(sql.FieldGTE(FieldIndexingHeight, v))
}

// IndexingHeightLT applies the LT predicate on the "indexing_height" field.
func IndexingHeightLT(v uint64) predicate.Chain {
	return predicate.Chain(sql.FieldLT(FieldIndexingHeight, v))
}

// IndexingHeightLTE applies the LTE predicate on the "indexing_height" field.
func IndexingHeightLTE(v uint64) predicate.Chain {
	return predicate.Chain(sql.FieldLTE(FieldIndexingHeight, v))
}

// PathEQ applies the EQ predicate on the "path" field.
func PathEQ(v string) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldPath, v))
}

// PathNEQ applies the NEQ predicate on the "path" field.
func PathNEQ(v string) predicate.Chain {
	return predicate.Chain(sql.FieldNEQ(FieldPath, v))
}

// PathIn applies the In predicate on the "path" field.
func PathIn(vs ...string) predicate.Chain {
	return predicate.Chain(sql.FieldIn(FieldPath, vs...))
}

// PathNotIn applies the NotIn predicate on the "path" field.
func PathNotIn(vs ...string) predicate.Chain {
	return predicate.Chain(sql.FieldNotIn(FieldPath, vs...))
}

// PathGT applies the GT predicate on the "path" field.
func PathGT(v string) predicate.Chain {
	return predicate.Chain(sql.FieldGT(FieldPath, v))
}

// PathGTE applies the GTE predicate on the "path" field.
func PathGTE(v string) predicate.Chain {
	return predicate.Chain(sql.FieldGTE(FieldPath, v))
}

// PathLT applies the LT predicate on the "path" field.
func PathLT(v string) predicate.Chain {
	return predicate.Chain(sql.FieldLT(FieldPath, v))
}

// PathLTE applies the LTE predicate on the "path" field.
func PathLTE(v string) predicate.Chain {
	return predicate.Chain(sql.FieldLTE(FieldPath, v))
}

// PathContains applies the Contains predicate on the "path" field.
func PathContains(v string) predicate.Chain {
	return predicate.Chain(sql.FieldContains(FieldPath, v))
}

// PathHasPrefix applies the HasPrefix predicate on the "path" field.
func PathHasPrefix(v string) predicate.Chain {
	return predicate.Chain(sql.FieldHasPrefix(FieldPath, v))
}

// PathHasSuffix applies the HasSuffix predicate on the "path" field.
func PathHasSuffix(v string) predicate.Chain {
	return predicate.Chain(sql.FieldHasSuffix(FieldPath, v))
}

// PathEqualFold applies the EqualFold predicate on the "path" field.
func PathEqualFold(v string) predicate.Chain {
	return predicate.Chain(sql.FieldEqualFold(FieldPath, v))
}

// PathContainsFold applies the ContainsFold predicate on the "path" field.
func PathContainsFold(v string) predicate.Chain {
	return predicate.Chain(sql.FieldContainsFold(FieldPath, v))
}

// HasCustomIndexerEQ applies the EQ predicate on the "has_custom_indexer" field.
func HasCustomIndexerEQ(v bool) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldHasCustomIndexer, v))
}

// HasCustomIndexerNEQ applies the NEQ predicate on the "has_custom_indexer" field.
func HasCustomIndexerNEQ(v bool) predicate.Chain {
	return predicate.Chain(sql.FieldNEQ(FieldHasCustomIndexer, v))
}

// HandledMessageTypesEQ applies the EQ predicate on the "handled_message_types" field.
func HandledMessageTypesEQ(v string) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldHandledMessageTypes, v))
}

// HandledMessageTypesNEQ applies the NEQ predicate on the "handled_message_types" field.
func HandledMessageTypesNEQ(v string) predicate.Chain {
	return predicate.Chain(sql.FieldNEQ(FieldHandledMessageTypes, v))
}

// HandledMessageTypesIn applies the In predicate on the "handled_message_types" field.
func HandledMessageTypesIn(vs ...string) predicate.Chain {
	return predicate.Chain(sql.FieldIn(FieldHandledMessageTypes, vs...))
}

// HandledMessageTypesNotIn applies the NotIn predicate on the "handled_message_types" field.
func HandledMessageTypesNotIn(vs ...string) predicate.Chain {
	return predicate.Chain(sql.FieldNotIn(FieldHandledMessageTypes, vs...))
}

// HandledMessageTypesGT applies the GT predicate on the "handled_message_types" field.
func HandledMessageTypesGT(v string) predicate.Chain {
	return predicate.Chain(sql.FieldGT(FieldHandledMessageTypes, v))
}

// HandledMessageTypesGTE applies the GTE predicate on the "handled_message_types" field.
func HandledMessageTypesGTE(v string) predicate.Chain {
	return predicate.Chain(sql.FieldGTE(FieldHandledMessageTypes, v))
}

// HandledMessageTypesLT applies the LT predicate on the "handled_message_types" field.
func HandledMessageTypesLT(v string) predicate.Chain {
	return predicate.Chain(sql.FieldLT(FieldHandledMessageTypes, v))
}

// HandledMessageTypesLTE applies the LTE predicate on the "handled_message_types" field.
func HandledMessageTypesLTE(v string) predicate.Chain {
	return predicate.Chain(sql.FieldLTE(FieldHandledMessageTypes, v))
}

// HandledMessageTypesContains applies the Contains predicate on the "handled_message_types" field.
func HandledMessageTypesContains(v string) predicate.Chain {
	return predicate.Chain(sql.FieldContains(FieldHandledMessageTypes, v))
}

// HandledMessageTypesHasPrefix applies the HasPrefix predicate on the "handled_message_types" field.
func HandledMessageTypesHasPrefix(v string) predicate.Chain {
	return predicate.Chain(sql.FieldHasPrefix(FieldHandledMessageTypes, v))
}

// HandledMessageTypesHasSuffix applies the HasSuffix predicate on the "handled_message_types" field.
func HandledMessageTypesHasSuffix(v string) predicate.Chain {
	return predicate.Chain(sql.FieldHasSuffix(FieldHandledMessageTypes, v))
}

// HandledMessageTypesEqualFold applies the EqualFold predicate on the "handled_message_types" field.
func HandledMessageTypesEqualFold(v string) predicate.Chain {
	return predicate.Chain(sql.FieldEqualFold(FieldHandledMessageTypes, v))
}

// HandledMessageTypesContainsFold applies the ContainsFold predicate on the "handled_message_types" field.
func HandledMessageTypesContainsFold(v string) predicate.Chain {
	return predicate.Chain(sql.FieldContainsFold(FieldHandledMessageTypes, v))
}

// UnhandledMessageTypesEQ applies the EQ predicate on the "unhandled_message_types" field.
func UnhandledMessageTypesEQ(v string) predicate.Chain {
	return predicate.Chain(sql.FieldEQ(FieldUnhandledMessageTypes, v))
}

// UnhandledMessageTypesNEQ applies the NEQ predicate on the "unhandled_message_types" field.
func UnhandledMessageTypesNEQ(v string) predicate.Chain {
	return predicate.Chain(sql.FieldNEQ(FieldUnhandledMessageTypes, v))
}

// UnhandledMessageTypesIn applies the In predicate on the "unhandled_message_types" field.
func UnhandledMessageTypesIn(vs ...string) predicate.Chain {
	return predicate.Chain(sql.FieldIn(FieldUnhandledMessageTypes, vs...))
}

// UnhandledMessageTypesNotIn applies the NotIn predicate on the "unhandled_message_types" field.
func UnhandledMessageTypesNotIn(vs ...string) predicate.Chain {
	return predicate.Chain(sql.FieldNotIn(FieldUnhandledMessageTypes, vs...))
}

// UnhandledMessageTypesGT applies the GT predicate on the "unhandled_message_types" field.
func UnhandledMessageTypesGT(v string) predicate.Chain {
	return predicate.Chain(sql.FieldGT(FieldUnhandledMessageTypes, v))
}

// UnhandledMessageTypesGTE applies the GTE predicate on the "unhandled_message_types" field.
func UnhandledMessageTypesGTE(v string) predicate.Chain {
	return predicate.Chain(sql.FieldGTE(FieldUnhandledMessageTypes, v))
}

// UnhandledMessageTypesLT applies the LT predicate on the "unhandled_message_types" field.
func UnhandledMessageTypesLT(v string) predicate.Chain {
	return predicate.Chain(sql.FieldLT(FieldUnhandledMessageTypes, v))
}

// UnhandledMessageTypesLTE applies the LTE predicate on the "unhandled_message_types" field.
func UnhandledMessageTypesLTE(v string) predicate.Chain {
	return predicate.Chain(sql.FieldLTE(FieldUnhandledMessageTypes, v))
}

// UnhandledMessageTypesContains applies the Contains predicate on the "unhandled_message_types" field.
func UnhandledMessageTypesContains(v string) predicate.Chain {
	return predicate.Chain(sql.FieldContains(FieldUnhandledMessageTypes, v))
}

// UnhandledMessageTypesHasPrefix applies the HasPrefix predicate on the "unhandled_message_types" field.
func UnhandledMessageTypesHasPrefix(v string) predicate.Chain {
	return predicate.Chain(sql.FieldHasPrefix(FieldUnhandledMessageTypes, v))
}

// UnhandledMessageTypesHasSuffix applies the HasSuffix predicate on the "unhandled_message_types" field.
func UnhandledMessageTypesHasSuffix(v string) predicate.Chain {
	return predicate.Chain(sql.FieldHasSuffix(FieldUnhandledMessageTypes, v))
}

// UnhandledMessageTypesEqualFold applies the EqualFold predicate on the "unhandled_message_types" field.
func UnhandledMessageTypesEqualFold(v string) predicate.Chain {
	return predicate.Chain(sql.FieldEqualFold(FieldUnhandledMessageTypes, v))
}

// UnhandledMessageTypesContainsFold applies the ContainsFold predicate on the "unhandled_message_types" field.
func UnhandledMessageTypesContainsFold(v string) predicate.Chain {
	return predicate.Chain(sql.FieldContainsFold(FieldUnhandledMessageTypes, v))
}

// HasEventListeners applies the HasEdge predicate on the "event_listeners" edge.
func HasEventListeners() predicate.Chain {
	return predicate.Chain(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, EventListenersTable, EventListenersColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEventListenersWith applies the HasEdge predicate on the "event_listeners" edge with a given conditions (other predicates).
func HasEventListenersWith(preds ...predicate.EventListener) predicate.Chain {
	return predicate.Chain(func(s *sql.Selector) {
		step := newEventListenersStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Chain) predicate.Chain {
	return predicate.Chain(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Chain) predicate.Chain {
	return predicate.Chain(func(s *sql.Selector) {
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
func Not(p predicate.Chain) predicate.Chain {
	return predicate.Chain(func(s *sql.Selector) {
		p(s.Not())
	})
}
