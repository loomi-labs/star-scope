package indexer

import (
	"context"
	connect_go "github.com/bufbuild/connect-go"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/grpc/indexer/indexerpb"
	"github.com/loomi-labs/star-scope/grpc/indexer/indexerpb/indexerpbconnect"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (i IndexerService) GetHeight(ctx context.Context, request *connect_go.Request[indexerpb.GetHeightRequest]) (*connect_go.Response[indexerpb.Height], error) {
	chain, err := i.chainManager.QueryByName(ctx, request.Msg.GetChainName())
	if err != nil {
		log.Sugar.Errorf("failed to query chain by name: %v", err)
		return nil, err
	}
	return connect_go.NewResponse(&indexerpb.Height{
		Height: chain.IndexingHeight,
	}), nil
}

func (i IndexerService) UpdateHeight(ctx context.Context, request *connect_go.Request[indexerpb.UpdateHeightRequest]) (*connect_go.Response[emptypb.Empty], error) {
	if request.Msg.GetHeight() < 0 {
		return nil, status.Error(codes.InvalidArgument, "height must be positive")
	}

	_, err := i.chainManager.UpdateIndexingHeight(ctx, request.Msg.GetChainName(), request.Msg.GetHeight())
	if err != nil {
		return nil, err
	}

	return connect_go.NewResponse(&emptypb.Empty{}), nil
}
