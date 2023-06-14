package types

import "time"

type Authorization struct {
	Type      string `json:"@type"`
	MaxTokens any    `json:"max_tokens"`
	AllowList struct {
		Address []string `json:"address"`
	} `json:"allow_list"`
	AuthorizationType string `json:"authorization_type"`
}

type Grant struct {
	Granter       string        `json:"granter"`
	Grantee       string        `json:"grantee"`
	Authorization Authorization `json:"authorization,omitempty"`
	Expiration    time.Time     `json:"expiration"`
}

type AuthzGrantsResponse struct {
	Grants     []Grant `json:"grants"`
	Pagination struct {
		NextKey any    `json:"next_key"`
		Total   string `json:"total"`
	} `json:"pagination"`
}
