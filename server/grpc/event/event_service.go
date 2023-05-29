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
	pbEvents := make([]*eventpb.Event, 0)
	for _, el := range els {
		events, err := e.eventListenerManager.QueryEvents(
			ctx,
			el,
			request.Msg.EventType,
			request.Msg.GetStartTime(),
			request.Msg.GetEndTime(),
			request.Msg.GetLimit(),
			request.Msg.GetOffset(),
		)
		if err != nil {
			log.Sugar.Error(err)
			return nil, types.UnknownErr
		}
		for _, event := range events {
			pbEvent, err := kafka.EntEventToProto(event, el.Edges.Chain)
			if err != nil {
				log.Sugar.Error(err)
				return nil, types.UnknownErr
			}
			pbEvents = append(pbEvents, pbEvent)
		}
	}
	return connect.NewResponse(&eventpb.EventList{
		Events: pbEvents,
	}), nil
}

func (e EventService) ListChains(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[eventpb.ChainList], error) {
	chains := e.chainManager.QueryEnabled(ctx)
	pbChains := make([]*eventpb.ChainData, len(chains))
	for i, chain := range chains {
		pbChains[i] = &eventpb.ChainData{
			Id:       int64(chain.ID),
			Name:     chain.PrettyName,
			ImageUrl: chain.Image,
		}
	}

	return connect.NewResponse(&eventpb.ChainList{
		Chains: pbChains,
	}), nil
}

func (e EventService) ListEventsCount(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[eventpb.ListEventsCountResponse], error) {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return nil, types.UserNotFoundErr
	}

	cntRead, err := e.eventListenerManager.QueryCountEventsByType(ctx, user, true)
	if err != nil {
		log.Sugar.Error(err)
		return nil, types.UnknownErr
	}
	cntUnread, err := e.eventListenerManager.QueryCountEventsByType(ctx, user, false)
	if err != nil {
		log.Sugar.Error(err)
		return nil, types.UnknownErr
	}
	counters := make([]*eventpb.EventsCount, len(eventpb.EventType_name))
	for i, name := range eventpb.EventType_name {
		counters[i] = &eventpb.EventsCount{
			EventType: eventpb.EventType(i),
			Count:     0,
		}
		for _, cnt := range cntRead {
			if cnt.EventType.String() == name {
				counters[i].Count = int32(cnt.Count)
				break
			}
		}
		for _, cnt := range cntUnread {
			if cnt.EventType.String() == name {
				counters[i].UnreadCount += int32(cnt.Count)
				break
			}
		}
	}

	return connect.NewResponse(&eventpb.ListEventsCountResponse{
		Counters: counters,
	}), nil
}

func (e EventService) MarkEventRead(ctx context.Context, request *connect.Request[eventpb.MarkEventReadRequest]) (*connect.Response[emptypb.Empty], error) {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return nil, types.UserNotFoundErr
	}

	err := e.eventListenerManager.UpdateMarkEventRead(ctx, user, int(request.Msg.GetEventId()))
	if err != nil {
		log.Sugar.Error(err)
		return nil, types.UnknownErr
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}
