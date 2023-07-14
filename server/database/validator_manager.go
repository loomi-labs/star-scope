package database

import (
	"context"
	cosmossdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/google/uuid"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/chain"
	"github.com/loomi-labs/star-scope/ent/event"
	"github.com/loomi-labs/star-scope/ent/eventlistener"
	"github.com/loomi-labs/star-scope/ent/validator"
	"github.com/shifty11/go-logger/log"
	"time"
)

const timeUntilConsideredInactive = time.Hour * 24

type ValidatorManager struct {
	client *ent.Client
}

func NewValidatorManager(client *ent.Client) *ValidatorManager {
	return &ValidatorManager{client: client}
}

func (m *ValidatorManager) getAccountAddress(operatorAddress string, chainEnt *ent.Chain) (string, error) {
	_, valAddr, err := bech32.DecodeAndConvert(operatorAddress)
	if err != nil {
		return "", err
	}
	accAddr, err := cosmossdk.Bech32ifyAddressBytes(chainEnt.Bech32Prefix, valAddr)
	if err != nil {
		return "", err
	}
	return accAddr, nil
}

func (m *ValidatorManager) Create(
	ctx context.Context,
	chainEnt *ent.Chain,
	operatorAddress string,
	moniker string,
	isActive bool,
) (*ent.Validator, error) {
	accountAddress, err := m.getAccountAddress(operatorAddress, chainEnt)
	if err != nil {
		log.Sugar.Errorf("Error while getting account address for validator %v: %v", operatorAddress, err)
		return nil, err
	}
	var firstInactiveTime *time.Time
	if !isActive {
		var now = time.Now().Add(-timeUntilConsideredInactive)
		firstInactiveTime = &now
	}
	return m.client.Validator.
		Create().
		SetChain(chainEnt).
		SetOperatorAddress(operatorAddress).
		SetAddress(accountAddress).
		SetMoniker(moniker).
		SetNillableFirstInactiveTime(firstInactiveTime).
		Save(ctx)
}

func (m *ValidatorManager) Update(ctx context.Context, validatorEnt *ent.Validator, moniker string, isActive bool) error {
	updateQuery := validatorEnt.Update()
	if isActive {
		updateQuery.ClearFirstInactiveTime()
	} else {
		updateQuery.SetFirstInactiveTime(time.Now())
	}
	return updateQuery.
		SetMoniker(moniker).
		Exec(ctx)
}

func (m *ValidatorManager) QueryActive(ctx context.Context) []*ent.Validator {
	return m.client.Validator.
		Query().
		Where(validator.Or(
			validator.FirstInactiveTimeIsNil(),
			validator.FirstInactiveTimeGT(time.Now().Add(-timeUntilConsideredInactive)),
		)).
		WithChain().
		Order(ent.Asc(validator.FieldMoniker)). // order by moniker so that the chain is random -> avoid spamming the same chain when querying validators
		AllX(ctx)
}

func (m *ValidatorManager) QueryByAddress(ctx context.Context, chainId int, address string) (*ent.Validator, error) {
	return m.client.Validator.
		Query().
		Where(validator.And(
			validator.AddressEQ(address),
			validator.HasChainWith(chain.IDEQ(chainId)),
		)).
		WithChain().
		First(ctx)
}

type ValidatorBundle struct {
	Validators []*ent.Validator
	Moniker    string
	LogoUrl    string
}

func (m *ValidatorManager) QueryActiveBundledByMoniker(ctx context.Context) []*ValidatorBundle {
	vals := m.client.Validator.
		Query().
		Where(validator.Or(
			validator.FirstInactiveTimeIsNil(),
			validator.FirstInactiveTimeGT(time.Now().Add(-timeUntilConsideredInactive)),
		)).
		Order(ent.Asc(validator.FieldMoniker)).
		AllX(ctx)
	var bundles []*ValidatorBundle
	var currentMoniker string
	for _, val := range vals {
		if currentMoniker != val.Moniker {
			currentMoniker = val.Moniker
			bundles = append(bundles, &ValidatorBundle{
				Validators: []*ent.Validator{val},
				Moniker:    val.Moniker,
			})
		} else {
			bundles[len(bundles)-1].Validators = append(bundles[len(bundles)-1].Validators, val)
		}
	}
	return bundles
}

func (m *ValidatorManager) UpdateSetSlashed(ctx context.Context, val *ent.Validator, period uint64) error {
	return val.Update().
		SetLastSlashValidatorPeriod(period).
		Exec(ctx)
}

func (m *ValidatorManager) DeleteOutOfActiveSetEvents(ctx context.Context, val *ent.Validator) (int, error) {
	events := val.
		QueryChain().
		QueryEventListeners().
		Where(eventlistener.DataTypeEQ(eventlistener.DataTypeChainEvent_ValidatorOutOfActiveSet)).
		QueryEvents().
		Where(event.NotifyTimeGT(time.Now())). // only cancel events that haven't been notified yet
		AllX(ctx)
	var toBeDeletedEvents []uuid.UUID
	for _, e := range events {
		if e.ChainEvent.GetValidatorOutOfActiveSet().GetValidatorAddress() == val.Address {
			toBeDeletedEvents = append(toBeDeletedEvents, e.ID)
		}
	}
	if len(toBeDeletedEvents) > 0 {
		return m.client.Event.Delete().
			Where(event.IDIn(toBeDeletedEvents...)).
			Exec(ctx)
	}
	return 0, nil
}
