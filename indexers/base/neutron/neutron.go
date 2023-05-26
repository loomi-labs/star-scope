package neutron

import (
	"buf.build/gen/go/loomi-labs/star-scope/bufbuild/connect-go/grpc/indexer/indexerpb/indexerpbconnect"
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/grpc/indexer/indexerpb"
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/queryevent"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/golang/protobuf/proto"
	"github.com/loomi-labs/star-scope/indexers/base/common"
	"github.com/loomi-labs/star-scope/indexers/base/kafka"
	"github.com/loomi-labs/star-scope/indexers/base/types"
	"github.com/robfig/cron/v3"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

const urlCosmWasm = "%v/cosmwasm/wasm/v1/contract/%v/smart/%v"

type NeutronCrawler struct {
	grpcClient    indexerpbconnect.IndexerServiceClient
	mainContract  string
	kafkaProducer *kafka.KafkaProducer
}

func NewNeutronCrawler(grpcClient indexerpbconnect.IndexerServiceClient, mainContract string, kafkaBrokers []string) *NeutronCrawler {
	return &NeutronCrawler{
		grpcClient:    grpcClient,
		mainContract:  mainContract,
		kafkaProducer: kafka.NewKafkaProducer(kafka.QueryEventsTopic, kafkaBrokers...),
	}
}

func (c *NeutronCrawler) createEvent(chain *indexerpb.ChainInfo, contractAddress string, id int, prop types.ContractProposal) ([]byte, error) {
	var now = timestamppb.Now()
	event := &queryevent.QueryEvent{
		ChainId:    chain.Id,
		Timestamp:  now,
		NotifyTime: now,
		Event: &queryevent.QueryEvent_ContractGovernanceProposal{
			ContractGovernanceProposal: &queryevent.ContractGovernanceProposalEvent{
				ProposalId:      uint64(id),
				Title:           prop.Title,
				Description:     prop.Description,
				ProposalType:    queryevent.ProposalType_PROPOSAL_TYPE_UNSPECIFIED,
				ProposalStatus:  queryevent.ContractProposalStatus(prop.Status),
				ContractAddress: contractAddress,
				FirstSeenTime:   now,
				VotingEndTime:   timestamppb.New(time.Time(prop.Expiration.AtTime)),
			},
		},
	}
	pbEvent, err := proto.Marshal(event)
	if err != nil {
		return nil, err
	}
	return pbEvent, nil
}

func (c *NeutronCrawler) createEvents(chain *indexerpb.ChainInfo, contractAddress string, propResp types.ListProposalsResponse) ([][]byte, error) {
	var pbEvents [][]byte
	for _, prop := range propResp.Data.Proposals {
		var found = false
		for _, currentProp := range chain.GetContractProposals() {
			if currentProp.GetProposalId() == uint64(prop.ID) && currentProp.GetContractAddress() == contractAddress {
				if currentProp.GetStatus() != queryevent.ContractProposalStatus(prop.Proposal.Status) {
					pbEvent, err := c.createEvent(chain, contractAddress, prop.ID, prop.Proposal)
					if err != nil {
						log.Sugar.Errorf("while creating event for proposal %v of %v: %v", prop.ID, contractAddress, err)
						break
					}
					pbEvents = append(pbEvents, pbEvent)
				}
				found = true
				break
			}
		}
		if !found {
			log.Sugar.Debugf("New proposal on %v #%v", contractAddress, prop.ID)
			pbEvent, err := c.createEvent(chain, contractAddress, prop.ID, prop.Proposal)
			if err != nil {
				log.Sugar.Errorf("while creating event for proposal %v of %v: %v", prop.ID, contractAddress, err)
				continue
			}
			pbEvents = append(pbEvents, pbEvent)
			continue
		}

	}
	return pbEvents, nil
}

func (c *NeutronCrawler) fetchProposals() {
	log.Sugar.Debug("Fetching proposals")

	var request = connect.NewRequest(&indexerpb.GetGovernanceProposalStatiRequest{
		ChainPaths: []string{"neutron", "neutron-pion"},
	})
	stati, err := c.grpcClient.GetGovernanceProposalStati(context.Background(), request)
	if err != nil {
		log.Sugar.Errorf("while fetching proposals: %v", err)
		return
	}

	for _, chain := range stati.Msg.GetChains() {
		payload := base64.StdEncoding.EncodeToString([]byte(`{"config":{}}`))
		url := fmt.Sprintf(urlCosmWasm, chain.RestEndpoint, c.mainContract, payload)
		var config types.ConfigResponse
		_, err := common.GetJson(url, 5, &config)
		if err != nil {
			log.Sugar.Errorf("while fetching config for chain %v: %v", chain.Name, err)
		}

		payload = base64.StdEncoding.EncodeToString([]byte(`{"proposal_modules":{}}`))
		url = fmt.Sprintf(urlCosmWasm, chain.RestEndpoint, c.mainContract, payload)
		var proposalModules types.ProposalModulesResponse
		_, err = common.GetJson(url, 5, &proposalModules)
		if err != nil {
			log.Sugar.Errorf("while fetching proposal modules for chain %v: %v", chain.Name, err)
		}

		var pbEvents [][]byte
		for _, module := range proposalModules.Data {
			payload = base64.StdEncoding.EncodeToString([]byte(`{"reverse_proposals":{}}`))
			url = fmt.Sprintf(urlCosmWasm, chain.RestEndpoint, module.Address, payload)
			var proposals types.ListProposalsResponse
			_, err = common.GetJson(url, 5, &proposals)
			if err != nil {
				log.Sugar.Errorf("while fetching proposal modules for chain %v: %v", chain.Name, err)
				continue
			}
			newPbEvents, err := c.createEvents(chain, module.Address, proposals)
			if err != nil {
				log.Sugar.Errorf("while creating events for chain %v: %v", chain.Name, err)
				continue
			}
			pbEvents = append(pbEvents, newPbEvents...)
		}

		if len(pbEvents) > 0 {
			log.Sugar.Debugf("Sending %v governance events", len(pbEvents))
			c.kafkaProducer.Produce(pbEvents)
		}
	}
}

func (c *NeutronCrawler) StartCrawling() {
	c.fetchProposals()
	log.Sugar.Info("Scheduling governance crawl")
	cr := cron.New()
	_, err := cr.AddFunc("*/10 * * * *", func() { c.fetchProposals() }) // every 10min
	if err != nil {
		log.Sugar.Errorf("while executing 'fetchProposals' via cron: %v", err)
	}
	cr.Start()
}
