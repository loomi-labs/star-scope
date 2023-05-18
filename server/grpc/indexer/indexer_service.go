package indexer

import (
	"context"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/grpc/indexer/indexerpb"
	"github.com/loomi-labs/star-scope/grpc/indexer/indexerpb/indexerpbconnect"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

//goland:noinspection GoNameStartsWithPackageName
type IndexerService struct {
	indexerpbconnect.UnimplementedIndexerServiceHandler
	chainManager *database.ChainManager
}

func NewIndexerServiceHandler(dbManagers *database.DbManagers) indexerpbconnect.IndexerServiceHandler {
	return &IndexerService{
		chainManager: dbManagers.ChainManager,
	}
}

func (i IndexerService) GetIndexingChains(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[indexerpb.GetIndexingChainsResponse], error) {
	chains := i.chainManager.QueryEnabled(ctx)
	pbChains := make([]*indexerpb.Chain, len(chains))
	for ix, chain := range chains {
		pbChains[ix] = &indexerpb.Chain{
			Id:                    uint64(chain.ID),
			Name:                  chain.Name,
			Path:                  chain.Path,
			RpcUrl:                fmt.Sprintf("https://rest.cosmos.directory/%s", chain.Path),
			IndexingHeight:        chain.IndexingHeight,
			UnhandledMessageTypes: strings.Split(chain.UnhandledMessageTypes, ","),
			HasCustomIndexer:      chain.HasCustomIndexer,
		}
	}

	return connect.NewResponse(&indexerpb.GetIndexingChainsResponse{Chains: pbChains}), nil
}

func (i IndexerService) UpdateIndexingChains(ctx context.Context, request *connect.Request[indexerpb.UpdateIndexingChainsRequest]) (*connect.Response[emptypb.Empty], error) {
	for _, chain := range request.Msg.GetChains() {
		err := i.chainManager.UpdateIndexStatus(ctx, int(chain.GetId()), chain.GetIndexingHeight(), chain.GetHandledMessageTypes(), chain.GetUnhandledMessageTypes())
		if err != nil {
			log.Sugar.Errorf("error while updating chain: %v", err)
		}
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}
