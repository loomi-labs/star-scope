package auth

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/eventlistener"
	kafkaevent "github.com/loomi-labs/star-scope/event"
	"github.com/loomi-labs/star-scope/kafka_internal"
	"github.com/loomi-labs/star-scope/types"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

const urlUnstaking = "%v/cosmos/staking/v1beta1/delegators/%v/unbonding_delegations"

type SetupCrawler struct {
	kafkaInternal kafka_internal.KafkaInternal
}

func NewSetupCrawler(kafkaInternal kafka_internal.KafkaInternal) *SetupCrawler {
	return &SetupCrawler{
		kafkaInternal: kafkaInternal,
	}
}

func (c *SetupCrawler) createEvent(chain *ent.Chain, account string, entry types.UnstakingEntry) ([]byte, error) {
	var now = timestamppb.Now()
	txEvent := &kafkaevent.WalletEvent{
		ChainId:       uint64(chain.ID),
		WalletAddress: account,
		Timestamp:     now,
		NotifyTime:    timestamppb.New(entry.CompletionTime),
		Event: &kafkaevent.WalletEvent_Unstake{
			Unstake: &kafkaevent.UnstakeEvent{
				Coin: &kafkaevent.Coin{
					Denom:  "",
					Amount: entry.Balance,
				},
				CompletionTime: timestamppb.New(entry.CompletionTime),
			},
		},
	}
	pbEvent, err := proto.Marshal(txEvent)
	if err != nil {
		return nil, err
	}
	return pbEvent, nil
}

func (c *SetupCrawler) fetchUnstakingEvents(els []*ent.EventListener) {
	log.Sugar.Debug("Fetch unstake events")
	var pbEvents [][]byte
	for _, el := range els {
		if el.DataType != eventlistener.DataTypeWalletEvent_Unstake {
			continue
		}
		if el.WalletAddress == "" {
			continue
		}
		chain, err := el.QueryChain().
			Only(context.Background())
		if err != nil {
			log.Sugar.Errorf("while fetching for event listener %v: %v", el.ID, err)
			continue
		}
		if strings.Contains(chain.Path, "neutron") {
			continue
		}

		url := fmt.Sprintf(urlUnstaking, chain.RestEndpoint, el.WalletAddress)
		var allocation types.UnstakingResponse
		_, err = common.GetJson(url, 5, &allocation)
		if err != nil {
			log.Sugar.Errorf("while fetching unstaking for %v on %v: %v", el.WalletAddress, chain.PrettyName, err)
		}
		for _, delegation := range allocation.UnbondingResponses {
			for _, entry := range delegation.Entries {
				pbEvent, err := c.createEvent(chain, el.WalletAddress, entry)
				if err != nil {
					log.Sugar.Errorf("while creating event for %v on %v: %v", el.WalletAddress, chain.PrettyName, err)
					continue
				}
				pbEvents = append(pbEvents, pbEvent)
			}
		}
	}
	if len(pbEvents) > 0 {
		log.Sugar.Debugf("Sending %v unstaking events", len(pbEvents))
		c.kafkaInternal.ProduceWalletEvents(pbEvents)
	}
}
