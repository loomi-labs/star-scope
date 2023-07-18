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
	userManager          *database.UserManager
	chainManager         *database.ChainManager
	eventListenerManager *database.EventListenerManager
}

func NewSettingsServiceHandler(dbManagers *database.DbManagers) settingspbconnect.SettingsServiceHandler {
	return &SettingsService{
		userManager:          dbManagers.UserManager,
		chainManager:         dbManagers.ChainManager,
		eventListenerManager: dbManagers.EventListenerManager,
	}
}

func (s *SettingsService) GetWallets(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[settingspb.GetWalletsResponse], error) {
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

func (s *SettingsService) AddWallet(ctx context.Context, request *connect.Request[settingspb.UpdateWalletRequest]) (*connect.Response[emptypb.Empty], error) {
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

	err = s.eventListenerManager.UpdateWallet(ctx, user, chain, request.Msg)
	if err != nil {
		log.Sugar.Error("failed to update wallet: ", err)
		return nil, types.UnknownErr
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s *SettingsService) UpdateWallet(ctx context.Context, request *connect.Request[settingspb.UpdateWalletRequest]) (*connect.Response[emptypb.Empty], error) {
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

	err = s.eventListenerManager.UpdateWallet(ctx, user, chain, request.Msg)
	if err != nil {
		log.Sugar.Error("failed to update wallet: ", err)
		return nil, types.UnknownErr
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s *SettingsService) RemoveWallet(ctx context.Context, request *connect.Request[settingspb.RemoveWalletRequest]) (*connect.Response[emptypb.Empty], error) {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return nil, types.UserNotFoundErr
	}

	update := &settingspb.UpdateWalletRequest{
		WalletAddress:           request.Msg.GetWalletAddress(),
		NotifyFunding:           false,
		NotifyStaking:           false,
		NotifyGovVotingReminder: false,
	}

	err := s.eventListenerManager.UpdateWallet(ctx, user, nil, update)
	if err != nil {
		log.Sugar.Error("failed to remove wallet: ", err)
		return nil, types.UnknownErr
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s *SettingsService) ValidateWallet(ctx context.Context, request *connect.Request[settingspb.ValidateWalletRequest]) (*connect.Response[settingspb.ValidateWalletResponse], error) {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return nil, types.UserNotFoundErr
	}

	err := common.ValidateBech32Address(request.Msg.GetAddress())
	isValid := err == nil
	hasWalletAddress := false

	if isValid {
		hasWalletAddress, err = s.eventListenerManager.QueryHasWalletAddress(ctx, user, request.Msg.GetAddress())
		if err != nil {
			log.Sugar.Error("failed to query has wallet address: ", err)
			return nil, types.UnknownErr
		}
	}

	var wallet = &settingspb.Wallet{
		Address:                            request.Msg.GetAddress(),
		LogoUrl:                            "",
		NotifyFunding:                      false,
		NotifyStaking:                      false,
		NotifyGovVotingReminder:            false,
		IsNotifyFundingSupported:           false,
		IsNotifyStakingSupported:           false,
		IsNotifyGovVotingReminderSupported: false,
	}
	var isSupported = false
	for _, chain := range s.chainManager.QueryEnabled(ctx) {
		if common.IsBech32AddressFromChain(request.Msg.GetAddress(), chain.Bech32Prefix) {
			wallet.LogoUrl = chain.Image
			isSupported = true
			wallet.IsNotifyFundingSupported = chain.IsIndexing
			wallet.IsNotifyStakingSupported = chain.IsIndexing
			wallet.IsNotifyGovVotingReminderSupported = chain.IsQuerying
			break
		}
	}

	return connect.NewResponse(&settingspb.ValidateWalletResponse{
		IsValid:        isValid,
		IsSupported:    isSupported,
		IsAlreadyAdded: hasWalletAddress,
		Wallet:         wallet,
	}), nil
}
