package database

import (
	"context"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/event"
	"github.com/loomi-labs/star-scope/ent/eventlistener"
	"github.com/loomi-labs/star-scope/ent/user"
	"github.com/loomi-labs/star-scope/grpc/event/eventpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type EventListenerManager struct {
	client *ent.Client
}

func NewEventListenerManager(client *ent.Client) *EventListenerManager {
	return &EventListenerManager{client: client}
}

func (m *EventListenerManager) QueryAll(ctx context.Context) []*ent.EventListener {
	return m.client.EventListener.
		Query().
		AllX(ctx)
}

func (m *EventListenerManager) QueryAllWithChain(ctx context.Context) []*ent.EventListener {
	return m.client.EventListener.
		Query().
		WithChain().
		AllX(ctx)
}

func (m *EventListenerManager) QueryByUser(ctx context.Context, entUser *ent.User) []*ent.EventListener {
	return m.client.EventListener.
		Query().
		Where(
			eventlistener.HasUserWith(
				user.IDEQ(entUser.ID)),
		).
		WithChain().
		AllX(ctx)
}

type EventsCount []struct {
	EventType event.EventType `json:"event_type,omitempty"`
	Count     int             `json:"count,omitempty"`
}

func (m *EventListenerManager) QueryCountEventsByType(ctx context.Context, entUser *ent.User, isRead bool) (EventsCount, error) {
	var eventsCount = EventsCount{}
	err := m.client.Event.
		Query().
		Where(
			event.And(
				event.HasEventListenerWith(eventlistener.HasUserWith(user.IDEQ(entUser.ID))),
				event.NotifyTimeLTE(time.Now()),
			),
		).
		GroupBy(event.FieldEventType).
		Aggregate(ent.Count()).
		Scan(ctx, &eventsCount)
	return eventsCount, err
}

func (m *EventListenerManager) QueryEvents(ctx context.Context, el *ent.EventListener, eventType eventpb.EventType, startTime *timestamppb.Timestamp, endTime *timestamppb.Timestamp, limit int32, offset int64) ([]*ent.Event, error) {
	if startTime == nil {
		startTime = timestamppb.Now()
	}
	if endTime == nil {
		endTime = timestamppb.New(time.Now().AddDate(-1, 0, 0))
	}
	if limit == 0 {
		limit = 100
	}
	return el.
		QueryEvents().
		Where(
			// TODO: fix this
			event.EventTypeEQ(event.EventType(eventType.String())),
			//event.NotifyTimeGTE(startTime.AsTime()),
			event.NotifyTimeLTE(time.Now()),
		).
		Offset(int(offset)).
		Limit(int(limit)).
		All(ctx)
}

func (m *EventListenerManager) UpdateAddEvent(
	ctx context.Context,
	el *ent.EventListener,
	eventType event.EventType,
	dataType event.DataType,
	notifyTime time.Time,
	data []byte,
	isTxEvent bool,
) (*ent.Event, error) {
	return m.client.Event.
		Create().
		SetEventListener(el).
		SetNotifyTime(notifyTime).
		SetEventType(eventType).
		SetDataType(dataType).
		SetData(data).
		SetIsTxEvent(isTxEvent).
		Save(ctx)
}

func (m *EventListenerManager) Create(
	ctx context.Context,
	entUser *ent.User,
	entChain *ent.Chain,
	walletAddress string,
) (*ent.EventListener, error) {
	return m.client.EventListener.
		Create().
		SetUser(entUser).
		SetChain(entChain).
		SetWalletAddress(walletAddress).
		Save(ctx)
}
