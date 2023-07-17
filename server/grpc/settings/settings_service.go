package settings

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/grpc/settings/settingspb"
	"github.com/loomi-labs/star-scope/grpc/settings/settingspb/settingspbconnect"
	"github.com/loomi-labs/star-scope/grpc/types"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

//goland:noinspection GoNameStartsWithPackageName
type SettingsService struct {
	settingspbconnect.UnimplementedSettingsServiceHandler
	userManager  *database.UserManager
	chainManager *database.ChainManager
}

func NewSettingsServiceHandler(dbManagers *database.DbManagers) settingspbconnect.SettingsServiceHandler {
	return &SettingsService{
		userManager:  dbManagers.UserManager,
		chainManager: dbManagers.ChainManager,
	}
}

func (s SettingsService) GetWallets(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[settingspb.GetWalletsResponse], error) {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return nil, types.UserNotFoundErr
	}

	wallets, err := s.userManager.QueryWallets(ctx, user)
	if err != nil {
		log.Sugar.Error("failed to query wallets: ", err)
		return nil, err
	}

	return connect.NewResponse(
		&settingspb.GetWalletsResponse{Wallets: wallets},
	), nil
}

func (s SettingsService) AddWallet(ctx context.Context, request *connect.Request[settingspb.UpdateWalletRequest]) (*connect.Response[emptypb.Empty], error) {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return nil, types.UserNotFoundErr
	}

	chain, err := s.chainManager.QueryByNewAddress(ctx, request.Msg.GetWalletAddress())
	if err != nil {
		log.Sugar.Error("failed to query chain by address: ", err)
		return nil, types.UnknownErr
	}

	err = s.userManager.UpdateWallet(ctx, user, chain, request.Msg)
	if err != nil {
		log.Sugar.Error("failed to update wallet: ", err)
		return nil, types.UnknownErr
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s SettingsService) UpdateWallet(ctx context.Context, request *connect.Request[settingspb.UpdateWalletRequest]) (*connect.Response[emptypb.Empty], error) {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return nil, types.UserNotFoundErr
	}

	chain, err := s.chainManager.QueryByNewAddress(ctx, request.Msg.GetWalletAddress())
	if err != nil {
		log.Sugar.Error("failed to query chain by address: ", err)
		return nil, types.UnknownErr
	}

	err = s.userManager.UpdateWallet(ctx, user, chain, request.Msg)
	if err != nil {
		log.Sugar.Error("failed to update wallet: ", err)
		return nil, types.UnknownErr
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s SettingsService) RemoveWallet(ctx context.Context, request *connect.Request[settingspb.RemoveWalletRequest]) (*connect.Response[emptypb.Empty], error) {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return nil, types.UserNotFoundErr
	}

	update := &settingspb.UpdateWalletRequest{
		WalletAddress:           request.Msg.GetWalletAddress(),
		NotifyStaking:           false,
		NotifyGovVotingReminder: false,
	}

	err := s.userManager.UpdateWallet(ctx, user, nil, update)
	if err != nil {
		log.Sugar.Error("failed to remove wallet: ", err)
		return nil, types.UnknownErr
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}
