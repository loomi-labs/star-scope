package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/chain"
	"github.com/loomi-labs/star-scope/ent/contractproposal"
	"github.com/loomi-labs/star-scope/ent/event"
	"github.com/loomi-labs/star-scope/ent/eventlistener"
	"github.com/loomi-labs/star-scope/ent/predicate"
	"github.com/loomi-labs/star-scope/ent/proposal"
	"github.com/loomi-labs/star-scope/ent/schema"
	"github.com/loomi-labs/star-scope/ent/user"
	kafkaevent "github.com/loomi-labs/star-scope/event"
	"github.com/loomi-labs/star-scope/kafka_internal"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
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
	return m.client.EventListener.
		Query().
		Where(eventlistener.DataTypeIn(dataTypes...)).
		AllX(ctx)
}

func (m *EventListenerManager) QueryWithChain(ctx context.Context, dataTypes ...eventlistener.DataType) []*ent.EventListener {
	return m.client.EventListener.
		Query().
		Where(eventlistener.DataTypeIn(dataTypes...)).
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

func (m *EventListenerManager) QueryCountEventsByType(ctx context.Context, entUser *ent.User, isRead bool, withBackgroundEvents bool) (*EventsCount, error) {
	var eventsCount = EventsCount{}
	var predicates = []predicate.Event{
		event.HasEventListenerWith(eventlistener.HasUserWith(user.IDEQ(entUser.ID))),
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

func getWalletEvents(entChain *ent.Chain) []eventlistener.DataType {
	dt := []eventlistener.DataType{
		eventlistener.DataTypeWalletEvent_CoinReceived,
	}
	if strings.Contains(entChain.Path, "neutron") {
		dt = append(dt, eventlistener.DataTypeWalletEvent_NeutronTokenVesting)
	} else {
		dt = append(dt, eventlistener.DataTypeWalletEvent_Unstake)
		dt = append(dt, eventlistener.DataTypeWalletEvent_Voted)
		dt = append(dt, eventlistener.DataTypeWalletEvent_VoteReminder)
	}
	if strings.Contains(entChain.Path, "osmosis") {
		dt = append(dt, eventlistener.DataTypeWalletEvent_OsmosisPoolUnlock)
	}
	return dt
}

func getChainEvents(entChain *ent.Chain) []eventlistener.DataType {
	var dt []eventlistener.DataType
	if strings.Contains(entChain.Path, "neutron") {
	} else {
		dt = append(dt, eventlistener.DataTypeChainEvent_GovernanceProposal_Ongoing)
		dt = append(dt, eventlistener.DataTypeChainEvent_GovernanceProposal_Finished)
		dt = append(dt, eventlistener.DataTypeChainEvent_ValidatorOutOfActiveSet)
	}
	return dt
}

func getContractEvents(entChain *ent.Chain) []eventlistener.DataType {
	var dt []eventlistener.DataType
	if strings.Contains(entChain.Path, "neutron") {
		dt = append(dt, eventlistener.DataTypeContractEvent_ContractGovernanceProposal_Ongoing)
		dt = append(dt, eventlistener.DataTypeContractEvent_ContractGovernanceProposal_Finished)
	}
	return dt
}

func (m *EventListenerManager) CreateBulk(
	ctx context.Context,
	entUser *ent.User,
	entChains []*ent.Chain,
	mainAddress string,
) ([]*ent.EventListener, error) {
	var bulk []*ent.EventListenerCreate
	for _, entChain := range entChains {
		walletAddress, err := common.ConvertWithOtherPrefix(mainAddress, entChain.Bech32Prefix)
		if err != nil {
			return nil, err
		}
		for _, dt := range getWalletEvents(entChain) {
			bulk = append(bulk, m.client.EventListener.
				Create().
				SetChain(entChain).
				SetUser(entUser).
				SetWalletAddress(walletAddress).
				SetDataType(dt))
		}
		for _, dt := range getChainEvents(entChain) {
			bulk = append(bulk, m.client.EventListener.
				Create().
				SetChain(entChain).
				SetUser(entUser).
				SetDataType(dt))
		}
		for _, dt := range getContractEvents(entChain) {
			bulk = append(bulk, m.client.EventListener.
				Create().
				SetChain(entChain).
				SetUser(entUser).
				SetDataType(dt))
		}
	}
	el, err := m.client.EventListener.
		CreateBulk(bulk...).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	m.kafkaInternal.ProduceDbChangeMsg(kafka_internal.EventListenerCreated)
	return el, nil
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
				event.HasEventListenerWith(eventlistener.HasUserWith(user.IDEQ(u.ID))),
				event.IDEQ(id),
			),
		).
		SetIsRead(true).
		Exec(ctx)
}
