package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/loomi-labs/star-scope/event"
)

// ContractProposal holds the schema definition for the ContractProposal entity.
type ContractProposal struct {
	ent.Schema
}

func (ContractProposal) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Fields of the ContractProposal.
func (ContractProposal) Fields() []ent.Field {
	var statusValues []string
	for _, status := range event.ContractProposalStatus_name {
		statusValues = append(statusValues, status)
	}
	return []ent.Field{
		field.Uint64("proposal_id"),
		field.String("title"),
		field.String("description"),
		field.Time("first_seen_time"),
		field.Time("voting_end_time"),
		field.String("contract_address"),
		field.Enum("status").
			Values(statusValues...),
	}
}

// Edges of the ContractProposal.
func (ContractProposal) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("chain", Chain.Type).
			Ref("contract_proposals").
			Unique(),
	}
}
