package client

import (
	"buf.build/gen/go/loomi-labs/star-scope/bufbuild/connect-go/grpc/indexer/indexerpb/indexerpbconnect"
	"github.com/bufbuild/connect-go"
	"github.com/loomi-labs/star-scope/indexers/osmosis/common"
	"net/http"
)

func GetClient() indexerpbconnect.IndexerServiceClient {
	authInterceptor := NewAuthInterceptor(common.GetEnvX("INDEXER_AUTH_TOKEN"))
	interceptors := connect.WithInterceptors(authInterceptor)
	return indexerpbconnect.NewIndexerServiceClient(
		http.DefaultClient,
		"http://localhost:50001",
		interceptors,
	)
}
