package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/chain"
	"github.com/loomi-labs/star-scope/ent/commchannel"
	"github.com/loomi-labs/star-scope/ent/contractproposal"
	"github.com/loomi-labs/star-scope/ent/event"
	"github.com/loomi-labs/star-scope/ent/eventlistener"
	"github.com/loomi-labs/star-scope/ent/predicate"
	"github.com/loomi-labs/star-scope/ent/proposal"
	"github.com/loomi-labs/star-scope/ent/schema"
	"github.com/loomi-labs/star-scope/ent/state"
	"github.com/loomi-labs/star-scope/ent/user"
	kafkaevent "github.com/loomi-labs/star-scope/event"
	"github.com/loomi-labs/star-scope/kafka_internal"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type EventListenerManager struct {
	client        *ent.Client
	kafkaInternal kafka_internal.KafkaInternal
}

func NewEventListenerManager(client *ent.Client, kafkaInternal kafka_internal.KafkaInternal) *EventListenerManager {
	return &EventListenerManager{client: client, kafkaInternal: kafkaInternal}
}

func (m *EventListenerManager) Query(ctx context.Context, dataTypes ...eventlistener.DataType) []*ent.EventListener {
	var predicates []predicate.EventListener
	if len(dataTypes) == 0 {
		predicates = append(predicates, eventlistener.DataTypeIn(dataTypes...))
	}
	return m.client.EventListener.
		Query().
		Where(predicates...).
		AllX(ctx)
}

func (m *EventListenerManager) QueryWithChain(ctx context.Context, dataTypes ...eventlistener.DataType) []*ent.EventListener {
	var predicates []predicate.EventListener
	if len(dataTypes) > 0 {
		predicates = append(predicates, eventlistener.DataTypeIn(dataTypes...))
	}
	return m.client.EventListener.
		Query().
		Where(predicates...).
		WithChain().
		AllX(ctx)
}

func (m *EventListenerManager) QueryByUser(ctx context.Context, entUser *ent.User) []*ent.EventListener {
	return m.client.EventListener.
		Query().
		Where(
			eventlistener.HasUsersWith(
				user.IDEQ(entUser.ID)),
		).
		WithChain().
		AllX(ctx)
}

type EventsCount []struct {
	EventType event.EventType `json:"event_type,omitempty"`
	Count     int             `json:"count,omitempty"`
}

func (m *EventListenerManager) QueryCountEventsByType(ctx context.Context, entUser *ent.User, isRead bool, withBackgroundEvents bool) (*EventsCount, error) {
	var eventsCount = EventsCount{}
	var predicates = []predicate.Event{
		event.HasEventListenerWith(eventlistener.HasUsersWith(user.IDEQ(entUser.ID))),
		event.NotifyTimeLTE(time.Now()),
		event.IsRead(isRead),
	}
	if !withBackgroundEvents {
		predicates = append(predicates, event.IsBackground(false))
	}
	err := m.client.Event.
		Query().
		Where(
			event.And(predicates...),
		).
		GroupBy(event.FieldEventType).
		Aggregate(ent.Count()).
		Scan(ctx, &eventsCount)
	return &eventsCount, err
}

func (m *EventListenerManager) QueryEvents(ctx context.Context, el *ent.EventListener, eventType *kafkaevent.EventType, startTime *timestamppb.Timestamp, endTime *timestamppb.Timestamp, limit int32, offset int64) ([]*ent.Event, error) {
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
		filters = append(filters, event.IsBackgroundEQ(false))
	}
	return el.
		QueryEvents().
		Where(filters...).
		Offset(int(offset)).
		Limit(int(limit)).
		All(ctx)
}

func (m *EventListenerManager) QueryEventsSince(ctx context.Context, startTime time.Time, endTime time.Time, entity state.Entity) ([]*ent.Event, error) {
	var predicates = []predicate.Event{
		event.NotifyTimeGTE(startTime),
		event.NotifyTimeLTE(endTime),
		event.IsBackgroundEQ(false),
	}
	var eventListenerQueries = []func(q *ent.EventListenerQuery){
		func(q *ent.EventListenerQuery) {
			q.WithChain()
		},
	}
	if entity == state.EntityDiscord {
		predicates = append(predicates, event.HasEventListenerWith(eventlistener.HasCommChannelsWith(commchannel.TypeEQ(commchannel.TypeDiscord))))
		eventListenerQueries = append(eventListenerQueries, func(q *ent.EventListenerQuery) {
			q.WithCommChannels(func(q *ent.CommChannelQuery) {
				q.Where(commchannel.TypeEQ(commchannel.TypeDiscord))
			})
		})
	} else if entity == state.EntityTelegram {
		predicates = append(predicates, event.HasEventListenerWith(eventlistener.HasCommChannelsWith(commchannel.TypeEQ(commchannel.TypeTelegram))))
		eventListenerQueries = append(eventListenerQueries, func(q *ent.EventListenerQuery) {
			q.WithCommChannels(func(q *ent.CommChannelQuery) {
				q.Where(commchannel.TypeEQ(commchannel.TypeTelegram))
			})
		})
	}
	return m.client.Event.
		Query().
		Where(
			event.And(predicates...),
		).
		Order(ent.Desc(event.FieldDataType)).
		WithEventListener(eventListenerQueries...).
		All(ctx)
}

func (m *EventListenerManager) QueryNotifierState(ctx context.Context, entity state.Entity) (*ent.State, error) {
	return m.client.State.
		Query().
		Where(
			state.EntityEQ(entity),
		).
		Only(ctx)
}

func (m *EventListenerManager) UpdateNotifierState(ctx context.Context, state *ent.State, updatetime time.Time) (*ent.State, error) {
	return state.
		Update().
		SetLastEventTime(updatetime).
		Save(ctx)
}

type VoteReminder struct {
	Chain         *ent.Chain
	EventListener *ent.EventListener
	Proposal      *ent.Proposal
}

// QueryForVoteReminderAddresses returns a list of addresses that should be checked if they have voted
// It filters out addresses that have already voted or have been reminded
// It only returns addresses that have a proposal that is in the voting period and ends in less than 24 hours
func (m *EventListenerManager) QueryForVoteReminderAddresses(ctx context.Context) ([]*VoteReminder, error) {
	oneDayInTheFuture := time.Now().Add(time.Hour * 24)
	els, err := m.client.EventListener.
		Query().
		Where(eventlistener.And(
			eventlistener.DataTypeEQ(eventlistener.DataTypeWalletEvent_VoteReminder),
			eventlistener.WalletAddressNEQ(""),
			eventlistener.Or(
				eventlistener.HasChainWith(chain.HasProposalsWith(proposal.And(
					proposal.StatusEQ(proposal.StatusPROPOSAL_STATUS_VOTING_PERIOD),
					proposal.VotingEndTimeLTE(oneDayInTheFuture),
				))),
				eventlistener.HasChainWith(chain.HasContractProposalsWith(contractproposal.And(
					contractproposal.StatusEQ(contractproposal.StatusOPEN),
					contractproposal.VotingEndTimeLTE(oneDayInTheFuture),
				))),
			),
		)).
		WithChain(func(q *ent.ChainQuery) {
			q.WithProposals(func(q *ent.ProposalQuery) {
				q.Where(proposal.StatusEQ(proposal.StatusPROPOSAL_STATUS_VOTING_PERIOD))
				q.Where(proposal.VotingEndTimeLTE(oneDayInTheFuture))
			})
			q.WithContractProposals(func(q *ent.ContractProposalQuery) {
				q.Where(contractproposal.StatusEQ(contractproposal.StatusOPEN))
				q.Where(contractproposal.VotingEndTimeLTE(oneDayInTheFuture))
			})
		}).
		All(ctx)
	if err != nil {
		return nil, err
	}
	var voteReminders []*VoteReminder
	for _, el := range els {
		votedEvents := m.client.Event.
			Query().
			Where(event.And(
				event.DataTypeEQ(event.DataTypeWalletEvent_Voted),
				event.HasEventListenerWith(eventlistener.WalletAddressEQ(el.WalletAddress)),
			)).
			AllX(ctx)
		voteReminderEvents := m.client.Event.
			Query().
			Where(event.And(
				event.DataTypeEQ(event.DataTypeWalletEvent_VoteReminder),
				event.HasEventListenerWith(eventlistener.WalletAddressEQ(el.WalletAddress)),
			)).
			AllX(ctx)
		for _, p := range el.Edges.Chain.Edges.Proposals {
			hasVoted := false
			for _, e := range votedEvents {
				if e.WalletEvent.GetVoted().GetProposalId() == p.ProposalID {
					hasVoted = true
					break
				}
			}
			hasBeenReminded := false
			for _, e := range voteReminderEvents {
				if e.WalletEvent.GetVoteReminder().GetProposalId() == p.ProposalID {
					hasBeenReminded = true
					break
				}
			}
			if !hasVoted && !hasBeenReminded {
				voteReminders = append(voteReminders, &VoteReminder{
					Chain:         el.Edges.Chain,
					EventListener: el,
					Proposal:      p,
				})
			}
		}
	}
	return voteReminders, nil
}

func (m *EventListenerManager) QuerySubscriptionsCountForTelegramChat(ctx context.Context, chatId int64) int {
	cnt, err := m.client.EventListener.
		Query().
		Where(eventlistener.HasCommChannelsWith(commchannel.TelegramChatIDEQ(chatId))).
		Count(ctx)
	if err != nil {
		log.Sugar.Errorf("Could not count subscriptions for telegram chat: %v", err)
	}
	return cnt
}

func (m *EventListenerManager) QuerySubscriptionsCountForDiscordChannel(ctx context.Context, channelId int64) int {
	cnt, err := m.client.EventListener.
		Query().
		Where(eventlistener.HasCommChannelsWith(commchannel.DiscordChannelIDEQ(channelId))).
		Count(ctx)
	if err != nil {
		log.Sugar.Errorf("Could not count subscriptions for discord channel: %v", err)
	}
	return cnt
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
	isBackground bool,
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
		SetIsBackground(isBackground).
		SetNotifyTime(walletEvent.NotifyTime.AsTime()).
		Save(ctx)
}

func (m *EventListenerManager) UpdateMarkEventRead(ctx context.Context, u *ent.User, id uuid.UUID) error {
	return m.client.Event.
		Update().
		Where(
			event.And(
				event.HasEventListenerWith(eventlistener.HasUsersWith(user.IDEQ(u.ID))),
				event.IDEQ(id),
			),
		).
		SetIsRead(true).
		Exec(ctx)
}
