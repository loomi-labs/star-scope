package indexer

import (
	"context"
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
	chains := i.chainManager.QueryIsIndexing(ctx)
	pbChains := make([]*indexerpb.IndexingChain, len(chains))
	for ix, chain := range chains {
		pbChains[ix] = &indexerpb.IndexingChain{
			Id:                    uint64(chain.ID),
			Name:                  chain.Name,
			Path:                  chain.Path,
			RestEndpoint:          chain.RestEndpoint,
			IndexingHeight:        chain.IndexingHeight,
			UnhandledMessageTypes: strings.Split(chain.UnhandledMessageTypes, ","),
			HasCustomIndexer:      chain.HasCustomIndexer,
		}
	}
	return connect.NewResponse(&indexerpb.GetIndexingChainsResponse{Chains: pbChains}), nil
}

func (i IndexerService) UpdateIndexingChains(ctx context.Context, request *connect.Request[indexerpb.UpdateIndexingChainsRequest]) (*connect.Response[indexerpb.UpdateIndexingChainsResponse], error) {
	disabledChainIds := make([]uint64, 0)
	for _, chain := range request.Msg.GetChains() {
		updatedChain, err := i.chainManager.UpdateIndexStatus(ctx, int(chain.GetId()), chain.GetIndexingHeight(), chain.GetHandledMessageTypes(), chain.GetUnhandledMessageTypes())
		if err != nil {
			log.Sugar.Errorf("error while updating chain: %v", err)
		}
		if !updatedChain.IsEnabled {
			disabledChainIds = append(disabledChainIds, uint64(updatedChain.ID))
		}
	}

	return connect.NewResponse(&indexerpb.UpdateIndexingChainsResponse{DisabledChainIds: disabledChainIds}), nil
}
