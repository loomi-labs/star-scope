package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/shifty11/blocklog-backend/indexevent"
	"reflect"
	"time"
)

// Event holds the schema definition for the Event entity.
type Event struct {
	ent.Schema
}

func (Event) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Fields of the Event.
func (Event) Fields() []ent.Field {
	var types []string
	var events = []interface{}{
		indexevent.TxEvent_CoinReceived{},
		indexevent.TxEvent_OsmosisPoolUnlock{},
	}
	for _, t := range events {
		types = append(types, reflect.TypeOf(t).Name())
	}
	return []ent.Field{
		field.Enum("type").
			Values(types...),
		field.Bytes("tx_event"),
		field.Time("notify_time").
			Default(time.Now()),
	}
}

// Edges of the Event.
func (Event) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event_listener", EventListener.Type).
			Ref("events").
			Unique(),
	}
}
