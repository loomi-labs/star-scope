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
