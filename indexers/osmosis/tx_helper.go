package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/tendermint/tendermint/abci/types"
	"golang.org/x/exp/slices"
	"io"
	"net/http"
	"strings"
)

type ChainInfo struct {
	Path         string
	RestEndpoint string
	Name         string
}

type EncodingConfig struct {
	InterfaceRegistry codectypes.InterfaceRegistry
	Codec             codec.Codec
	TxConfig          client.TxConfig
}

type TxHelper struct {
	chainInfo      ChainInfo
	encodingConfig EncodingConfig
}

func NewTxHelper(chainInfo ChainInfo, encodingConfig EncodingConfig) TxHelper {
	return TxHelper{
		chainInfo:      chainInfo,
		encodingConfig: encodingConfig,
	}
}

func (h *TxHelper) GetTxResult(tx []byte) (*txtypes.GetTxResponse, error) {
	hash := sha256.Sum256(tx)
	hashString := hex.EncodeToString(hash[:])

	var url = fmt.Sprintf("%v/cosmos/tx/v1beta1/txs/%v", h.chainInfo.RestEndpoint, hashString)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Status code: %v", resp.StatusCode))
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var txResponse txtypes.GetTxResponse
	if err := h.encodingConfig.Codec.UnmarshalJSON(body, &txResponse); err != nil {
		return nil, err
	}
	return &txResponse, nil
}

func (h *TxHelper) GetTxResponse(tx []byte) (*sdktypes.TxResponse, error) {
	resp, err := h.GetTxResult(tx)
	if err != nil {
		return nil, err
	}
	if resp.GetTxResponse().Code == 0 {
		return resp.GetTxResponse(), nil
	}
	return nil, nil
}

func (h *TxHelper) WasTxSuccessful(tx []byte) (bool, error) {
	txResponse, err := h.GetTxResponse(tx)
	if err != nil {
		return false, err
	}
	if txResponse == nil {
		return false, nil
	}
	return len(txResponse.Events) > 0, nil
}

type RawEvent struct {
	Type       string
	Attributes []string
}

func getRawEventResult(events []types.Event, event RawEvent) (map[string]string, error) {
	var result = make(map[string]string)
	for _, e := range events {
		if e.Type == event.Type {
			for _, attribute := range e.Attributes {
				if slices.Contains(event.Attributes, string(attribute.GetKey())) {
					result[string(attribute.GetKey())] = string(attribute.GetValue())
				}
			}
		}

	}
	if len(result) != len(event.Attributes) {
		var missing []string
		for _, attr := range event.Attributes {
			if _, ok := result[attr]; !ok {
				missing = append(missing, attr)
			}
		}
		return nil, errors.New(fmt.Sprintf("missing attributes: %v", strings.Join(missing, ", ")))
	}
	return result, nil
}
