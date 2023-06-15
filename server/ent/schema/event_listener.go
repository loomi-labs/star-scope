package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// EventListener holds the schema definition for the EventListener entity.
type EventListener struct {
	ent.Schema
}

func (EventListener) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Fields of the EventListener.
func (EventListener) Fields() []ent.Field {
	return []ent.Field{
		field.String("wallet_address").
			Optional(),
		field.Enum("data_type").
			Values(getDataTypes()...),
	}
}

// Edges of the EventListener.
func (EventListener) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("users", User.Type).
			Ref("event_listeners"),
		edge.From("chain", Chain.Type).
			Ref("event_listeners").
			Unique(),
		edge.To("events", Event.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.From("comm_channels", CommChannel.Type).
			Ref("event_listeners"),
	}
}
