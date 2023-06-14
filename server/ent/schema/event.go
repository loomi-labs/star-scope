package schema

import (
	"database/sql/driver"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	eventpb "github.com/loomi-labs/star-scope/event"
	kafkaevent "github.com/loomi-labs/star-scope/event"
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

	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Enum("event_type").
			Values(eventTypes...),
		field.Bytes("chain_event").
			Optional().
			GoType(&ChainEventWithScan{}),
		field.Bytes("contract_event").
			Optional().
			GoType(&ContractEventWithScan{}),
		field.Bytes("wallet_event").
			Optional().
			GoType(&WalletEventWithScan{}),
		field.Enum("data_type").
			Values(getDataTypes()...),
		field.Time("notify_time").
			Default(time.Now()),
		field.Bool("is_read").
			Default(false),
		field.Bool("is_background").
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

func getDataTypes() []string {
	var dataTypes []string
	var events = []interface{}{
		kafkaevent.WalletEvent_CoinReceived{},
		kafkaevent.WalletEvent_OsmosisPoolUnlock{},
		kafkaevent.WalletEvent_Unstake{},
		kafkaevent.WalletEvent_NeutronTokenVesting{},
		kafkaevent.WalletEvent_Voted{},
		kafkaevent.WalletEvent_VoteReminder{},
		kafkaevent.ChainEvent_ValidatorOutOfActiveSet{},
		kafkaevent.ChainEvent_ValidatorSlash{},
	}
	for _, t := range events {
		dataTypes = append(dataTypes, reflect.TypeOf(t).Name())
	}
	var govBaseEvents = []interface{}{
		kafkaevent.ChainEvent_GovernanceProposal{},
		kafkaevent.ContractEvent_ContractGovernanceProposal{},
	}
	var govEvents = []string{"Ongoing", "Finished"}
	for _, baseEvent := range govBaseEvents {
		var govBase = reflect.TypeOf(baseEvent).Name()
		for _, govEvent := range govEvents {
			dataTypes = append(dataTypes, fmt.Sprintf("%s_%s", govBase, govEvent))
		}
	}
	return dataTypes
}

type ChainEventWithScan struct {
	*kafkaevent.ChainEvent
}

func (x *ChainEventWithScan) Value() (driver.Value, error) {
	return proto.Marshal(x)
}

func (x *ChainEventWithScan) Scan(src any) error {
	if src == nil {
		return nil
	}
	if b, ok := src.([]byte); ok {
		var chainEvent kafkaevent.ChainEvent
		if err := proto.Unmarshal(b, &chainEvent); err != nil {
			return err
		}
		x.ChainEvent = &chainEvent
		return nil
	}
	return fmt.Errorf("unexpected type %T", src)
}

type ContractEventWithScan struct {
	*kafkaevent.ContractEvent
}

func (x *ContractEventWithScan) Value() (driver.Value, error) {
	return proto.Marshal(x)
}

func (x *ContractEventWithScan) Scan(src any) error {
	if src == nil {
		return nil
	}
	if b, ok := src.([]byte); ok {
		var contractEvent kafkaevent.ContractEvent
		if err := proto.Unmarshal(b, &contractEvent); err != nil {
			return err
		}
		x.ContractEvent = &contractEvent
		return nil
	}
	return fmt.Errorf("unexpected type %T", src)
}

type WalletEventWithScan struct {
	*kafkaevent.WalletEvent
}

func (x *WalletEventWithScan) Value() (driver.Value, error) {
	return proto.Marshal(x)
}

func (x *WalletEventWithScan) Scan(src any) error {
	if src == nil {
		return nil
	}
	if b, ok := src.([]byte); ok {
		var walletEvent kafkaevent.WalletEvent
		if err := proto.Unmarshal(b, &walletEvent); err != nil {
			return err
		}
		x.WalletEvent = &walletEvent
		return nil
	}
	return fmt.Errorf("unexpected type %T", src)
}
