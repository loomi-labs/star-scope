package event

import (
	"context"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/shifty11/blocklog-backend/grpc/event/v1"
	eventconnect "github.com/shifty11/blocklog-backend/grpc/event/v1/v1connect"
	"google.golang.org/protobuf/types/known/emptypb"
)

//goland:noinspection GoNameStartsWithPackageName
type EventService struct {
	eventconnect.UnimplementedEventServiceHandler
}

func NewEventServiceHandler() eventconnect.EventServiceHandler {
	return &EventService{}
}

func (e EventService) EventStream(ctx context.Context, c *connect_go.Request[emptypb.Empty], c2 *connect_go.ServerStream[v1.Event]) error {
	//TODO implement me
	panic("implement me")
}
