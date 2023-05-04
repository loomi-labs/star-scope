package event

import (
	"context"
	connect_go "github.com/bufbuild/connect-go"
	"github.com/shifty11/blocklog-backend/common"
	"github.com/shifty11/blocklog-backend/ent"
	"github.com/shifty11/blocklog-backend/grpc/event/eventpb"
	"github.com/shifty11/blocklog-backend/grpc/event/eventpb/eventpbconnect"
	"github.com/shifty11/blocklog-backend/grpc/types"
	"github.com/shifty11/blocklog-backend/kafka"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

//goland:noinspection GoNameStartsWithPackageName
type EventService struct {
	eventpbconnect.UnimplementedEventServiceHandler
	kafka *kafka.Kafka
}

func NewEventServiceHandler(kafka *kafka.Kafka) eventpbconnect.EventServiceHandler {
	return &EventService{
		kafka: kafka,
	}
}

func (e EventService) EventStream(ctx context.Context, _ *connect_go.Request[emptypb.Empty], stream *connect_go.ServerStream[eventpb.Event]) error {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return types.UserNotFoundErr
	}

	// TODO: cancel the stream when the client disconnects
	processedEvents := make(chan *eventpb.Event, 100)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go e.kafka.ConsumeProcessedEvents(ctx, user, processedEvents)

	for {
		event, ok := <-processedEvents
		if !ok {
			log.Sugar.Debugf("processed events channel closed")
			return types.UnknownErr
		}
		log.Sugar.Debugf("received processed event: %v", event)
		err := stream.Send(event)
		if err != nil {
			log.Sugar.Error(err)
			return types.UnknownErr
		}
	}
}
