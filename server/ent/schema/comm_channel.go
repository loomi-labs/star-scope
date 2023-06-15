package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// CommChannel holds the schema definition for the CommChannel entity.
type CommChannel struct {
	ent.Schema
}

func (CommChannel) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Fields of the CommChannel.
func (CommChannel) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Enum("type").
			Values("webpush", "telegram", "discord"),
		field.Int64("telegram_chat_id").
			Unique().
			Immutable(),
		field.Int64("discord_channel_id").
			Unique().
			Immutable(),
		field.Bool("is_group").
			Immutable(),
	}
}

// Edges of the CommChannel.
func (CommChannel) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("event_listeners", EventListener.Type),
		edge.From("users", User.Type).
			Ref("comm_channels"),
	}
}
