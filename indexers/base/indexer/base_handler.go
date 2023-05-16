package indexer

import (
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/indexevent"
	"errors"
	"fmt"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibcChannel "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/golang/protobuf/proto"
	"github.com/shifty11/go-logger/log"
	"github.com/tendermint/tendermint/abci/types"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"time"
)

type baseMessageHandler struct {
	MessageHandler
	chainInfo ChainInfo
	txHelper  TxHelper
}

func NewBaseMessageHandler(chainInfo ChainInfo, encodingConfig EncodingConfig) MessageHandler {
	return &baseMessageHandler{
		chainInfo: chainInfo,
		txHelper:  NewTxHelper(chainInfo, encodingConfig),
	}
}
func (m *baseMessageHandler) DecodeTx(tx []byte) (sdktypes.Tx, error) {
	txDecoded, _ := m.txHelper.encodingConfig.TxConfig.TxDecoder()(tx)
	return txDecoded, nil
}

func (m *baseMessageHandler) HandleMessage(anyMsg sdktypes.Msg, tx []byte) ([]byte, error) {
	switch msg := anyMsg.(type) {
	case *banktypes.MsgSend:
		return m.handleMsgSend(msg, tx)
	case *ibcChannel.MsgRecvPacket:
		return m.handleMsgRecvPacket(msg, tx)
	}
	return nil, nil
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
	var timestamp, err = time.Parse(time.RFC3339, txResponse.Timestamp)
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
		log.Sugar.Warn(err)
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
			NotifyTime:    timestamppb.Now(),
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
