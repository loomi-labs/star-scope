package indexer

import (
	"buf.build/gen/go/loomi-labs/star-scope/bufbuild/connect-go/grpc/indexer/indexerpb/indexerpbconnect"
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/grpc/indexer/indexerpb"
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/loomi-labs/star-scope/indexers/base/client"
	"github.com/shifty11/go-logger/log"
)

type customMessageHandler struct {
	TxHandler
	txHandlerClient    indexerpbconnect.TxHandlerServiceClient
	baseMessageHandler TxHandler
}

func NewCustomMessageHandler(chainInfo ChainInfo, encodingConfig EncodingConfig, grpcEndpoint string) TxHandler {
	return &customMessageHandler{
		txHandlerClient:    client.NewTxHandlerServiceClient(grpcEndpoint),
		baseMessageHandler: NewBaseMessageHandler(chainInfo, encodingConfig),
	}
}

func (m *customMessageHandler) HandleTxs(txs [][]byte) (*indexerpb.HandleTxsResponse, error) {
	baseResult, err := m.baseMessageHandler.HandleTxs(txs)
	if err != nil {
		return nil, err
	}
	var request = connect.NewRequest(&indexerpb.HandleTxsRequest{Txs: txs})
	response, err := m.txHandlerClient.HandleTxs(context.Background(), request)
	if err != nil {
		log.Sugar.Errorf("Error handling tx: %s", err.Error())
		return nil, err
	}
	baseResult.ProtoMessages = append(baseResult.ProtoMessages, response.Msg.ProtoMessages...)
	baseResult.CountProcessed += response.Msg.CountProcessed
	baseResult.CountSkipped += response.Msg.CountSkipped
	baseResult.HandledMessageTypes = append(baseResult.HandledMessageTypes, response.Msg.HandledMessageTypes...)
	baseResult.UnhandledMessageTypes = append(baseResult.UnhandledMessageTypes, response.Msg.UnhandledMessageTypes...)
	return baseResult, nil
}
