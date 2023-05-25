package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// Proposal holds the schema definition for the Proposal entity.
type Proposal struct {
	ent.Schema
}

func (Proposal) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Fields of the Proposal.
func (Proposal) Fields() []ent.Field {
	var statusValues []string
	for _, status := range v1beta1.ProposalStatus_name {
		statusValues = append(statusValues, status)
	}
	return []ent.Field{
		field.Uint64("proposal_id"),
		field.String("title"),
		field.String("description"),
		field.Time("voting_start_time"),
		field.Time("voting_end_time"),
		field.Enum("status").
			Values(statusValues...),
	}
}

// Edges of the Proposal.
func (Proposal) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("chain", Chain.Type).
			Ref("proposals").
			Unique(),
	}
}
