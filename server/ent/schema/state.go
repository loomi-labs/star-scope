package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// State holds the schema definition for the State entity.
type State struct {
	ent.Schema
}

func (State) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Fields of the State.
func (State) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("entity").
			Values("discord", "telegram"),
		field.Time("last_event_time").
			Optional(),
	}
}

// Edges of the State.
func (State) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (State) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("entity").
			Unique(),
	}
}
