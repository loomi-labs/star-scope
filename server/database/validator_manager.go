package database

import (
	"context"
	cosmossdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/shifty11/go-logger/log"
	"time"
)

const timeUntilConsideredInactive = 24 * time.Hour // TODO: set this to 1 month or so

type ValidatorManager struct {
	client *ent.Client
}

func NewValidatorManager(client *ent.Client) *ValidatorManager {
	return &ValidatorManager{client: client}
}

func (manager *ValidatorManager) StartTx(ctx context.Context) (context.Context, error) {
	return startTx(ctx, manager.client)
}

func (manager *ValidatorManager) getAccountAddress(operatorAddress string, chainEnt *ent.Chain) (string, error) {
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

func (manager *ValidatorManager) Create(
	ctx context.Context,
	chainEnt *ent.Chain,
	operatorAddress string,
	moniker string,
	isActive bool,
) (*ent.Validator, error) {
	accountAddress, err := manager.getAccountAddress(operatorAddress, chainEnt)
	if err != nil {
		log.Sugar.Errorf("Error while getting account address for validator %v: %v", operatorAddress, err)
		return nil, err
	}
	var firstInactiveTime *time.Time
	if !isActive {
		var now = time.Now()
		firstInactiveTime = &now
	}
	return manager.client.Validator.
		Create().
		SetChain(chainEnt).
		SetOperatorAddress(operatorAddress).
		SetAddress(accountAddress).
		SetMoniker(moniker).
		SetNillableFirstInactiveTime(firstInactiveTime).
		Save(ctx)
}

func (manager *ValidatorManager) Update(ctx context.Context, validatorEnt *ent.Validator, moniker string, isActive bool) error {
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
