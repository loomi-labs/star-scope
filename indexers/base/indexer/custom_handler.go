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
	txHandlerClient indexerpbconnect.TxHandlerServiceClient
}

func NewCustomMessageHandler(grpcEndpoint string) TxHandler {
	return &customMessageHandler{
		txHandlerClient: client.NewTxHandlerServiceClient(grpcEndpoint),
	}
}

func (m *customMessageHandler) HandleTxs(txs [][]byte) (*indexerpb.HandleTxsResponse, error) {
	var request = connect.NewRequest(&indexerpb.HandleTxsRequest{Txs: txs})
	response, err := m.txHandlerClient.HandleTxs(context.Background(), request)
	if err != nil {
		log.Sugar.Errorf("Error handling tx: %s", err.Error())
		return nil, err
	}
	return response.Msg, nil
}
