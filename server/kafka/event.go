package kafka

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/grpc/event/eventpb"
	"github.com/loomi-labs/star-scope/indexevent"
)

func TxEventToProto(data []byte) (*eventpb.Event, error) {
	var txEvent indexevent.TxEvent
	err := proto.Unmarshal(data, &txEvent)
	if err != nil {
		return nil, err
	}
	switch txEvent.GetEvent().(type) {
	case *indexevent.TxEvent_CoinReceived:
		return &eventpb.Event{
			Title:       "Token Received",
			Description: fmt.Sprintf("%v received %v%v from %v", txEvent.WalletAddress, txEvent.GetCoinReceived().GetCoin().Amount, txEvent.GetCoinReceived().GetCoin().Denom, txEvent.GetCoinReceived().Sender),
			Timestamp:   txEvent.Timestamp,
		}, nil
	case *indexevent.TxEvent_OsmosisPoolUnlock:
		return &eventpb.Event{
			Title:       "Pool Unlock",
			Description: fmt.Sprintf("%v will unlock pool at %v", txEvent.WalletAddress, txEvent.GetOsmosisPoolUnlock().UnlockTime),
			Timestamp:   txEvent.Timestamp,
		}, nil
	}
	return nil, errors.New(fmt.Sprintf("No type defined for event %v", txEvent.GetEvent()))
}

func EntEventToProto(entEvent *ent.Event) (*eventpb.Event, error) {
	var pbEvent, err = TxEventToProto(entEvent.TxEvent)
	if err != nil {
		return nil, err
	}
	pbEvent.Id = int64(entEvent.ID)
	//pbEvent.ChannelId = int64(entEvent.Edges.EventListener.Edges.Channel.ID)
	return pbEvent, nil
}
