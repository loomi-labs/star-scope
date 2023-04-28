package subscription

import (
	"context"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/shifty11/blocklog-backend/grpc/subscription/v1"
	subscriptionconnect "github.com/shifty11/blocklog-backend/grpc/subscription/v1/v1connect"
	"google.golang.org/protobuf/types/known/emptypb"
)

//goland:noinspection GoNameStartsWithPackageName
type SubscriptionService struct {
	subscriptionconnect.UnimplementedSubscriptionServiceHandler
}

func NewSubscriptionServiceHandler() subscriptionconnect.SubscriptionServiceHandler {
	return &SubscriptionService{}
}

func (s SubscriptionService) ListSubscriptions(ctx context.Context, c *connect_go.Request[emptypb.Empty]) (*connect_go.Response[v1.ListSubscriptionsResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (s SubscriptionService) Subscribe(ctx context.Context, c *connect_go.Request[v1.SubscribeRequest]) (*connect_go.Response[v1.SubscribeResponse], error) {
	//TODO implement me
	panic("implement me")
}
