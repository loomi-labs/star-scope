package database

import (
	"context"
	"github.com/shifty11/blocklog-backend/ent"
)

type EventListenerManager struct {
	client *ent.Client
}

func NewEventListenerManager(client *ent.Client) *EventListenerManager {
	return &EventListenerManager{client: client}
}

func (m *EventListenerManager) QueryAll(ctx context.Context) []*ent.EventListener {
	return m.client.EventListener.Query().AllX(ctx)
}
