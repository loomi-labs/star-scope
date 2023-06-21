package user

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/commchannel"
	"github.com/loomi-labs/star-scope/grpc/types"
	"github.com/loomi-labs/star-scope/grpc/user/userpb"
	"github.com/loomi-labs/star-scope/grpc/user/userpb/userpbconnect"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

//goland:noinspection GoNameStartsWithPackageName
type UserService struct {
	userpbconnect.UnimplementedUserServiceHandler
	userManager *database.UserManager
}

func NewUserServiceHandler(dbManagers *database.DbManagers) userpbconnect.UserServiceHandler {
	return &UserService{
		userManager: dbManagers.UserManager,
	}
}

func (s UserService) GetUser(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[userpb.User], error) {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return nil, types.UserNotFoundErr
	}

	username := ""
	if user.DiscordUsername != "" {
		username = user.DiscordUsername
	} else if user.TelegramUsername != "" {
		username = user.TelegramUsername
	} else {
		username = user.WalletAddress
	}

	return connect.NewResponse(&userpb.User{
		Id:          int64(user.ID),
		Name:        username,
		HasDiscord:  user.DiscordUserID != 0,
		HasTelegram: user.TelegramUserID != 0,
	}), nil
}

func (s UserService) DeleteAccount(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return nil, types.UserNotFoundErr
	}

	err := s.userManager.Delete(ctx, user)
	if err != nil {
		log.Sugar.Errorf("error deleting user: %v", err)
		return nil, types.UnknownErr
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s UserService) GetDiscordChannels(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[userpb.DiscordChannels], error) {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return nil, types.UserNotFoundErr
	}

	t := commchannel.TypeDiscord
	channels, err := s.userManager.QueryCommChannels(ctx, user, &t)
	if err != nil {
		log.Sugar.Errorf("error querying comm channels: %v", err)
		return nil, types.UnknownErr
	}

	pbChannels := make([]*userpb.DiscordChannel, len(channels))
	for i, channel := range channels {
		pbChannels[i] = &userpb.DiscordChannel{
			Id:        int64(channel.ID),
			ChannelId: channel.DiscordChannelID,
			Name:      channel.Name,
			IsGroup:   channel.IsGroup,
		}
	}

	return connect.NewResponse(&userpb.DiscordChannels{
		Channels: pbChannels,
	}), nil
}
