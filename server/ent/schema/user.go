package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("role").
			Values("user", "admin").
			Default("user"),
		field.Int64("telegram_user_id").
			Unique().
			Optional(),
		field.String("telegram_username").
			Optional(),
		field.Int64("discord_user_id").
			Unique().
			Optional(),
		field.String("discord_username").
			Optional(),
		field.String("wallet_address").
			Unique().
			Optional(),
		field.Time("last_login_time").
			Optional().
			Nillable(),
		field.Bool("is_setup_complete").
			Default(false),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("event_listeners", EventListener.Type),
		edge.To("comm_channels", CommChannel.Type),
		edge.To("setup", UserSetup.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("telegram_user_id"),
		index.Fields("discord_user_id"),
		index.Fields("wallet_address"),
	}
}
