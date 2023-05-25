package types

import (
	"encoding/json"
	cosmossdktypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"time"
)

type ProposalStatus cosmossdktypes.ProposalStatus

func (s *ProposalStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(cosmossdktypes.ProposalStatus(*s).String())
}

func (s *ProposalStatus) UnmarshalJSON(data []byte) error {
	var name = ""
	err := json.Unmarshal(data, &name)
	if err != nil {
		return err
	}
	*s = ProposalStatus(cosmossdktypes.ProposalStatus_value[name])
	return nil
}

type Content struct {
	Type        string `json:"@type"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Proposal struct {
	ProposalId      int            `json:"proposal_id,string"`
	Content         Content        `json:"content"`
	Status          ProposalStatus `json:"status"`
	VotingStartTime time.Time      `json:"voting_start_time"`
	VotingEndTime   time.Time      `json:"voting_end_time"`
}

type Pagination struct {
	TotalCount int `json:"total_count"`
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
	NextKey    int `json:"next_key"`
}

type ProposalResponse struct {
	Proposal Proposal `json:"proposal"`
}

type ProposalsResponse struct {
	Proposals []Proposal `json:"proposals"`
	//Pagination Pagination `json:"pagination"`
}
