package event

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/grpc/event/eventpb"
	"github.com/loomi-labs/star-scope/grpc/event/eventpb/eventpbconnect"
	"github.com/loomi-labs/star-scope/grpc/types"
	"github.com/loomi-labs/star-scope/kafka"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

//goland:noinspection GoNameStartsWithPackageName
type EventService struct {
	eventpbconnect.UnimplementedEventServiceHandler
	kafka                *kafka.Kafka
	chainManager         *database.ChainManager
	eventListenerManager *database.EventListenerManager
}

func NewEventServiceHandler(dbManagers *database.DbManagers, kafka *kafka.Kafka) eventpbconnect.EventServiceHandler {
	return &EventService{
		kafka:                kafka,
		chainManager:         dbManagers.ChainManager,
		eventListenerManager: dbManagers.EventListenerManager,
	}
}

func (e EventService) EventStream(ctx context.Context, _ *connect.Request[emptypb.Empty], stream *connect.ServerStream[eventpb.EventList]) error {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return types.UserNotFoundErr
	}

	processedEvents := make(chan *eventpb.EventList, 100)

	go e.kafka.ConsumeProcessedEvents(ctx, user, processedEvents)

	// Timer for sending empty message every 30 seconds
	timer := time.NewTicker(30 * time.Second)
	defer timer.Stop()

	for {
		select {
		case eventList, ok := <-processedEvents:
			if !ok {
				log.Sugar.Debugf("processed events channel closed")
				return types.UnknownErr
			}
			log.Sugar.Debugf("received processed %v events", len(eventList.GetEvents()))
			err := stream.Send(eventList)
			if err != nil {
				log.Sugar.Debugf("error sending processed eventList: %v", err)
				return types.UnknownErr
			}
		case <-timer.C:
			log.Sugar.Debugf("sending empty message")
			err := stream.Send(&eventpb.EventList{})
			if err != nil {
				log.Sugar.Debugf("error sending empty message: %v", err)
				return types.UnknownErr
			}
		}
	}
}

func (e EventService) ListEvents(ctx context.Context, request *connect.Request[eventpb.ListEventsRequest]) (*connect.Response[eventpb.EventList], error) {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return nil, types.UserNotFoundErr
	}

	els := e.eventListenerManager.QueryByUser(ctx, user)
	events := make([]*eventpb.Event, 0)
	for _, el := range els {
		for _, event := range el.Edges.Events {
			pbEvent, err := kafka.EntEventToProto(event, el.Edges.Chain)
			if err != nil {
				log.Sugar.Error(err)
				return nil, types.UnknownErr
			}
			events = append(events, pbEvent)
		}
	}
	return connect.NewResponse(&eventpb.EventList{
		Events: events,
	}), nil
}

func (e EventService) ListChains(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[eventpb.ChainList], error) {
	chains := e.chainManager.QueryEnabled(ctx)
	pbChains := make([]*eventpb.ChainData, len(chains))
	for i, chain := range chains {
		pbChains[i] = &eventpb.ChainData{
			Id:       int64(chain.ID),
			Name:     chain.Name,
			ImageUrl: chain.Image,
		}
	}

	return connect.NewResponse(&eventpb.ChainList{
		Chains: pbChains,
	}), nil
}
