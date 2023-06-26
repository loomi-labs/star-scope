package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// UserSetup holds the schema definition for the UserSetup entity.
type UserSetup struct {
	ent.Schema
}

func (UserSetup) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Fields of the UserSetup.
func (UserSetup) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("step").
			Values("one", "two", "three", "four", "five").
			Default("one"),
		field.Bool("is_validator").
			Default(false),
		field.Strings("wallet_addresses").
			Optional(),
		field.Bool("notify_funding").
			Default(false),
		field.Bool("notify_staking").
			Default(false),
		field.Bool("notify_gov_new_proposal").
			Default(false),
		field.Bool("notify_gov_voting_end").
			Default(false),
		field.Bool("notify_gov_voting_reminder").
			Default(false),
	}
}

// Edges of the UserSetup.
func (UserSetup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("setup").
			Unique().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("selected_validators", Validator.Type),
		edge.To("selected_chains", Chain.Type),
	}
}
