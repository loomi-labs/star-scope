package types

import (
	"time"
)

type UnstakingEntry struct {
	CreationHeight string    `json:"creation_height"`
	CompletionTime time.Time `json:"completion_time"`
	InitialBalance string    `json:"initial_balance"`
	Balance        string    `json:"balance"`
}

type UnstakingResponse struct {
	UnbondingResponses []struct {
		DelegatorAddress string           `json:"delegator_address"`
		ValidatorAddress string           `json:"validator_address"`
		Entries          []UnstakingEntry `json:"entries"`
	} `json:"unbonding_responses"`
	Pagination struct {
		NextKey any    `json:"next_key"`
		Total   string `json:"total"`
	} `json:"pagination"`
}
