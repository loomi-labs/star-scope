package kafka

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/grpc/event/eventpb"
	"github.com/loomi-labs/star-scope/indexevent"
)

func txEventToProto(data []byte) (*indexevent.TxEvent, *eventpb.Event, error) {
	var txEvent indexevent.TxEvent
	err := proto.Unmarshal(data, &txEvent)
	if err != nil {
		return nil, nil, err
	}
	switch txEvent.GetEvent().(type) {
	case *indexevent.TxEvent_CoinReceived:
		return &txEvent, &eventpb.Event{
			Title:       "Token Received",
			Description: fmt.Sprintf("%v received %v%v from %v", txEvent.WalletAddress, txEvent.GetCoinReceived().GetCoin().Amount, txEvent.GetCoinReceived().GetCoin().Denom, txEvent.GetCoinReceived().Sender),
			Timestamp:   txEvent.Timestamp,
			EventType:   eventpb.EventType_FUNDING,
		}, nil
	case *indexevent.TxEvent_OsmosisPoolUnlock:
		return &txEvent, &eventpb.Event{
			Title:       "Pool Unlock",
			Description: fmt.Sprintf("%v will unlock pool at %v", txEvent.WalletAddress, txEvent.GetOsmosisPoolUnlock().UnlockTime),
			Timestamp:   txEvent.Timestamp,
			EventType:   eventpb.EventType_DEX,
		}, nil
	case *indexevent.TxEvent_Unstake:
		return &txEvent, &eventpb.Event{
			Title:       "Unstake",
			Description: fmt.Sprintf("%v will unstake %v%v at %v", txEvent.WalletAddress, txEvent.GetUnstake().GetCoin().Amount, txEvent.GetUnstake().GetCoin().Denom, txEvent.GetUnstake().CompletionTime),
			Timestamp:   txEvent.Timestamp,
			EventType:   eventpb.EventType_STAKING,
		}, nil
	}
	return nil, nil, errors.New(fmt.Sprintf("No type defined for event %v", txEvent.GetEvent()))
}

func kafkaMsgToProto(data []byte, chains []*ent.Chain) (*eventpb.Event, error) {
	var txEvent, pbEvent, err = txEventToProto(data)
	if err != nil {
		return nil, err
	}
	for _, chain := range chains {
		if chain.Path == txEvent.ChainPath {
			pbEvent.Chain = &eventpb.ChainData{
				Id:       int64(chain.ID),
				Name:     chain.Name,
				ImageUrl: chain.Image,
			}
			break
		}
	}
	return pbEvent, nil
}

func EntEventToProto(entEvent *ent.Event, chain *ent.Chain) (*eventpb.Event, error) {
	var _, pbEvent, err = txEventToProto(entEvent.TxEvent)
	if err != nil {
		return nil, err
	}
	pbEvent.Id = int64(entEvent.ID)
	pbEvent.Chain = &eventpb.ChainData{
		Id:       int64(chain.ID),
		Name:     chain.Name,
		ImageUrl: chain.Image,
	}
	return pbEvent, nil
}
