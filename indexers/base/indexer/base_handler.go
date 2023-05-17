package indexer

import (
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/grpc/indexer/indexerpb"
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/indexevent"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcChannel "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/golang/protobuf/proto"
	"github.com/shifty11/go-logger/log"
	"github.com/tendermint/tendermint/abci/types"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

type baseMessageHandler struct {
	TxHandler
	chainInfo ChainInfo
	txHelper  TxHelper
}

func NewBaseMessageHandler(chainInfo ChainInfo, encodingConfig EncodingConfig) TxHandler {
	return &baseMessageHandler{
		chainInfo: chainInfo,
		txHelper:  NewTxHelper(chainInfo, encodingConfig),
	}
}

func addToResultIfNoError(result *indexerpb.HandleTxsResponse, msg sdktypes.Msg, protoMsg []byte) {
	if protoMsg != nil {
		result.ProtoMessages = append(result.ProtoMessages, protoMsg)
		result.HandledMessageTypes = append(result.HandledMessageTypes, fmt.Sprintf("%T", msg))
	}
	result.CountProcessed++
}

func (m *baseMessageHandler) handleMsg(tx []byte, anyMsg sdktypes.Msg, result *indexerpb.HandleTxsResponse) error {
	switch msg := anyMsg.(type) {
	case *banktypes.MsgSend:
		protoMsg, err := m.handleMsgSend(msg, tx)
		if err != nil {
			return err
		}
		addToResultIfNoError(result, msg, protoMsg)
	case *ibcChannel.MsgRecvPacket:
		protoMsg, err := m.handleMsgRecvPacket(msg, tx)
		if err != nil {
			return err
		}
		addToResultIfNoError(result, msg, protoMsg)
	case *stakingtypes.MsgUndelegate:
		protoMsg, err := m.handleMsgUndelegate(msg, tx)
		if err != nil {
			return err
		}
		addToResultIfNoError(result, msg, protoMsg)
	case *authz.MsgExec:
		for _, authzEncMsg := range msg.Msgs {
			authzMsg, err := sdktypes.GetMsgFromTypeURL(m.txHelper.encodingConfig.Codec, authzEncMsg.GetTypeUrl())
			if err != nil {
				return err
			}
			err = m.txHelper.encodingConfig.Codec.Unmarshal(authzEncMsg.GetValue(), authzMsg.(codec.ProtoMarshaler))
			if err != nil {
				return err
			}
			err = m.handleMsg(tx, authzMsg, result)
			if err != nil {
				return err
			}
		}
	default:
		result.UnhandledMessageTypes = append(result.UnhandledMessageTypes, fmt.Sprintf("%T", msg))
		result.CountSkipped++
	}
	return nil
}

func (m *baseMessageHandler) HandleTxs(txs [][]byte) (*indexerpb.HandleTxsResponse, error) {
	var result indexerpb.HandleTxsResponse
	for _, tx := range txs {
		txDecoded, err := m.txHelper.encodingConfig.TxConfig.TxDecoder()(tx)
		if err != nil {
			split := strings.Split(err.Error(), "/")
			if len(split) > 1 {
				result.UnhandledMessageTypes = append(result.UnhandledMessageTypes, strings.TrimSuffix(split[1], ": tx parse error"))
			} else {
				log.Sugar.Errorf("Error decoding tx: %s", err)
			}
			continue
		}
		if txDecoded == nil {
			continue
		}
		for _, anyMsg := range txDecoded.GetMsgs() {
			err = m.handleMsg(tx, anyMsg, &result)
			if err != nil {
				return nil, err
			}
		}
	}
	return &result, nil
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

func (i *baseMessageHandler) handleFungibleTokenPacketEvent(txResponse *sdktypes.TxResponse) ([]byte, error) {
	if txResponse == nil || len(txResponse.Events) == 0 {
		return nil, nil
	}
	var timestamp, err = i.txHelper.GetTxTimestamp(txResponse)
	if err != nil {
		return nil, err
	}
	txEvent := &indexevent.TxEvent{
		ChainName:  i.chainInfo.Name,
		Timestamp:  timestamppb.New(timestamp),
		NotifyTime: timestamppb.Now(),
		Event: &indexevent.TxEvent_CoinReceived{
			CoinReceived: &indexevent.CoinReceivedEvent{
				Coin: &indexevent.Coin{},
			},
		},
	}

	var receiver, sender, amount, denom, success = "receiver", "sender", "amount", "denom", "success"
	result, err := getRawEventResult(txResponse.Events, RawEvent{
		Type:       "fungible_token_packet",
		Attributes: []string{receiver, sender, amount, denom, success},
	})
	if err != nil {
		// check out this tx -> https://www.mintscan.io/osmosis/txs/8822ACEB04702476DB2D6ACA8E9AE398C7412B012DFEBDEE39BCBBCE50B872E1?height=9415274
		//log.Sugar.Warn(err)
		return nil, nil
	}
	if result[success] != "true" {
		return nil, nil
	}
	txEvent.WalletAddress = result[receiver]
	txEvent.GetCoinReceived().Sender = result[sender]
	txEvent.GetCoinReceived().GetCoin().Amount = result[amount]
	txEvent.GetCoinReceived().GetCoin().Denom = result[denom]
	return proto.Marshal(txEvent)
}

func (i *baseMessageHandler) handleMsgSend(msg *banktypes.MsgSend, tx []byte) ([]byte, error) {
	wasSuccessful, err := i.txHelper.WasTxSuccessful(tx)
	if err != nil {
		return nil, err
	}
	if wasSuccessful {
		var txEvent = &indexevent.TxEvent{
			ChainName:     i.chainInfo.Name,
			WalletAddress: msg.ToAddress,
			//Timestamp:     TODO: get timestamp from tx
			NotifyTime: timestamppb.Now(),
			Event: &indexevent.TxEvent_CoinReceived{
				CoinReceived: &indexevent.CoinReceivedEvent{
					Sender: msg.FromAddress,
					Coin: &indexevent.Coin{
						Amount: msg.Amount[0].Amount.String(),
						Denom:  msg.Amount[0].Denom,
					},
				},
			},
		}
		return proto.Marshal(txEvent)
	}
	return nil, nil
}

func (i *baseMessageHandler) handleMsgRecvPacket(_ *ibcChannel.MsgRecvPacket, tx []byte) ([]byte, error) {
	txResponse, err := i.txHelper.GetTxResponse(tx)
	if err != nil {
		return nil, err
	}
	return i.handleFungibleTokenPacketEvent(txResponse)
}

func (i *baseMessageHandler) handleMsgUndelegate(msg *stakingtypes.MsgUndelegate, tx []byte) ([]byte, error) {
	txResponse, err := i.txHelper.GetTxResponse(tx)
	if err != nil {
		return nil, err
	}
	if txResponse == nil || len(txResponse.Events) == 0 {
		return nil, nil
	}

	timestamp, err := i.txHelper.GetTxTimestamp(txResponse)
	if err != nil {
		return nil, err
	}

	var completionTimeStr, amount = "completion_time", "amount"
	result, err := getRawEventResult(txResponse.Events, RawEvent{
		Type:       "unbond",
		Attributes: []string{completionTimeStr, amount},
	})
	completionTime, err := parseTime(result[completionTimeStr])
	if err != nil {
		return nil, err
	}
	txEvent := &indexevent.TxEvent{
		ChainName:  i.chainInfo.Name,
		Timestamp:  timestamppb.New(timestamp),
		NotifyTime: timestamppb.Now(),
		Event: &indexevent.TxEvent_Unstake{
			Unstake: &indexevent.UnstakeEvent{
				Coin: &indexevent.Coin{
					Denom:  msg.Amount.Denom,
					Amount: msg.Amount.Amount.String(),
				},
				CompletionTime: completionTime,
			},
		},
	}
	return proto.Marshal(txEvent)
}
