// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/loomi-labs/star-scope/ent/chain"
	"github.com/loomi-labs/star-scope/ent/contractproposal"
)

// ContractProposal is the model entity for the ContractProposal schema.
type ContractProposal struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// ProposalID holds the value of the "proposal_id" field.
	ProposalID uint64 `json:"proposal_id,omitempty"`
	// Title holds the value of the "title" field.
	Title string `json:"title,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// FirstSeenTime holds the value of the "first_seen_time" field.
	FirstSeenTime time.Time `json:"first_seen_time,omitempty"`
	// VotingEndTime holds the value of the "voting_end_time" field.
	VotingEndTime time.Time `json:"voting_end_time,omitempty"`
	// ContractAddress holds the value of the "contract_address" field.
	ContractAddress string `json:"contract_address,omitempty"`
	// Status holds the value of the "status" field.
	Status contractproposal.Status `json:"status,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ContractProposalQuery when eager-loading is set.
	Edges                    ContractProposalEdges `json:"edges"`
	chain_contract_proposals *int
	selectValues             sql.SelectValues
}

// ContractProposalEdges holds the relations/edges for other nodes in the graph.
type ContractProposalEdges struct {
	// Chain holds the value of the chain edge.
	Chain *Chain `json:"chain,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// ChainOrErr returns the Chain value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ContractProposalEdges) ChainOrErr() (*Chain, error) {
	if e.loadedTypes[0] {
		if e.Chain == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: chain.Label}
		}
		return e.Chain, nil
	}
	return nil, &NotLoadedError{edge: "chain"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ContractProposal) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case contractproposal.FieldID, contractproposal.FieldProposalID:
			values[i] = new(sql.NullInt64)
		case contractproposal.FieldTitle, contractproposal.FieldDescription, contractproposal.FieldContractAddress, contractproposal.FieldStatus:
			values[i] = new(sql.NullString)
		case contractproposal.FieldCreateTime, contractproposal.FieldUpdateTime, contractproposal.FieldFirstSeenTime, contractproposal.FieldVotingEndTime:
			values[i] = new(sql.NullTime)
		case contractproposal.ForeignKeys[0]: // chain_contract_proposals
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ContractProposal fields.
func (cp *ContractProposal) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case contractproposal.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			cp.ID = int(value.Int64)
		case contractproposal.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field create_time", values[i])
			} else if value.Valid {
				cp.CreateTime = value.Time
			}
		case contractproposal.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field update_time", values[i])
			} else if value.Valid {
				cp.UpdateTime = value.Time
			}
		case contractproposal.FieldProposalID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field proposal_id", values[i])
			} else if value.Valid {
				cp.ProposalID = uint64(value.Int64)
			}
		case contractproposal.FieldTitle:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field title", values[i])
			} else if value.Valid {
				cp.Title = value.String
			}
		case contractproposal.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				cp.Description = value.String
			}
		case contractproposal.FieldFirstSeenTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field first_seen_time", values[i])
			} else if value.Valid {
				cp.FirstSeenTime = value.Time
			}
		case contractproposal.FieldVotingEndTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field voting_end_time", values[i])
			} else if value.Valid {
				cp.VotingEndTime = value.Time
			}
		case contractproposal.FieldContractAddress:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field contract_address", values[i])
			} else if value.Valid {
				cp.ContractAddress = value.String
			}
		case contractproposal.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				cp.Status = contractproposal.Status(value.String)
			}
		case contractproposal.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field chain_contract_proposals", value)
			} else if value.Valid {
				cp.chain_contract_proposals = new(int)
				*cp.chain_contract_proposals = int(value.Int64)
			}
		default:
			cp.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ContractProposal.
// This includes values selected through modifiers, order, etc.
func (cp *ContractProposal) Value(name string) (ent.Value, error) {
	return cp.selectValues.Get(name)
}

// QueryChain queries the "chain" edge of the ContractProposal entity.
func (cp *ContractProposal) QueryChain() *ChainQuery {
	return NewContractProposalClient(cp.config).QueryChain(cp)
}

// Update returns a builder for updating this ContractProposal.
// Note that you need to call ContractProposal.Unwrap() before calling this method if this ContractProposal
// was returned from a transaction, and the transaction was committed or rolled back.
func (cp *ContractProposal) Update() *ContractProposalUpdateOne {
	return NewContractProposalClient(cp.config).UpdateOne(cp)
}

// Unwrap unwraps the ContractProposal entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (cp *ContractProposal) Unwrap() *ContractProposal {
	_tx, ok := cp.config.driver.(*txDriver)
	if !ok {
		panic("ent: ContractProposal is not a transactional entity")
	}
	cp.config.driver = _tx.drv
	return cp
}

// String implements the fmt.Stringer.
func (cp *ContractProposal) String() string {
	var builder strings.Builder
	builder.WriteString("ContractProposal(")
	builder.WriteString(fmt.Sprintf("id=%v, ", cp.ID))
	builder.WriteString("create_time=")
	builder.WriteString(cp.CreateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("update_time=")
	builder.WriteString(cp.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("proposal_id=")
	builder.WriteString(fmt.Sprintf("%v", cp.ProposalID))
	builder.WriteString(", ")
	builder.WriteString("title=")
	builder.WriteString(cp.Title)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(cp.Description)
	builder.WriteString(", ")
	builder.WriteString("first_seen_time=")
	builder.WriteString(cp.FirstSeenTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("voting_end_time=")
	builder.WriteString(cp.VotingEndTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("contract_address=")
	builder.WriteString(cp.ContractAddress)
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", cp.Status))
	builder.WriteByte(')')
	return builder.String()
}

// ContractProposals is a parsable slice of ContractProposal.
type ContractProposals []*ContractProposal
