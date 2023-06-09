package governance_crawler

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
)

const urlProposals = "%v/cosmos/gov/v1beta1/proposals"

type GovernanceCrawler struct {
	chainManager  *database.ChainManager
	kafkaInternal *kafka_internal.KafkaInternal
}

func NewGovernanceCrawler(dbManagers *database.DbManagers, kafkaBrokers []string) *GovernanceCrawler {
	return &GovernanceCrawler{
		chainManager:  dbManagers.ChainManager,
		kafkaInternal: kafka_internal.NewKafkaInternal(kafkaBrokers),
	}
}

func (c *GovernanceCrawler) createEvent(chain *ent.Chain, prop *types.Proposal) ([]byte, error) {
	var now = timestamppb.Now()
	chainEvent := &kafkaevent.ChainEvent{
		ChainId:    uint64(chain.ID),
		Timestamp:  now,
		NotifyTime: now,
		Event: &kafkaevent.ChainEvent_GovernanceProposal{
			GovernanceProposal: &kafkaevent.GovernanceProposalEvent{
				ProposalId:      uint64(prop.ProposalId),
				Title:           prop.Content.Title,
				Description:     prop.Content.Description,
				ProposalType:    kafkaevent.ProposalType_PROPOSAL_TYPE_UNSPECIFIED,
				ProposalStatus:  kafkaevent.ProposalStatus(prop.Status),
				VotingStartTime: timestamppb.New(prop.VotingStartTime),
				VotingEndTime:   timestamppb.New(prop.VotingEndTime),
			},
		},
	}
	pbEvent, err := proto.Marshal(chainEvent)
	if err != nil {
		return nil, err
	}
	return pbEvent, nil
}

func (c *GovernanceCrawler) fetchProposals() {
	log.Sugar.Debug("Fetching governance proposals")

	var chains = c.chainManager.QueryEnabledWithProposals(context.Background())

	var pbEvents [][]byte
	for _, chain := range chains {
		if strings.Contains(chain.Path, "neutron") {
			continue
		}

		url := fmt.Sprintf(urlProposals+"?pagination.reverse=true&limit=100", chain.RestEndpoint)

		var resp types.ProposalsResponse
		_, err := common.GetJson(url, 5, &resp)
		if err != nil {
			log.Sugar.Errorf("while fetching proposals for chain %v: %v", chain.Name, err)
			continue
		}
		for _, prop := range resp.Proposals {
			var found = false
			for _, currentProp := range chain.Edges.Proposals {
				if uint64(prop.ProposalId) == currentProp.ProposalID {
					if prop.Status.String() != currentProp.Status.String() {
						log.Sugar.Debugf("Proposal %v changed status from %v to %v", prop.ProposalId, currentProp.Status.String(), prop.Status.String())
						_, err := c.chainManager.CreateOrUpdateProposal(context.Background(), chain, &prop)
						if err != nil {
							log.Sugar.Errorf("while updating proposal %v: %v", prop.ProposalId, err)
							break
						}
						pbEvent, err := c.createEvent(chain, &prop)
						if err != nil {
							log.Sugar.Errorf("while creating event for proposal %v: %v", prop.ProposalId, err)
							break
						}
						pbEvents = append(pbEvents, pbEvent)
					}
					found = true
					break
				}
			}
			if !found {
				log.Sugar.Debugf("New proposal on %v #%v", chain.Name, prop.ProposalId)
				_, err := c.chainManager.CreateOrUpdateProposal(context.Background(), chain, &prop)
				if err != nil {
					log.Sugar.Errorf("while creating proposal %v: %v", prop.ProposalId, err)
					break
				}
				pbEvent, err := c.createEvent(chain, &prop)
				if err != nil {
					log.Sugar.Errorf("while creating event for proposal %v: %v", prop.ProposalId, err)
					continue
				}
				pbEvents = append(pbEvents, pbEvent)
				continue
			}
		}
	}
	if len(pbEvents) > 0 {
		log.Sugar.Debugf("Sending %v governance events", len(pbEvents))
		c.kafkaInternal.ProduceEvents(pbEvents)
	}
}

func (c *GovernanceCrawler) StartCrawling() {
	c.fetchProposals()
	log.Sugar.Info("Scheduling governance crawl")
	cr := cron.New()
	_, err := cr.AddFunc("*/10 * * * *", func() { c.fetchProposals() }) // every 10min
	if err != nil {
		log.Sugar.Errorf("while executing 'fetchProposals' via cron: %v", err)
	}
	cr.Start()
	select {}
}
