package event

import (
	"context"
	connect_go "github.com/bufbuild/connect-go"
	"github.com/shifty11/blocklog-backend/grpc/event/eventpb"
	"github.com/shifty11/blocklog-backend/grpc/event/eventpb/eventpbconnect"
	"google.golang.org/protobuf/types/known/emptypb"
)

//goland:noinspection GoNameStartsWithPackageName
type EventService struct {
	eventpbconnect.UnimplementedEventServiceHandler
}

func NewEventServiceHandler() eventpbconnect.EventServiceHandler {
	return &EventService{}
}

func (e EventService) EventStream(ctx context.Context, c *connect_go.Request[emptypb.Empty], c2 *connect_go.ServerStream[eventpb.Event]) error {
	//TODO implement me
	panic("implement me")
}
