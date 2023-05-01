package event

import (
	"context"
	"fmt"
	connect_go "github.com/bufbuild/connect-go"
	"github.com/shifty11/blocklog-backend/common"
	"github.com/shifty11/blocklog-backend/ent"
	"github.com/shifty11/blocklog-backend/grpc/event/eventpb"
	"github.com/shifty11/blocklog-backend/grpc/event/eventpb/eventpbconnect"
	"github.com/shifty11/blocklog-backend/grpc/types"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

//goland:noinspection GoNameStartsWithPackageName
type EventService struct {
	eventpbconnect.UnimplementedEventServiceHandler
}

func NewEventServiceHandler() eventpbconnect.EventServiceHandler {
	return &EventService{}
}

func (e EventService) EventStream(ctx context.Context, c *connect_go.Request[emptypb.Empty], c2 *connect_go.ServerStream[eventpb.Event]) error {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return types.UserNotFoundErr
	}

	for {
		err := c2.Send(&eventpb.Event{
			Id:          int64(0),
			Title:       "Hello",
			Description: fmt.Sprintf("Hello from the server to %v!", user.Name),
		})
		if err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
}
