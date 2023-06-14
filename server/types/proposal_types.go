package types

import (
	"encoding/json"
	cosmossdktypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/loomi-labs/star-scope/event"
	"strings"
	"time"
)

type ProposalStatus event.ProposalStatus

func (s *ProposalStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(event.ProposalStatus(*s).String())
}

func (s *ProposalStatus) UnmarshalJSON(data []byte) error {
	var name = ""
	err := json.Unmarshal(data, &name)
	if err != nil {
		return err
	}
	*s = ProposalStatus(event.ProposalStatus_value[strings.ToUpper(name)])
	return nil
}

func (s *ProposalStatus) String() string {
	return event.ProposalStatus(*s).String()
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

type ChainProposalVoteOption cosmossdktypes.VoteOption

func (o *ChainProposalVoteOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(cosmossdktypes.VoteOption(*o).String())
}

func (o *ChainProposalVoteOption) UnmarshalJSON(data []byte) error {
	var name string
	err := json.Unmarshal(data, &name)
	if err != nil {
		return err
	}
	*o = ChainProposalVoteOption(cosmossdktypes.VoteOption_value[name])
	return nil
}

func (o ChainProposalVoteOption) ToCosmosType() cosmossdktypes.VoteOption {
	return cosmossdktypes.VoteOption(o)
}

type ChainProposalVoteResponse struct {
	Vote struct {
		ProposalID string                  `json:"proposal_id"`
		Voter      string                  `json:"voter"`
		Option     ChainProposalVoteOption `json:"option"`
		Options    []struct {
			Option ChainProposalVoteOption `json:"option"`
			Weight string                  `json:"weight"`
		} `json:"options"`
	} `json:"vote"`
}
