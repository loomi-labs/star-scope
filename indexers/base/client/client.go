package client

import (
	"buf.build/gen/go/loomi-labs/star-scope/bufbuild/connect-go/grpc/indexer/indexerpb/indexerpbconnect"
	"github.com/bufbuild/connect-go"
	"net/http"
)

func NewIndexerServiceClient(endpoint string, authToken string) indexerpbconnect.IndexerServiceClient {
	authInterceptor := NewAuthInterceptor(authToken)
	interceptors := connect.WithInterceptors(authInterceptor)
	return indexerpbconnect.NewIndexerServiceClient(
		http.DefaultClient,
		endpoint,
		interceptors,
	)
}

func NewTxHandlerServiceClient(endpoint string) indexerpbconnect.TxHandlerServiceClient {
	return indexerpbconnect.NewTxHandlerServiceClient(
		http.DefaultClient,
		endpoint,
	)
}
