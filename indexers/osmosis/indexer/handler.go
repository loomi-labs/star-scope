package indexer

import (
	"errors"
	"fmt"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibcChannel "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	"github.com/golang/protobuf/proto"
	lockuptypes "github.com/osmosis-labs/osmosis/v15/x/lockup/types"
	indexEvent "github.com/shifty11/blocklog-backend/indexers/osmosis/index_event"
	"github.com/shifty11/go-logger/log"
	"github.com/tendermint/tendermint/abci/types"
	"golang.org/x/exp/slices"
	"strings"
)

type RawEvent struct {
	Type       string
	Attributes []string
}

func (i *Indexer) getRawEventResult(events []types.Event, event RawEvent) (map[string]string, error) {
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

func (i *Indexer) handleFungibleTokenPacketEvent(events []types.Event) ([]byte, error) {
	if len(events) == 0 {
		return nil, nil
	}
	txEvent := &indexEvent.TxEvent{
		ChainName: i.chainInfo.ChainName,
		Event: &indexEvent.TxEvent_CoinReceived{
			CoinReceived: &indexEvent.CoinReceivedEvent{
				Coin: &indexEvent.Coin{},
			},
		},
	}

	var receiver, sender, amount, denom, success = "receiver", "sender", "amount", "denom", "success"
	result, err := i.getRawEventResult(events, RawEvent{
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

func (i *Indexer) handleMsgSend(msg *banktypes.MsgSend, tx []byte) ([]byte, error) {
	wasSuccessful, err := i.wasTxSuccessful(tx)
	if err != nil {
		return nil, err
	}
	if wasSuccessful {
		var txEvent = &indexEvent.TxEvent{
			ChainName:     i.chainInfo.ChainName,
			WalletAddress: msg.ToAddress,
			Event: &indexEvent.TxEvent_CoinReceived{
				CoinReceived: &indexEvent.CoinReceivedEvent{
					Sender: msg.FromAddress,
					Coin: &indexEvent.Coin{
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

func (i *Indexer) handleMsgMultiSend(_ *banktypes.MsgMultiSend, _ []byte, height int64) {
	log.Sugar.Errorf("MsgMultiSend not implemented: height: %v on %v", height, i.chainInfo.ChainName)
}

func (i *Indexer) handleMsgRecvPacket(_ *ibcChannel.MsgRecvPacket, tx []byte) ([]byte, error) {
	events, err := i.getTxEvents(tx)
	if err != nil {
		return nil, err
	}
	return i.handleFungibleTokenPacketEvent(events)
}

func (i *Indexer) handleMsgBeginUnlockingAll(_ *lockuptypes.MsgBeginUnlockingAll, _ []byte, height int64) {
	log.Sugar.Errorf("MsgBeginUnlockingAll not implemented: height: %v on %v", height, i.chainInfo.ChainName)
}

func (i *Indexer) handleMsgBeginUnlocking(_ *lockuptypes.MsgBeginUnlocking, tx []byte) ([]byte, error) {
	txEvent := &indexEvent.TxEvent{
		ChainName: i.chainInfo.ChainName,
		Event: &indexEvent.TxEvent_OsmosisPoolUnlock{
			OsmosisPoolUnlock: &indexEvent.OsmosisPoolUnlockEvent{},
		},
	}
	var owner, duration, unlockTime = "owner", "duration", "unlock_time"
	events, err := i.getTxEvents(tx)
	if err != nil {
		return nil, err
	}
	result, err := i.getRawEventResult(events, RawEvent{
		Type:       "begin_unlock",
		Attributes: []string{owner, duration, unlockTime},
	})
	if err != nil {
		log.Sugar.Error(err)
		return nil, nil
	}
	txEvent.WalletAddress = result[owner]
	dur, err := parseDuration(result[duration])
	if err != nil {
		log.Sugar.Errorf("Failed to parse duration: %v", err)
		return nil, nil
	}
	txEvent.GetOsmosisPoolUnlock().Duration = dur
	unlTime, err := parseTime(result[unlockTime])
	if err != nil {
		log.Sugar.Errorf("Failed to parse unlock time: %v", err)
		return nil, nil
	}
	txEvent.GetOsmosisPoolUnlock().UnlockTime = unlTime
	return proto.Marshal(txEvent)
}
