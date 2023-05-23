package database

import (
	"context"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/event"
	"github.com/loomi-labs/star-scope/ent/eventlistener"
	"github.com/loomi-labs/star-scope/ent/user"
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

func (m *EventListenerManager) QueryByUser(ctx context.Context, entUser *ent.User) []*ent.EventListener {
	return m.client.EventListener.
		Query().
		Where(
			eventlistener.HasUserWith(
				user.IDEQ(entUser.ID)),
		).
		WithEvents().
		WithChain().
		AllX(ctx)
}

func (m *EventListenerManager) UpdateAddEvent(
	ctx context.Context,
	el *ent.EventListener,
	eventType event.Type,
	notifyTime time.Time,
	txEvent []byte,
) (*ent.Event, error) {
	return m.client.Event.
		Create().
		SetEventListener(el).
		SetType(eventType).
		SetNotifyTime(notifyTime).
		SetTxEvent(txEvent).
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
