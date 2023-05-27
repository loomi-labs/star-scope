package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"fmt"
	"github.com/loomi-labs/star-scope/grpc/event/eventpb"
	"github.com/loomi-labs/star-scope/indexevent"
	"github.com/loomi-labs/star-scope/queryevent"
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
	var eventTypes []string
	for _, t := range eventpb.EventType_name {
		eventTypes = append(eventTypes, t)
	}

	var dataTypes []string
	var events = []interface{}{
		indexevent.TxEvent_CoinReceived{},
		indexevent.TxEvent_OsmosisPoolUnlock{},
		indexevent.TxEvent_Unstake{},
	}
	for _, t := range events {
		dataTypes = append(dataTypes, reflect.TypeOf(t).Name())
	}
	var govBase = reflect.TypeOf(queryevent.QueryEvent_GovernanceProposal{}).Name()
	var govEvents = []string{"Ongoing", "Finished"}
	for _, govEvent := range govEvents {
		dataTypes = append(dataTypes, fmt.Sprintf("%s_%s", govBase, govEvent))
	}
	return []ent.Field{
		field.Enum("event_type").
			Values(eventTypes...),
		field.Bytes("data"),
		field.Enum("data_type").
			Values(dataTypes...),
		field.Bool("is_tx_event"),
		field.Time("notify_time").
			Default(time.Now()),
		field.Bool("is_read").
			Default(false),
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
