package common

import "github.com/cosmos/cosmos-sdk/types/bech32"

func ValidateBech32Address(address string) error {
	_, _, err := bech32.DecodeAndConvert(address)
	return err
}

func ConvertWithOtherPrefix(address string, newPrefix string) (string, error) {
	_, bytes, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return "", err
	}
	return bech32.ConvertAndEncode(newPrefix, bytes)
}
