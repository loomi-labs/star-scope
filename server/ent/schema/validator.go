package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/loomi-labs/star-scope/common"
)

// Validator holds the schema definition for the Validator entity.
type Validator struct {
	ent.Schema
}

func (Validator) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Fields of the Validator.
func (Validator) Fields() []ent.Field {
	return []ent.Field{
		field.String("operator_address").
			Immutable().
			Validate(func(s string) error {
				return common.ValidateBech32Address(s)
			}),
		field.String("address").
			Immutable().
			Validate(func(s string) error {
				return common.ValidateBech32Address(s)
			}),
		field.String("moniker"),
		field.Time("first_inactive_time").
			Nillable().
			Optional(),
		field.Uint64("last_slash_validator_period").
			Nillable().
			Optional(),
	}
}

// Edges of the Validator.
func (Validator) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("chain", Chain.Type).
			Ref("validators").
			Unique().
			Required(),
		edge.To("selected_by_setups", UserSetup.Type),
	}
}

func (Validator) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("operator_address"),
		index.Fields("address"),
		index.Fields("moniker"),
		index.Fields("moniker", "operator_address").
			Edges("chain").
			Unique(),
		index.Fields("moniker", "address").
			Edges("chain").
			Unique(),
		index.Fields("address"). // address is unique per chain but not globally (ex: terra and terra2)
						Edges("chain").
						Unique(),
		index.Fields("operator_address"). // same as address
							Edges("chain").
							Unique(),
	}
}
