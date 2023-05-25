package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"fmt"
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
	var types []string
	var events = []interface{}{
		indexevent.TxEvent_CoinReceived{},
		indexevent.TxEvent_OsmosisPoolUnlock{},
		indexevent.TxEvent_Unstake{},
	}
	for _, t := range events {
		types = append(types, reflect.TypeOf(t).Name())
	}
	var govBase = reflect.TypeOf(queryevent.QueryEvent_GovernanceProposal{}).Name()
	var govEvents = []string{"Ongoing", "Finished"}
	for _, govEvent := range govEvents {
		types = append(types, fmt.Sprintf("%s_%s", govBase, govEvent))
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
