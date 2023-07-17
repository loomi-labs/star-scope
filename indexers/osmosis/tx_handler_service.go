package main

import (
	"buf.build/gen/go/loomi-labs/star-scope/bufbuild/connect-go/grpc/indexer/indexerpb/indexerpbconnect"
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/event"
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/grpc/indexer/indexerpb"
	"context"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/golang/protobuf/proto"
	"github.com/osmosis-labs/osmosis/osmoutils/noapptest"
	"github.com/osmosis-labs/osmosis/v15/app/keepers"
	lockuptypes "github.com/osmosis-labs/osmosis/v15/x/lockup/types"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

//goland:noinspection GoNameStartsWithPackageName
type TxHandlerService struct {
	indexerpbconnect.UnimplementedTxHandlerServiceHandler
	txHelper TxHelper
}

func NewTxHandlerServiceHandler() indexerpbconnect.TxHandlerServiceHandler {
	var encodingConfig = noapptest.MakeTestEncodingConfig(keepers.AppModuleBasics...)
	return &TxHandlerService{
		txHelper: NewTxHelper(ChainInfo{
			Path:         "osmosis",
			RestEndpoint: "https://rest.cosmos.directory/osmosis",
			Name:         "Osmosis",
		}, EncodingConfig{
			InterfaceRegistry: encodingConfig.InterfaceRegistry,
			Codec:             encodingConfig.Codec,
			TxConfig:          encodingConfig.TxConfig,
		}),
	}
}

func (t TxHandlerService) HandleTxs(_ context.Context, request *connect.Request[indexerpb.HandleTxsRequest]) (*connect.Response[indexerpb.HandleTxsResponse], error) {
	log.Sugar.Debugf("Handle %v Tx's", len(request.Msg.GetTxs()))
	var response = indexerpb.HandleTxsResponse{ProtoMessages: make([][]byte, 0)}
	for _, tx := range request.Msg.GetTxs() {
		txDecoded, err := t.txHelper.encodingConfig.TxConfig.TxDecoder()(tx)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		for _, anyMsg := range txDecoded.GetMsgs() {
			switch anyMsg.(type) {
			case *lockuptypes.MsgBeginUnlocking:
				protoMsg, err := t.handleMsgBeginUnlocking(anyMsg.(*lockuptypes.MsgBeginUnlocking), tx)
				if err != nil {
					log.Sugar.Errorf("error handling msg: %v", err)
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				response.ProtoMessages = append(response.ProtoMessages, protoMsg)
				response.HandledMessageTypes = append(response.HandledMessageTypes, fmt.Sprintf("%T", anyMsg))
				response.CountProcessed++
			default:
				response.UnhandledMessageTypes = append(response.UnhandledMessageTypes, fmt.Sprintf("%T", anyMsg))
				response.CountSkipped++
			}
		}
	}
	return connect.NewResponse(&response), nil
}

func (t TxHandlerService) handleMsgBeginUnlocking(_ *lockuptypes.MsgBeginUnlocking, tx []byte) ([]byte, error) {
	txResponse, err := t.txHelper.GetTxResponse(tx)
	if err != nil {
		return nil, err
	}
	if txResponse == nil || len(txResponse.Events) == 0 {
		return nil, nil
	}
	timestamp, err := time.Parse(time.RFC3339, txResponse.Timestamp)
	if err != nil {
		return nil, err
	}
	txEvent := &event.WalletEvent{
		ChainName:  t.txHelper.chainInfo.Name,
		Timestamp:  timestamppb.New(timestamp),
		NotifyTime: timestamppb.Now(),
		Event: &event.WalletEvent_OsmosisPoolUnlock{
			OsmosisPoolUnlock: &event.OsmosisPoolUnlockEvent{},
		},
	}
	var owner, duration, unlockTime = "owner", "duration", "unlock_time"

	result, err := getRawEventResult(txResponse.Events, RawEvent{
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
