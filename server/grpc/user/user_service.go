package user

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/grpc/types"
	"github.com/loomi-labs/star-scope/grpc/user/userpb"
	"github.com/loomi-labs/star-scope/grpc/user/userpb/userpbconnect"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

//goland:noinspection GoNameStartsWithPackageName
type UserService struct {
	userpbconnect.UnimplementedUserServiceHandler
}

func NewUserServiceHandler() userpbconnect.UserServiceHandler {
	return &UserService{}
}

func (e UserService) GetUser(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[userpb.User], error) {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return nil, types.UserNotFoundErr
	}

	return connect.NewResponse(&userpb.User{
		Id:   int64(user.ID),
		Name: user.Name,
	}), nil
}

func (e UserService) ListChannels(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[userpb.ListChannelsResponse], error) {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return nil, types.UserNotFoundErr
	}

	channels, err := user.
		QueryProjects().
		QueryChannels().
		//WithEventListeners(
		//	func(query *ent.EventListenerQuery) {
		//		query.WithEvents()
		//	}).
		All(ctx)
	if err != nil {
		log.Sugar.Errorf("failed to query channels: %v", err)
		return nil, err
	}

	var channelPbs []*userpb.Channel
	for _, channel := range channels {
		channelPbs = append(channelPbs, &userpb.Channel{
			Id:   int64(channel.ID),
			Name: channel.Name,
		})
	}

	return connect.NewResponse(&userpb.ListChannelsResponse{
		Channels: channelPbs,
	}), nil
}
