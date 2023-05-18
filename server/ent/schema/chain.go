package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// Chain holds the schema definition for the Chain entity.
type Chain struct {
	ent.Schema
}

func (Chain) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Fields of the Chain.
func (Chain) Fields() []ent.Field {
	return []ent.Field{
		field.String("chain_id"),
		field.String("name").
			Unique(),
		field.String("pretty_name").
			Unique(),
		field.String("path").
			Unique().
			Immutable(),
		field.String("image"),
		field.String("bech32_prefix"),

		field.Uint64("indexing_height").
			Default(0),
		field.Bool("has_custom_indexer").
			Default(false),

		// comma separated list of message types that are handled by the indexer
		field.String("handled_message_types").
			Default(""),
		// comma separated list of message types that are not handled by the indexer
		field.String("unhandled_message_types").
			Default(""),

		field.Bool("is_enabled").
			Default(false),
	}
}

// Edges of the Chain.
func (Chain) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("event_listeners", EventListener.Type),
	}
}

func (Chain) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Unique(),
	}
}
