package indexer

import (
	"buf.build/gen/go/loomi-labs/star-scope/bufbuild/connect-go/grpc/indexer/indexerpb/indexerpbconnect"
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/grpc/indexer/indexerpb"
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/indexevent"
	"context"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/golang/protobuf/proto"
	"github.com/loomi-labs/star-scope/indexers/base/common"
	"github.com/loomi-labs/star-scope/indexers/base/kafka"
	"github.com/loomi-labs/star-scope/indexers/base/types"
	"github.com/robfig/cron/v3"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

const urlUnstaking = "%v/cosmos/staking/v1beta1/delegators/%v/unbonding_delegations"

type SetupCrawler struct {
	grpcClient       indexerpbconnect.IndexerServiceClient
	kafkaProducer    *kafka.KafkaProducer
	processedWallets map[string]bool
}

func NewSetupCrawler(grpcClient indexerpbconnect.IndexerServiceClient, kafkaBrokers []string) *SetupCrawler {
	return &SetupCrawler{
		grpcClient:       grpcClient,
		kafkaProducer:    kafka.NewKafkaProducer(kafka.IndexEventsTopic, kafkaBrokers...),
		processedWallets: make(map[string]bool),
	}
}

func (c *SetupCrawler) createEvent(chain *indexerpb.NewAccountsChainInfo, account string, entry types.UnstakingEntry) ([]byte, error) {
	var now = timestamppb.Now()
	txEvent := &indexevent.TxEvent{
		ChainId:       chain.Id,
		WalletAddress: account,
		Timestamp:     now,
		NotifyTime:    timestamppb.New(entry.CompletionTime),
		Event: &indexevent.TxEvent_Unstake{
			Unstake: &indexevent.UnstakeEvent{
				Coin: &indexevent.Coin{
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

func (c *SetupCrawler) fetchUnstakings() {
	log.Sugar.Debug("Fetch unstakings")
	stati, err := c.grpcClient.GetNewAccounts(
		context.Background(),
		connect.NewRequest(&indexerpb.GetNewAccountsRequest{}),
	)
	if err != nil {
		log.Sugar.Errorf("Error getting indexing chains: %v", err)
	}
	var pbEvents [][]byte
	for _, chain := range stati.Msg.GetChains() {
		if strings.Contains(chain.Path, "neutron") {
			continue
		}

		for _, account := range chain.GetNewAccounts() {
			if _, ok := c.processedWallets[account]; ok {
				continue
			}

			url := fmt.Sprintf(urlUnstaking, chain.RestEndpoint, account)
			var allocation types.UnstakingResponse
			_, err := common.GetJson(url, 5, &allocation)
			if err != nil {
				log.Sugar.Errorf("while fetching unstaking for %v on %v: %v", account, chain.Name, err)
			}
			for _, delegation := range allocation.UnbondingResponses {
				for _, entry := range delegation.Entries {
					event, err := c.createEvent(chain, account, entry)
					if err != nil {
						log.Sugar.Errorf("while creating event for %v on %v: %v", account, chain.Name, err)
						continue
					}
					pbEvents = append(pbEvents, event)
				}
			}
			c.processedWallets[account] = true
		}
	}
	if len(pbEvents) > 0 {
		log.Sugar.Debugf("Sending %v governance events", len(pbEvents))
		c.kafkaProducer.Produce(pbEvents)
	}
}

func (c *SetupCrawler) StartCrawling() {
	c.fetchUnstakings()
	log.Sugar.Info("Scheduling setup crawl")
	cr := cron.New()
	_, err := cr.AddFunc("*/10 * * * *", func() { c.fetchUnstakings() }) // every 10min
	if err != nil {
		log.Sugar.Errorf("while executing 'fetchProposals' via cron: %v", err)
	}
	cr.Start()
}
