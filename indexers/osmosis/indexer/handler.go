package indexer

import (
	"errors"
	"fmt"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibcChannel "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	lockuptypes "github.com/osmosis-labs/osmosis/v15/x/lockup/types"
	indexEvent "github.com/shifty11/blocklog-backend/indexers/osmosis/index_event"
	"github.com/shifty11/go-logger/log"
	"github.com/tendermint/tendermint/abci/types"
	"golang.org/x/exp/slices"
	"math/big"
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

func (i *Indexer) handleFungibleTokenPacketEvent(events []types.Event) {
	if len(events) == 0 {
		return
	}
	txEvent := &indexEvent.TxEvent{
		ChainName: i.chainInfo.ChainName,
		Event: &indexEvent.TxEvent_CoinReceived{
			CoinReceived: &indexEvent.CoinReceivedEvent{},
		},
	}

	var receiver, sender, amount, denom, success = "receiver", "sender", "amount", "denom", "success"
	result, err := i.getRawEventResult(events, RawEvent{
		Type:       "fungible_token_packet",
		Attributes: []string{receiver, sender, amount, denom, success},
	})
	if err != nil {
		// check out this tx -> https://www.mintscan.io/osmosis/txs/8822ACEB04702476DB2D6ACA8E9AE398C7412B012DFEBDEE39BCBBCE50B872E1?height=9415274
		log.Sugar.Error(err)
		return
	}
	if result[success] != "true" {
		return
	}
	txEvent.WalletAddress = result[receiver]
	txEvent.GetCoinReceived().Sender = result[sender]
	num := new(big.Int)
	num, ok := num.SetString(result[amount], 10) // convert to big.Int to make sure it is a valid number
	if !ok {
		log.Sugar.Errorf("Failed to parse amount: %v", result[amount])
		return
	}
	log.Sugar.Infof("Amount: %v", num)
	txEvent.GetCoinReceived().Amount = num.String()
	txEvent.GetCoinReceived().Denom = result[denom]
	i.kafkaProducer.Produce(txEvent)
}

func (i *Indexer) handleMsgSend(msg *banktypes.MsgSend, tx []byte) {
	if i.wasTxSuccessful(tx) {
		i.kafkaProducer.Produce(&indexEvent.TxEvent{
			ChainName:     i.chainInfo.ChainName,
			WalletAddress: msg.ToAddress,
			Event: &indexEvent.TxEvent_CoinReceived{
				CoinReceived: &indexEvent.CoinReceivedEvent{
					Denom:  msg.Amount[0].Denom,
					Amount: msg.Amount[0].Amount.String(),
					Sender: msg.FromAddress,
				},
			},
		})
	}
}

func (i *Indexer) handleMsgMultiSend(_ *banktypes.MsgMultiSend, _ []byte, height int64) {
	log.Sugar.Errorf("MsgMultiSend not implemented: height: %v on %v", height, i.chainInfo.ChainName)
}

func (i *Indexer) handleMsgRecvPacket(_ *ibcChannel.MsgRecvPacket, tx []byte) {
	i.handleFungibleTokenPacketEvent(i.getTxEvents(tx))
}

func (i *Indexer) handleMsgBeginUnlockingAll(_ *lockuptypes.MsgBeginUnlockingAll, _ []byte, height int64) {
	log.Sugar.Errorf("MsgBeginUnlockingAll not implemented: height: %v on %v", height, i.chainInfo.ChainName)
}

func (i *Indexer) handleMsgBeginUnlocking(_ *lockuptypes.MsgBeginUnlocking, tx []byte) {
	txEvent := &indexEvent.TxEvent{
		ChainName: i.chainInfo.ChainName,
		Event: &indexEvent.TxEvent_OsmosisPoolUnlock{
			OsmosisPoolUnlock: &indexEvent.OsmosisPoolUnlockEvent{},
		},
	}
	var owner, duration, unlockTime = "owner", "duration", "unlock_time"
	result, err := i.getRawEventResult(i.getTxEvents(tx), RawEvent{
		Type:       "begin_unlock",
		Attributes: []string{owner, duration, unlockTime},
	})
	if err != nil {
		log.Sugar.Error(err)
		return
	}
	txEvent.WalletAddress = result[owner]
	dur, err := parseDuration(result[duration])
	if err != nil {
		log.Sugar.Errorf("Failed to parse duration: %v", err)
		return
	}
	txEvent.GetOsmosisPoolUnlock().Duration = dur
	unlTime, err := parseTime(result[unlockTime])
	if err != nil {
		log.Sugar.Errorf("Failed to parse unlock time: %v", err)
		return
	}
	txEvent.GetOsmosisPoolUnlock().UnlockTime = unlTime
	i.kafkaProducer.Produce(txEvent)
}
