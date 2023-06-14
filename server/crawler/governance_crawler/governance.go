package governance_crawler

import (
	"context"
	"fmt"
	cosmossdktypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
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
const urlVote = urlProposals + "/%v/votes/%v"

type GovernanceCrawler struct {
	chainManager         *database.ChainManager
	eventListenerManager *database.EventListenerManager
	kafkaInternal        kafka_internal.KafkaInternal
}

func NewGovernanceCrawler(dbManagers *database.DbManagers, kafkaInternal kafka_internal.KafkaInternal) *GovernanceCrawler {
	return &GovernanceCrawler{
		chainManager:         dbManagers.ChainManager,
		eventListenerManager: dbManagers.EventListenerManager,
		kafkaInternal:        kafkaInternal,
	}
}

func (c *GovernanceCrawler) createGovEvent(chain *ent.Chain, prop *types.Proposal) ([]byte, error) {
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
			if chain.Path == "neutron" {
				NewNeutronCrawler(c.chainManager, c.kafkaInternal).
					fetchProposals(chain, "neutron1suhgf5svhu4usrurvxzlgn54ksxmn8gljarjtxqnapv8kjnp4nrstdxvff")
			} else if chain.Path == "neutron-testnet" || chain.Path == "neutron-pion" {
				NewNeutronCrawler(c.chainManager, c.kafkaInternal).
					fetchProposals(chain, "neutron1suhgf5svhu4usrurvxzlgn54ksxmn8gljarjtxqnapv8kjnp4nrstdxvff")
			}
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
						pbEvent, err := c.createGovEvent(chain, &prop)
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
				pbEvent, err := c.createGovEvent(chain, &prop)
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
		c.kafkaInternal.ProduceChainEvents(pbEvents)
	}
}

func (c *GovernanceCrawler) createVoteReminderEvent(chain *ent.Chain, prop *ent.Proposal, walletAddress string) ([]byte, error) {
	var now = timestamppb.Now()
	var chainEvent = &kafkaevent.WalletEvent{
		ChainId:       uint64(chain.ID),
		WalletAddress: walletAddress,
		Timestamp:     now,
		NotifyTime:    now,
		Event: &kafkaevent.WalletEvent_VoteReminder{
			VoteReminder: &kafkaevent.VoteReminderEvent{
				ProposalId:  prop.ProposalID,
				VoteEndTime: timestamppb.New(prop.VotingEndTime),
			},
		},
	}
	pbEvent, err := proto.Marshal(chainEvent)
	if err != nil {
		return nil, err
	}
	return pbEvent, nil
}

func (c *GovernanceCrawler) createVotedEvent(chain *ent.Chain, proposalId uint64, option string, walletAddress string) ([]byte, error) {
	var now = timestamppb.Now()
	var chainEvent = &kafkaevent.WalletEvent{
		ChainId:       uint64(chain.ID),
		WalletAddress: walletAddress,
		Timestamp:     now,
		NotifyTime:    now,
		Event: &kafkaevent.WalletEvent_Voted{
			Voted: &kafkaevent.VotedEvent{
				ProposalId: proposalId,
				Option:     option,
			},
		},
	}
	pbEvent, err := proto.Marshal(chainEvent)
	if err != nil {
		return nil, err
	}
	return pbEvent, nil
}

func (c *GovernanceCrawler) fetchVotingReminders() {
	log.Sugar.Info("Checking for voting reminders")
	voteReminders, err := c.eventListenerManager.QueryForVoteReminderAddresses(context.Background())
	if err != nil {
		log.Sugar.Errorf("while querying for vote reminders: %v", err)
		return
	}
	pbEvents := make([][]byte, 0)
	for _, vr := range voteReminders {
		var chain = vr.Chain
		var walletAddress = vr.EventListener.WalletAddress
		log.Sugar.Debugf("Check vote reminder for proposal %v on %v for wallet %v", vr.Proposal.ProposalID, chain.PrettyName, walletAddress)
		var proposalId = vr.Proposal.ProposalID
		url := fmt.Sprintf(urlVote, chain.RestEndpoint, proposalId, walletAddress)
		var voteResponse types.ChainProposalVoteResponse
		statusCode, err := common.GetJson(url, 5, &voteResponse)
		if err != nil && statusCode == 400 {
			pbEvent, err := c.createVoteReminderEvent(chain, vr.Proposal, walletAddress)
			if err != nil {
				log.Sugar.Errorf("while creating vote reminder event for proposal %v: %v", proposalId, err)
				continue
			}
			pbEvents = append(pbEvents, pbEvent)
		} else if err != nil {
			log.Sugar.Errorf("while fetching vote for proposal %v: %v", proposalId, err)
			continue
		} else {
			if voteResponse.Vote.Option.ToCosmosType() == cosmossdktypes.OptionEmpty {
				pbEvent, err := c.createVoteReminderEvent(chain, vr.Proposal, walletAddress)
				if err != nil {
					log.Sugar.Errorf("while creating vote reminder event for proposal %v: %v", proposalId, err)
					continue
				}
				pbEvents = append(pbEvents, pbEvent)

			} else {
				pbEvent, err := c.createVotedEvent(chain, proposalId, voteResponse.Vote.Option.ToCosmosType().String(), walletAddress)
				if err != nil {
					log.Sugar.Errorf("while creating voted event for proposal %v: %v", proposalId, err)
					continue
				}
				pbEvents = append(pbEvents, pbEvent)
			}
		}
	}
	if len(pbEvents) > 0 {
		log.Sugar.Debugf("Sending %v governance events", len(pbEvents))
		c.kafkaInternal.ProduceWalletEvents(pbEvents)
	}
}

func (c *GovernanceCrawler) StartCrawling() {
	c.fetchProposals()
	c.fetchVotingReminders()
	log.Sugar.Info("Scheduling governance crawl")
	cr := cron.New()
	_, err := cr.AddFunc("*/15 * * * *", func() { c.fetchProposals() }) // every 10min
	if err != nil {
		log.Sugar.Errorf("while executing 'fetchProposals' via cron: %v", err)
	}
	_, err = cr.AddFunc("5-50/15 * * * *", func() { c.fetchVotingReminders() }) // every 15min but not at the same time as fetchProposals
	if err != nil {
		log.Sugar.Errorf("while executing 'fetchProposals' via cron: %v", err)
	}
	cr.Start()
	select {}
}
