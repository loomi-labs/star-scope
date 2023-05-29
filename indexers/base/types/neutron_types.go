package types

import (
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/queryevent"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type ContractProposalStatus queryevent.ContractProposalStatus

func (s *ContractProposalStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(queryevent.ProposalStatus(*s).String())
}

func (s *ContractProposalStatus) UnmarshalJSON(data []byte) error {
	var name = ""
	err := json.Unmarshal(data, &name)
	if err != nil {
		return err
	}
	*s = ContractProposalStatus(queryevent.ContractProposalStatus_value[strings.ToUpper(name)])
	return nil
}

type AtTime time.Time

func (t *AtTime) UnmarshalJSON(b []byte) error {
	r := strings.Trim(string(b), "\"")
	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(0, q)
	return nil
}

func (t AtTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

type ConfigResponse struct {
	Data struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		DaoURI      any    `json:"dao_uri"`
	} `json:"data"`
}

type ProposalModulesResponse struct {
	Data []struct {
		Address string `json:"address"`
		Prefix  string `json:"prefix"`
		Status  string `json:"status"`
	} `json:"data"`
}

type ContractProposal struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Proposer        string `json:"proposer"`
	StartHeight     int    `json:"start_height"`
	MinVotingPeriod any    `json:"min_voting_period"`
	Expiration      struct {
		AtTime AtTime `json:"at_time"`
	} `json:"expiration"`
	Threshold struct {
		ThresholdQuorum struct {
			Threshold struct {
				Percent string `json:"percent"`
			} `json:"threshold"`
			Quorum struct {
				Percent string `json:"percent"`
			} `json:"quorum"`
		} `json:"threshold_quorum"`
	} `json:"threshold"`
	TotalPower string `json:"total_power"`
	Msgs       []struct {
		Custom struct {
			SubmitAdminProposal struct {
				AdminProposal struct {
					ParamChangeProposal struct {
						Title        string `json:"title"`
						Description  string `json:"description"`
						ParamChanges []struct {
							Subspace string `json:"subspace"`
							Key      string `json:"key"`
							Value    string `json:"value"`
						} `json:"param_changes"`
					} `json:"param_change_proposal"`
				} `json:"admin_proposal"`
			} `json:"submit_admin_proposal"`
		} `json:"custom"`
	} `json:"msgs"`
	Status ContractProposalStatus `json:"status"`
	Votes  struct {
		Yes     string `json:"yes"`
		No      string `json:"no"`
		Abstain string `json:"abstain"`
	} `json:"votes"`
	AllowRevoting bool `json:"allow_revoting"`
}

type ListProposalsResponse struct {
	Data struct {
		Proposals []struct {
			ID       int              `json:"id"`
			Proposal ContractProposal `json:"proposal"`
		} `json:"proposals"`
	} `json:"data"`
}

type CreditsVaultAllocationResponse struct {
	Data struct {
		AllocatedAmount string `json:"allocated_amount"`
		WithdrawnAmount string `json:"withdrawn_amount"`
		Schedule        struct {
			StartTime int64 `json:"start_time"`
			Cliff     int   `json:"cliff"`
			Duration  int   `json:"duration"`
		} `json:"schedule"`
	} `json:"data"`
}
