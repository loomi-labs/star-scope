package neutron

import (
	"buf.build/gen/go/loomi-labs/star-scope/bufbuild/connect-go/grpc/indexer/indexerpb/indexerpbconnect"
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/grpc/indexer/indexerpb"
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/indexevent"
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
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"
)

const urlCosmWasm = "%v/cosmwasm/wasm/v1/contract/%v/smart/%v"

type NeutronCrawler struct {
	grpcClient           indexerpbconnect.IndexerServiceClient
	chainPath            string
	mainContract         string
	creditsVaultContract string
	kafkaQueryProducer   *kafka.KafkaProducer
	kafkaIndexProducer   *kafka.KafkaProducer
}

func NewNeutronCrawler(grpcClient indexerpbconnect.IndexerServiceClient, chainPath string, mainContract string, creditsVaultContract string, kafkaBrokers []string) *NeutronCrawler {
	return &NeutronCrawler{
		grpcClient:           grpcClient,
		chainPath:            chainPath,
		mainContract:         mainContract,
		creditsVaultContract: creditsVaultContract,
		kafkaQueryProducer:   kafka.NewKafkaProducer(kafka.QueryEventsTopic, kafkaBrokers...),
		kafkaIndexProducer:   kafka.NewKafkaProducer(kafka.IndexEventsTopic, kafkaBrokers...),
	}
}

func (c *NeutronCrawler) createGovPropEvent(chain *indexerpb.GovernanceChainInfo, contractAddress string, id int, prop types.ContractProposal) ([]byte, error) {
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

func (c *NeutronCrawler) createGovPropEvents(chain *indexerpb.GovernanceChainInfo, contractAddress string, propResp types.ListProposalsResponse) ([][]byte, error) {
	var pbEvents [][]byte
	for _, prop := range propResp.Data.Proposals {
		var found = false
		for _, currentProp := range chain.GetContractProposals() {
			if currentProp.GetProposalId() == uint64(prop.ID) && currentProp.GetContractAddress() == contractAddress {
				if currentProp.GetStatus() != queryevent.ContractProposalStatus(prop.Proposal.Status) {
					pbEvent, err := c.createGovPropEvent(chain, contractAddress, prop.ID, prop.Proposal)
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
			pbEvent, err := c.createGovPropEvent(chain, contractAddress, prop.ID, prop.Proposal)
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
		ChainPaths: []string{c.chainPath},
	})
	response, err := c.grpcClient.GetGovernanceProposalStati(context.Background(), request)
	if err != nil {
		log.Sugar.Errorf("while fetching proposals: %v", err)
		return
	}

	for _, chain := range response.Msg.GetChains() {
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
			newPbEvents, err := c.createGovPropEvents(chain, module.Address, proposals)
			if err != nil {
				log.Sugar.Errorf("while creating events for chain %v: %v", chain.Name, err)
				continue
			}
			pbEvents = append(pbEvents, newPbEvents...)
		}

		if len(pbEvents) > 0 {
			log.Sugar.Debugf("Sending %v neutron governance events", len(pbEvents))
			c.kafkaQueryProducer.Produce(pbEvents)
		}
	}
}

func (c *NeutronCrawler) createVestingEvent(chain *indexerpb.NewAccountsChainInfo, contractAddress string, propResp types.CreditsVaultAllocationResponse, account string) ([]byte, error) {
	var now = timestamppb.Now()
	amount, err := strconv.ParseUint(propResp.Data.AllocatedAmount, 10, 64)
	if err != nil {
		return nil, err
	}
	duration := durationpb.New(time.Duration(propResp.Data.Schedule.Duration))
	unlockTime := timestamppb.New(time.Unix(propResp.Data.Schedule.StartTime, 0))

	var event = &indexevent.TxEvent{
		ChainId:       chain.Id,
		WalletAddress: account,
		Timestamp:     now,
		NotifyTime:    unlockTime,
		Event: &indexevent.TxEvent_NeutronTokenVesting{
			NeutronTokenVesting: &indexevent.NeutronTokenVestingEvent{
				WalletAddress: account,
				Amount:        amount,
				Duration:      duration,
				UnlockTime:    unlockTime,
			},
		},
	}
	pbEvent, err := proto.Marshal(event)
	if err != nil {
		return nil, err
	}
	return pbEvent, nil
}

func (c *NeutronCrawler) processNewAccounts() {
	if c.creditsVaultContract == "" {
		log.Sugar.Debug("No credits vault contract set")
		return
	}
	log.Sugar.Debug("Processing new accounts")

	var request = connect.NewRequest(&indexerpb.GetNewAccountsRequest{
		ChainPaths: []string{c.chainPath},
	})
	response, err := c.grpcClient.GetNewAccounts(context.Background(), request)
	if err != nil {
		log.Sugar.Errorf("while fetching proposals: %v", err)
		return
	}

	var pbEvents [][]byte
	for _, chain := range response.Msg.GetChains() {
		for _, account := range chain.GetNewAccounts() {
			payload := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`{"allocation":{"address":"%v"}}`, account)))
			url := fmt.Sprintf(urlCosmWasm, chain.RestEndpoint, c.creditsVaultContract, payload)
			var allocation types.CreditsVaultAllocationResponse
			_, err := common.GetJson(url, 5, &allocation)
			if err != nil {
				log.Sugar.Errorf("while fetching allocation for %v on %v: %v", account, chain.Name, err)
			}
			newPbEvent, err := c.createVestingEvent(chain, c.creditsVaultContract, allocation, account)
			if err != nil {
				log.Sugar.Errorf("while creating events for chain %v: %v", chain.Name, err)
				continue
			}
			pbEvents = append(pbEvents, newPbEvent)
		}
	}

	if len(pbEvents) > 0 {
		log.Sugar.Debugf("Sending %v neutron vesting events", len(pbEvents))
		c.kafkaIndexProducer.Produce(pbEvents)
	}
}

func (c *NeutronCrawler) StartCrawling() {
	c.fetchProposals()
	c.processNewAccounts()
	log.Sugar.Info("Scheduling governance crawl")
	cr := cron.New()
	_, err := cr.AddFunc("*/10 * * * *", func() { c.fetchProposals() }) // every 10 min
	if err != nil {
		log.Sugar.Errorf("while executing 'fetchProposals' via cron: %v", err)
	}
	_, err = cr.AddFunc("*/10 * * * *", func() { c.processNewAccounts() }) // every 10 min
	if err != nil {
		log.Sugar.Errorf("while executing 'processNewAccounts' via cron: %v", err)
	}
	cr.Start()
}
