package kafka_internal

import (
	"errors"
	"fmt"
	"github.com/loomi-labs/star-scope/ent"
	kafkaevent "github.com/loomi-labs/star-scope/event"
	"github.com/loomi-labs/star-scope/grpc/event/eventpb"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProto(entEvent *ent.Event) (*eventpb.Event, error) {
	if entEvent.ChainEvent.ChainEvent != nil {
		var chainEvent = entEvent.ChainEvent.ChainEvent
		switch chainEvent.GetEvent().(type) {
		case *kafkaevent.ChainEvent_GovernanceProposal:
			var statusText = "Proposal %v"
			switch chainEvent.GetGovernanceProposal().GetProposalStatus() {
			case kafkaevent.ProposalStatus_PROPOSAL_STATUS_VOTING_PERIOD:
				statusText = "New Proposal - %v"
			case kafkaevent.ProposalStatus_PROPOSAL_STATUS_PASSED:
				statusText = "Proposal %v Passed"
			case kafkaevent.ProposalStatus_PROPOSAL_STATUS_REJECTED:
				statusText = "Proposal %v Rejected"
			case kafkaevent.ProposalStatus_PROPOSAL_STATUS_FAILED:
				statusText = "Proposal %v Failed"
			default:
				log.Sugar.Errorf("Unknown proposal status %v", chainEvent.GetGovernanceProposal().GetProposalStatus())
			}
			return &eventpb.Event{
				Title:       fmt.Sprintf(statusText, chainEvent.GetGovernanceProposal().GetProposalId()),
				Subtitle:    chainEvent.GetGovernanceProposal().GetTitle(),
				Description: chainEvent.GetGovernanceProposal().GetDescription(),
				CreatedAt:   chainEvent.Timestamp,
				EventType:   kafkaevent.EventType_GOVERNANCE,
			}, nil
		default:
			return nil, errors.New(fmt.Sprintf("No type defined for event %v", chainEvent.GetEvent()))
		}
	} else if entEvent.ContractEvent.ContractEvent != nil {
		var contractEvent = entEvent.ContractEvent.ContractEvent
		switch contractEvent.GetEvent().(type) {
		case *kafkaevent.ContractEvent_ContractGovernanceProposal:
			var statusText = "Proposal %v"
			switch contractEvent.GetContractGovernanceProposal().GetProposalStatus() {
			case kafkaevent.ContractProposalStatus_OPEN:
				statusText = "New Proposal - %v"
			case kafkaevent.ContractProposalStatus_PASSED:
				statusText = "Proposal %v Passed"
			case kafkaevent.ContractProposalStatus_REJECTED:
				statusText = "Proposal %v Rejected"
			case kafkaevent.ContractProposalStatus_EXECUTION_FAILED:
				statusText = "Proposal %v Failed"
			case kafkaevent.ContractProposalStatus_CLOSED:
				statusText = "Proposal %v Closed"
			}
			return &eventpb.Event{
				Title:       fmt.Sprintf(statusText, contractEvent.GetContractGovernanceProposal().GetProposalId()),
				Subtitle:    contractEvent.GetContractGovernanceProposal().GetTitle(),
				Description: contractEvent.GetContractGovernanceProposal().GetDescription(),
				CreatedAt:   contractEvent.Timestamp,
				NotifyAt:    contractEvent.NotifyTime,
				EventType:   kafkaevent.EventType_GOVERNANCE,
			}, nil
		default:
			return nil, errors.New(fmt.Sprintf("No type defined for event %v", contractEvent.GetEvent()))
		}
	} else if entEvent.WalletEvent.WalletEvent != nil {
		var walletEvent = entEvent.WalletEvent.WalletEvent
		switch walletEvent.GetEvent().(type) {
		case *kafkaevent.WalletEvent_CoinReceived:
			var coinReceived = walletEvent.GetCoinReceived()
			return &eventpb.Event{
				Title:       "Token Received",
				Description: fmt.Sprintf("%v received %v%v from %v", walletEvent.GetWalletAddress(), coinReceived.GetCoin().Amount, coinReceived.GetCoin().Denom, coinReceived.Sender),
				CreatedAt:   walletEvent.Timestamp,
				EventType:   kafkaevent.EventType_FUNDING,
			}, nil
		case *kafkaevent.WalletEvent_OsmosisPoolUnlock:
			return &eventpb.Event{
				Title:       "Pool Unlock",
				Description: fmt.Sprintf("Unlock period of Osmosis pool for %v is over", walletEvent.GetWalletAddress()),
				CreatedAt:   walletEvent.Timestamp,
				EventType:   kafkaevent.EventType_DEX,
			}, nil
		case *kafkaevent.WalletEvent_Unstake:
			var unstake = walletEvent.GetUnstake()
			return &eventpb.Event{
				Title:       "Unstake",
				Description: fmt.Sprintf("Unbonding period for %v is over. %v %v Available", walletEvent.GetWalletAddress(), unstake.GetCoin().Amount, unstake.GetCoin().Denom),
				CreatedAt:   walletEvent.Timestamp,
				EventType:   kafkaevent.EventType_STAKING,
			}, nil
		case *kafkaevent.WalletEvent_NeutronTokenVesting:
			var neutronTokenVesting = walletEvent.GetNeutronTokenVesting()
			return &eventpb.Event{
				Title:       "Vesting Unlock",
				Description: fmt.Sprintf("Vesting period for %v is over. %v Neutron available.", walletEvent.GetWalletAddress(), neutronTokenVesting.GetAmount()/1000000),
				CreatedAt:   walletEvent.Timestamp,
				EventType:   kafkaevent.EventType_FUNDING,
			}, nil
		case *kafkaevent.WalletEvent_VoteReminder:
			var voteReminder = walletEvent.GetVoteReminder()
			return &eventpb.Event{
				Title:       fmt.Sprintf("Vote Reminder for Proposal %v", voteReminder.GetProposalId()),
				Description: fmt.Sprintf("%v did not vote yet", walletEvent.GetWalletAddress()),
				CreatedAt:   walletEvent.Timestamp,
				EventType:   kafkaevent.EventType_GOVERNANCE,
			}, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No type defined for event %v", entEvent))
}

func EntEventToProto(entEvent *ent.Event, chain *ent.Chain) (*eventpb.Event, error) {
	pbEvent, err := toProto(entEvent)
	if err != nil {
		return nil, err
	}
	pbEvent.Id = entEvent.ID.String()
	pbEvent.NotifyAt = timestamppb.New(entEvent.NotifyTime)
	pbEvent.Chain = &eventpb.ChainData{
		Id:       int64(chain.ID),
		Name:     chain.Name,
		ImageUrl: chain.Image,
	}
	pbEvent.Read = entEvent.IsRead
	return pbEvent, nil
}
