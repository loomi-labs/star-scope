package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/event"
	"github.com/loomi-labs/star-scope/ent/eventlistener"
	"github.com/loomi-labs/star-scope/ent/predicate"
	"github.com/loomi-labs/star-scope/ent/schema"
	"github.com/loomi-labs/star-scope/ent/user"
	kafkaevent "github.com/loomi-labs/star-scope/event"
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
				event.IsRead(isRead),
			),
		).
		GroupBy(event.FieldEventType).
		Aggregate(ent.Count()).
		Scan(ctx, &eventsCount)
	return eventsCount, err
}

func (m *EventListenerManager) QueryEvents(ctx context.Context, el *ent.EventListener, eventType *eventpb.EventType, startTime *timestamppb.Timestamp, endTime *timestamppb.Timestamp, limit int32, offset int64) ([]*ent.Event, error) {
	if startTime == nil {
		startTime = timestamppb.Now()
	}
	if endTime == nil {
		endTime = timestamppb.New(time.Now())
	}
	if limit == 0 {
		limit = 100
	}
	var filters = []predicate.Event{
		event.NotifyTimeLTE(endTime.AsTime()),
	}
	if eventType != nil {
		filters = append(filters, event.EventTypeEQ(event.EventType(eventType.String())))
	}
	return el.
		QueryEvents().
		Where(filters...).
		Offset(int(offset)).
		Limit(int(limit)).
		All(ctx)
}

func (m *EventListenerManager) UpdateAddChainEvent(
	ctx context.Context,
	el *ent.EventListener,
	chainEvent *kafkaevent.ChainEvent,
	eventType event.EventType,
	dataType event.DataType,
) (*ent.Event, error) {
	var withScan = &schema.ChainEventWithScan{
		ChainEvent: chainEvent,
	}
	return m.client.Event.
		Create().
		SetEventListener(el).
		SetChainEvent(withScan).
		SetEventType(eventType).
		SetDataType(dataType).
		SetNotifyTime(chainEvent.NotifyTime.AsTime()).
		Save(ctx)
}

func (m *EventListenerManager) UpdateAddContractEvent(
	ctx context.Context,
	el *ent.EventListener,
	contractEvent *kafkaevent.ContractEvent,
	eventType event.EventType,
	dataType event.DataType,
) (*ent.Event, error) {
	var withScan = &schema.ContractEventWithScan{
		ContractEvent: contractEvent,
	}
	return m.client.Event.
		Create().
		SetEventListener(el).
		SetContractEvent(withScan).
		SetEventType(eventType).
		SetDataType(dataType).
		SetNotifyTime(contractEvent.NotifyTime.AsTime()).
		Save(ctx)
}

func (m *EventListenerManager) UpdateAddWalletEvent(
	ctx context.Context,
	el *ent.EventListener,
	walletEvent *kafkaevent.WalletEvent,
	eventType event.EventType,
	dataType event.DataType,
) (*ent.Event, error) {
	var withScan = &schema.WalletEventWithScan{
		WalletEvent: walletEvent,
	}
	return m.client.Event.
		Create().
		SetEventListener(el).
		SetWalletEvent(withScan).
		SetEventType(eventType).
		SetDataType(dataType).
		SetNotifyTime(walletEvent.NotifyTime.AsTime()).
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

func (m *EventListenerManager) UpdateMarkEventRead(ctx context.Context, u *ent.User, id uuid.UUID) error {
	return m.client.Event.
		Update().
		Where(
			event.And(
				event.HasEventListenerWith(eventlistener.HasUserWith(user.IDEQ(u.ID))),
				event.IDEQ(id),
			),
		).
		SetIsRead(true).
		Exec(ctx)
}
