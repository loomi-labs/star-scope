package validator_crawler

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/ent"
	kafkaevent "github.com/loomi-labs/star-scope/event"
	"github.com/loomi-labs/star-scope/kafka_internal"
	"github.com/loomi-labs/star-scope/types"
	"github.com/robfig/cron/v3"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"time"
)

const urlValidators = "%v/cosmos/staking/v1beta1/validators?pagination.limit=1000"
const urlValidatorSet = "%v/cosmos/base/tendermint/v1beta1/validatorsets/latest"
const urlValidatorSlashes = "%v/cosmos/distribution/v1beta1/validators/%v/slashes"

type ValidatorCrawler struct {
	chainManager     *database.ChainManager
	validatorManager *database.ValidatorManager
	kafkaInternal    kafka_internal.KafkaInternal
}

func NewValidatorCrawler(dbManagers *database.DbManagers, kafkaInternal kafka_internal.KafkaInternal) *ValidatorCrawler {
	return &ValidatorCrawler{
		chainManager:     dbManagers.ChainManager,
		validatorManager: dbManagers.ValidatorManager,
		kafkaInternal:    kafkaInternal,
	}
}

func validatorNeedsUpdate(validatorEnt *ent.Validator, data *types.Validator, isValidatorInActiveSet bool) bool {
	return validatorEnt.Moniker != data.Description.Moniker ||
		(validatorEnt.FirstInactiveTime != nil && isValidatorInActiveSet) ||
		(validatorEnt.FirstInactiveTime == nil && !isValidatorInActiveSet)
}

func isValidatorValid(data *types.Validator) bool {
	return data.OperatorAddress != ""
}

func getExistingValidator(validators []*ent.Validator, validator *types.Validator) *ent.Validator {
	for _, val := range validators {
		if val.OperatorAddress == validator.OperatorAddress {
			return val
		}
	}
	return nil
}

// isValidatorInActiveSet compares a pubkey of a validator with the pubkeys of the validators in the active set.
// We can not use the address from the active set because it is a `valcons` address which we would have to convert first.
func isValidatorInActiveSet(pubKey string, activeValidatorSet []types.ValidatorSetValidator) bool {
	for _, validator := range activeValidatorSet {
		if pubKey == validator.PubKey.Key {
			return true
		}
	}
	return false
}

func (c *ValidatorCrawler) createOutOfActiveSetEvent(chain *ent.Chain, validator *ent.Validator) ([]byte, error) {
	var now = timestamppb.Now()
	chainEvent := &kafkaevent.ChainEvent{
		ChainId:    uint64(chain.ID),
		Timestamp:  now,
		NotifyTime: timestamppb.New(time.Now().Add(time.Hour * 24)), // notify if validator is still out of active set after 24 hours
		Event: &kafkaevent.ChainEvent_ValidatorOutOfActiveSet{
			ValidatorOutOfActiveSet: &kafkaevent.ValidatorOutOfActiveSetEvent{
				ValidatorAddress:         validator.Address,
				ValidatorOperatorAddress: validator.OperatorAddress,
				ValidatorMoniker:         validator.Moniker,
				FirstInactiveTime:        now,
			},
		},
	}
	pbEvent, err := proto.Marshal(chainEvent)
	if err != nil {
		return nil, err
	}
	return pbEvent, nil
}

func (c *ValidatorCrawler) reportErrorIfNecessary(chain *ent.Chain, url string, err error) {
	// If the last successful validator query was less than 3 days ago, don't report the error
	if chain.LastSuccessfulValidatorQuery != nil && chain.LastSuccessfulValidatorQuery.Add(24*time.Hour*3).After(time.Now()) {
		log.Sugar.Debugf("error calling %v on %v: %v", url, chain.PrettyName, err)
	} else {
		log.Sugar.Errorf("error calling %v on %v: %v", url, chain.PrettyName, err)
	}
}

func (c *ValidatorCrawler) addOrUpdateValidators() {
	log.Sugar.Info("Getting all validators")
	var ctx = context.Background()
	var pbEvents [][]byte
	for _, entChain := range c.chainManager.QueryIsQuerying(ctx) {
		if strings.Contains(entChain.Path, "neutron") {
			continue
		}

		log.Sugar.Infof("Getting validators for chain %v", entChain.PrettyName)
		url := fmt.Sprintf(urlValidators, entChain.RestEndpoint)
		var validatorsResponse types.ValidatorsResponse
		_, err := common.GetJson(url, 5, &validatorsResponse)
		if err != nil {
			c.reportErrorIfNecessary(entChain, url, err)
			continue
		}
		if validatorsResponse.Pagination.Total != "0" {
			log.Sugar.Errorf("pagination is not implemented yet")
		}

		var validatorSetResponse types.ValidatorSetResponse
		url = fmt.Sprintf(urlValidatorSet, entChain.RestEndpoint)
		_, err = common.GetJson(url, 5, &validatorSetResponse)
		if err != nil {
			c.reportErrorIfNecessary(entChain, url, err)
			continue
		}

		existingValidators, err := entChain.QueryValidators().All(ctx)
		if err != nil {
			log.Sugar.Panicf("error getting validators for chain %v: %v", entChain.PrettyName, err)
		}

		for _, validator := range validatorsResponse.Validators {
			if isValidatorValid(&validator) {
				var isInActiveSet = isValidatorInActiveSet(validator.ConsensusPubkey.Key, validatorSetResponse.Validators)
				var existingValidator = getExistingValidator(existingValidators, &validator)
				if existingValidator != nil {
					if validatorNeedsUpdate(existingValidator, &validator, isInActiveSet) {
						log.Sugar.Infof("Updating validator %v %v", validator.OperatorAddress, validator.Description.Moniker)
						err := c.validatorManager.Update(ctx, existingValidator, validator.Description.Moniker, isInActiveSet)
						if err != nil {
							log.Sugar.Errorf("error updating validator %v: %v", existingValidator.Address, err)
							continue
						}
						if !isInActiveSet {
							log.Sugar.Infof("Validator %v %v is out of active set", validator.OperatorAddress, validator.Description.Moniker)
							event, err := c.createOutOfActiveSetEvent(entChain, existingValidator)
							if err != nil {
								log.Sugar.Errorf("error creating out of active set event for validator %v: %v", existingValidator.Address, err)
								continue
							}
							pbEvents = append(pbEvents, event)
						} else {
							log.Sugar.Debugf("Validator %v %v is back in active set", validator.OperatorAddress, validator.Description.Moniker)
							cnt, err := c.validatorManager.DeleteOutOfActiveSetEvents(ctx, existingValidator)
							if err != nil {
								log.Sugar.Errorf("error deleting out of active set events for validator %v: %v", existingValidator.Address, err)
								continue
							}
							log.Sugar.Debugf("Deleted %v out-of-active-set events for validator %v", cnt, existingValidator.Address)
						}
					}
				} else {
					log.Sugar.Infof("Creating validator %v %v", validator.OperatorAddress, validator.Description.Moniker)
					_, err = c.validatorManager.Create(ctx, entChain, validator.OperatorAddress, validator.Description.Moniker, isInActiveSet)
					if err != nil {
						log.Sugar.Errorf("error creating validator %v: %v", validator.OperatorAddress, err)
						continue
					}
				}
			} else {
				log.Sugar.Debugf("Validator %v %v is invalid", validator.OperatorAddress, validator.Description.Moniker)
			}
		}
		c.chainManager.UpdateSetLastSuccessfulValidatorQuery(ctx, entChain)
	}
	if len(pbEvents) > 0 {
		log.Sugar.Debugf("Send %v out-of-active-set events", len(pbEvents))
		c.kafkaInternal.ProduceChainEvents(pbEvents)
	}
}

func (c *ValidatorCrawler) createSlashEvent(chain *ent.Chain, validator *ent.Validator, slashEvent types.SlashEvent) ([]byte, error) {
	var now = timestamppb.Now()
	chainEvent := &kafkaevent.ChainEvent{
		ChainId:    uint64(chain.ID),
		Timestamp:  now,
		NotifyTime: now,
		Event: &kafkaevent.ChainEvent_ValidatorSlash{
			ValidatorSlash: &kafkaevent.ValidatorSlashEvent{
				ValidatorAddress:         validator.Address,
				ValidatorOperatorAddress: validator.OperatorAddress,
				ValidatorMoniker:         validator.Moniker,
				ValidatorPeriod:          slashEvent.ValidatorPeriod,
				Fraction:                 slashEvent.Fraction,
			},
		},
	}
	pbEvent, err := proto.Marshal(chainEvent)
	if err != nil {
		return nil, err
	}
	return pbEvent, nil
}

func (c *ValidatorCrawler) fetchSlashEvents() {
	log.Sugar.Info("Fetching slash events")
	var pbEvents [][]byte
	var validators = c.validatorManager.QueryActive(context.Background())
	for _, validator := range validators {
		if strings.Contains(validator.Edges.Chain.Path, "neutron") {
			continue
		}

		log.Sugar.Debugf("Fetching slash events for validator %v", validator.Address)

		var validatorSetResponse types.ValidatorSlashResponse
		url := fmt.Sprintf(urlValidatorSlashes, validator.Edges.Chain.RestEndpoint, validator.OperatorAddress)
		_, err := common.GetJson(url, 5, &validatorSetResponse)
		if err != nil {
			log.Sugar.Errorf("error calling %v: %v", url, err)
			continue
		}

		for _, slashEvent := range validatorSetResponse.Slashes {
			if validator.LastSlashValidatorPeriod != nil && *validator.LastSlashValidatorPeriod == slashEvent.ValidatorPeriod {
				continue
			}
			err := c.validatorManager.UpdateSetSlashed(context.Background(), validator, slashEvent.ValidatorPeriod)
			if err != nil {
				log.Sugar.Errorf("error updating validator %v: %v", validator.Address, err)
				continue
			}

			pbEvent, err := c.createSlashEvent(validator.Edges.Chain, validator, slashEvent)
			if err != nil {
				log.Sugar.Errorf("error creating slash event for validator %v: %v", validator.Address, err)
				continue
			}
			pbEvents = append(pbEvents, pbEvent)
		}
		time.Sleep(100 * time.Millisecond)
	}
	if len(pbEvents) > 0 {
		log.Sugar.Debugf("Send %v slashing events", len(pbEvents))
		c.kafkaInternal.ProduceChainEvents(pbEvents)
	}
}

func (c *ValidatorCrawler) StartCrawling() {
	c.addOrUpdateValidators()
	log.Sugar.Info("Scheduling validator crawl")
	cr := cron.New()
	_, err := cr.AddFunc("0 10 * * *", func() { c.addOrUpdateValidators() }) // every day at 10:00
	if err != nil {
		log.Sugar.Errorf("while executing 'addOrUpdateValidators' via cron: %v", err)
	}
	_, err = cr.AddFunc("0 11 * * *", func() { c.fetchSlashEvents() }) // every day at 11:00
	if err != nil {
		log.Sugar.Errorf("while executing 'fetchSlashEvents' via cron: %v", err)
	}
	cr.Start()
	select {}
}
