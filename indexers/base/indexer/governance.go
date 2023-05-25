package indexer

import (
	"buf.build/gen/go/loomi-labs/star-scope/bufbuild/connect-go/grpc/indexer/indexerpb/indexerpbconnect"
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/grpc/indexer/indexerpb"
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/queryevent"
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

const urlProposals = "%v/cosmos/gov/v1beta1/proposals"

type GovernanceCrawler struct {
	grpcClient    indexerpbconnect.IndexerServiceClient
	kafkaProducer *kafka.KafkaProducer
}

func NewGovernanceCrawler(grpcClient indexerpbconnect.IndexerServiceClient, kafkaBrokers []string) *GovernanceCrawler {
	return &GovernanceCrawler{
		grpcClient:    grpcClient,
		kafkaProducer: kafka.NewKafkaProducer(kafka.QueryEventsTopic, kafkaBrokers...),
	}
}

func (c *GovernanceCrawler) createEvent(chain *indexerpb.ChainInfo, prop types.Proposal) ([]byte, error) {
	var now = timestamppb.Now()
	event := &queryevent.QueryEvent{
		ChainId:    chain.Id,
		Timestamp:  now,
		NotifyTime: now,
		Event: &queryevent.QueryEvent_GovernanceProposal{
			GovernanceProposal: &queryevent.GovernanceProposalEvent{
				ProposalId:      uint64(prop.ProposalId),
				Title:           prop.Content.Title,
				Description:     prop.Content.Description,
				ProposalType:    queryevent.ProposalType_PROPOSAL_TYPE_UNSPECIFIED,
				ProposalStatus:  queryevent.ProposalStatus(prop.Status),
				VotingStartTime: timestamppb.New(prop.VotingStartTime),
				VotingEndTime:   timestamppb.New(prop.VotingEndTime),
			},
		},
	}
	pbEvent, err := proto.Marshal(event)
	if err != nil {
		return nil, err
	}
	return pbEvent, nil
}

func (c *GovernanceCrawler) fetchProposals() {
	log.Sugar.Debug("Fetching proposals")
	stati, err := c.grpcClient.GetGovernanceProposalStati(
		context.Background(),
		connect.NewRequest(&indexerpb.GetGovernanceProposalStatiRequest{}),
	)
	if err != nil {
		log.Sugar.Errorf("Error getting indexing chains: %v", err)
	}
	var pbEvents [][]byte
	for _, chain := range stati.Msg.GetChains() {
		if strings.Contains(chain.Path, "neutron") {
			continue
		}

		url := fmt.Sprintf(urlProposals+"?pagination.reverse=true&limit=100", chain.RpcUrl)

		var resp types.ProposalsResponse
		_, err := common.GetJson(url, 5, &resp)
		if err != nil {
			log.Sugar.Errorf("while fetching proposals for chain %v: %v", chain.Name, err)
			continue
		}
		for _, prop := range resp.Proposals {
			var found = false
			for _, currentProp := range chain.GetProposals() {
				if uint64(prop.ProposalId) == currentProp.GetProposalId() {
					if prop.Status != types.ProposalStatus(currentProp.GetStatus().Number()) {
						log.Sugar.Debugf("Proposal %v changed status from %v to %v", prop.ProposalId, currentProp.GetStatus().Number(), prop.Status)
						pbEvent, err := c.createEvent(chain, prop)
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
				pbEvent, err := c.createEvent(chain, prop)
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
		c.kafkaProducer.Produce(pbEvents)
	}
}

func (c *GovernanceCrawler) StartGovernanceCrawling() {
	c.fetchProposals()
	log.Sugar.Info("Scheduling governance crawl")
	cr := cron.New()
	_, err := cr.AddFunc("*/10 * * * *", func() { c.fetchProposals() }) // every 10min
	if err != nil {
		log.Sugar.Errorf("while executing 'fetchProposals' via cron: %v", err)
	}
	cr.Start()
}
