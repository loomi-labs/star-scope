package types

import "encoding/json"

type DiscordIdentity struct {
	Id               int64  `json:"id"`
	Username         string `json:"username"`
	GlobalName       string `json:"global_name"`
	Avatar           string `json:"avatar"`
	Discriminator    string `json:"discriminator"`
	PublicFlags      int    `json:"public_flags"`
	Flags            int    `json:"flags"`
	Banner           any    `json:"banner"`
	BannerColor      any    `json:"banner_color"`
	AccentColor      any    `json:"accent_color"`
	Locale           string `json:"locale"`
	MfaEnabled       bool   `json:"mfa_enabled"`
	PremiumType      int    `json:"premium_type"`
	AvatarDecoration any    `json:"avatar_decoration"`
}

func (d *DiscordIdentity) UnmarshalJSON(data []byte) error {
	// Define a struct with fields of the original type
	type discordIdentityAlias DiscordIdentity
	aux := &struct {
		Id json.Number `json:"id"`
		*discordIdentityAlias
	}{
		discordIdentityAlias: (*discordIdentityAlias)(d),
	}

	// Unmarshal the JSON into the auxiliary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Convert the Id from json.Number to int64
	id, err := aux.Id.Int64()
	if err != nil {
		return err
	}

	// Update the Id field of the main struct
	d.Id = id
	return nil
}
