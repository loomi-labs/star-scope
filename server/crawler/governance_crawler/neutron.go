package governance_crawler

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/ent"
	kafkaevent "github.com/loomi-labs/star-scope/event"
	"github.com/loomi-labs/star-scope/kafka_internal"
	"github.com/loomi-labs/star-scope/types"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

const urlCosmWasm = "%v/cosmwasm/wasm/v1/contract/%v/smart/%v"

type NeutronCrawler struct {
	chainManager  *database.ChainManager
	kafkaInternal kafka_internal.KafkaInternal
}

func NewNeutronCrawler(chainManager *database.ChainManager, kafkaInternal kafka_internal.KafkaInternal) *NeutronCrawler {
	return &NeutronCrawler{
		chainManager:  chainManager,
		kafkaInternal: kafkaInternal,
	}
}

func (c *NeutronCrawler) createGovPropEvent(chain *ent.Chain, contractAddress string, id int, prop types.ContractProposal) ([]byte, error) {
	var now = timestamppb.Now()
	contractEvent := &kafkaevent.ContractEvent{
		ChainId:    uint64(chain.ID),
		Timestamp:  now,
		NotifyTime: now,
		Event: &kafkaevent.ContractEvent_ContractGovernanceProposal{
			ContractGovernanceProposal: &kafkaevent.ContractGovernanceProposalEvent{
				ProposalId:      uint64(id),
				Title:           prop.Title,
				Description:     prop.Description,
				ProposalType:    kafkaevent.ProposalType_PROPOSAL_TYPE_UNSPECIFIED.String(),
				ProposalStatus:  kafkaevent.ContractProposalStatus(prop.Status),
				ContractAddress: contractAddress,
				FirstSeenTime:   now,
				VotingEndTime:   timestamppb.New(time.Time(prop.Expiration.AtTime)),
			},
		},
	}
	pbEvent, err := proto.Marshal(contractEvent)
	if err != nil {
		return nil, err
	}
	return pbEvent, nil
}

func (c *NeutronCrawler) createGovPropEvents(chain *ent.Chain, contractAddress string, propResp types.ListProposalsResponse) ([][]byte, error) {
	var pbEvents [][]byte
	var contractProposals = chain.QueryContractProposals().AllX(context.Background())
	for _, prop := range propResp.Data.Proposals {
		var found = false
		for _, currentProp := range contractProposals {
			if currentProp.ProposalID == uint64(prop.ID) && currentProp.ContractAddress == contractAddress {
				found = true
				if currentProp.Status.String() != prop.Proposal.Status.String() {
					_, err := c.chainManager.CreateOrUpdateContractProposal(context.Background(), chain, uint64(prop.ID), contractAddress, &prop.Proposal)
					if err != nil {
						log.Sugar.Errorf("while updating proposal %v of %v: %v", prop.ID, contractAddress, err)
						break
					}
					pbEvent, err := c.createGovPropEvent(chain, contractAddress, prop.ID, prop.Proposal)
					if err != nil {
						log.Sugar.Errorf("while creating event for proposal %v of %v: %v", prop.ID, contractAddress, err)
						break
					}
					pbEvents = append(pbEvents, pbEvent)
				}
				break
			}
		}
		if !found {
			log.Sugar.Debugf("New proposal on %v #%v", contractAddress, prop.ID)
			_, err := c.chainManager.CreateOrUpdateContractProposal(context.Background(), chain, uint64(prop.ID), contractAddress, &prop.Proposal)
			if err != nil {
				log.Sugar.Errorf("while updating proposal %v of %v: %v", prop.ID, contractAddress, err)
				break
			}
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

func (c *NeutronCrawler) fetchProposals(chain *ent.Chain, mainContract string) {
	log.Sugar.Debug("Fetching proposals for chain ", chain.PrettyName)

	payload := base64.StdEncoding.EncodeToString([]byte(`{"config":{}}`))
	url := fmt.Sprintf(urlCosmWasm, chain.RestEndpoint, mainContract, payload)
	var config types.ConfigResponse
	_, err := common.GetJson(url, 5, &config)
	if err != nil {
		log.Sugar.Errorf("while fetching config for chain %v: %v", chain.Name, err)
	}

	payload = base64.StdEncoding.EncodeToString([]byte(`{"proposal_modules":{}}`))
	url = fmt.Sprintf(urlCosmWasm, chain.RestEndpoint, mainContract, payload)
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
		c.kafkaInternal.ProduceContractEvents(pbEvents)
	}
}
