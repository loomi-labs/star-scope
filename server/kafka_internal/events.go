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
			var emoji = ""
			var statusText = ""
			switch chainEvent.GetGovernanceProposal().GetProposalStatus() {
			case kafkaevent.ProposalStatus_PROPOSAL_STATUS_VOTING_PERIOD:
				emoji = "🗳"
				statusText = "Proposal %v"
			case kafkaevent.ProposalStatus_PROPOSAL_STATUS_PASSED:
				emoji = "✅"
				statusText = "Proposal %v passed"
			case kafkaevent.ProposalStatus_PROPOSAL_STATUS_REJECTED:
				emoji = "❌"
				statusText = "Proposal %v rejected"
			case kafkaevent.ProposalStatus_PROPOSAL_STATUS_FAILED:
				emoji = "❌"
				statusText = "Proposal %v failed"
			default:
				log.Sugar.Errorf("Unknown proposal status %v", chainEvent.GetGovernanceProposal().GetProposalStatus())
			}
			return &eventpb.Event{
				Title:       fmt.Sprintf(statusText, chainEvent.GetGovernanceProposal().GetProposalId()),
				Subtitle:    chainEvent.GetGovernanceProposal().GetTitle(),
				Description: chainEvent.GetGovernanceProposal().GetDescription(),
				Emoji:       emoji,
				CreatedAt:   chainEvent.Timestamp,
				EventType:   eventpb.EventType_GOVERNANCE,
			}, nil
		default:
			return nil, errors.New(fmt.Sprintf("No type defined for event %v", chainEvent.GetEvent()))
		}
	} else if entEvent.ContractEvent.ContractEvent != nil {
		var contractEvent = entEvent.ContractEvent.ContractEvent
		switch contractEvent.GetEvent().(type) {
		case *kafkaevent.ContractEvent_ContractGovernanceProposal:
			var emoji = ""
			var statusText = ""
			switch contractEvent.GetContractGovernanceProposal().GetProposalStatus() {
			case kafkaevent.ContractProposalStatus_OPEN:
				emoji = "🗳"
				statusText = "Proposal %v"
			case kafkaevent.ContractProposalStatus_PASSED:
				emoji = "✅"
				statusText = "Proposal %v passed"
			case kafkaevent.ContractProposalStatus_REJECTED:
				emoji = "❌"
				statusText = "Proposal %v rejected"
			case kafkaevent.ContractProposalStatus_EXECUTION_FAILED:
				emoji = "❌"
				statusText = "Proposal %v failed"
			case kafkaevent.ContractProposalStatus_CLOSED:
				emoji = "❌"
				statusText = "Proposal %v closed"
			}
			return &eventpb.Event{
				Title:       fmt.Sprintf(statusText, contractEvent.GetContractGovernanceProposal().GetProposalId()),
				Subtitle:    contractEvent.GetContractGovernanceProposal().GetTitle(),
				Description: contractEvent.GetContractGovernanceProposal().GetDescription(),
				Emoji:       emoji,
				CreatedAt:   contractEvent.Timestamp,
				NotifyAt:    contractEvent.NotifyTime,
				EventType:   eventpb.EventType_GOVERNANCE,
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
				Emoji:       "💰",
				CreatedAt:   walletEvent.Timestamp,
				EventType:   eventpb.EventType_FUNDING,
			}, nil
		case *kafkaevent.WalletEvent_OsmosisPoolUnlock:
			return &eventpb.Event{
				Title:       "Pool Unlock",
				Description: fmt.Sprintf("Unlock period of Osmosis pool for %v is over", walletEvent.GetWalletAddress()),
				Emoji:       "🔓",
				CreatedAt:   walletEvent.Timestamp,
				EventType:   eventpb.EventType_DEX,
			}, nil
		case *kafkaevent.WalletEvent_Unstake:
			var unstake = walletEvent.GetUnstake()
			return &eventpb.Event{
				Title:       "Unstake",
				Description: fmt.Sprintf("Unbonding period for %v is over. %v %v Available", walletEvent.GetWalletAddress(), unstake.GetCoin().Amount, unstake.GetCoin().Denom),
				Emoji:       "🔓",
				CreatedAt:   walletEvent.Timestamp,
				EventType:   eventpb.EventType_STAKING,
			}, nil
		case *kafkaevent.WalletEvent_NeutronTokenVesting:
			var neutronTokenVesting = walletEvent.GetNeutronTokenVesting()
			return &eventpb.Event{
				Title:       "Vesting Unlock",
				Description: fmt.Sprintf("Vesting period for %v is over. %v Neutron available.", walletEvent.GetWalletAddress(), neutronTokenVesting.GetAmount()/1000000),
				Emoji:       "🔓",
				CreatedAt:   walletEvent.Timestamp,
				EventType:   eventpb.EventType_FUNDING,
			}, nil
		case *kafkaevent.WalletEvent_VoteReminder:
			var voteReminder = walletEvent.GetVoteReminder()
			wallet := ""
			if walletEvent.GetWalletName() != "" {
				wallet = fmt.Sprintf("%v (%v)", walletEvent.GetWalletName(), walletEvent.GetWalletAddress())
			} else {
				wallet = walletEvent.GetWalletAddress()
			}
			return &eventpb.Event{
				Title:       fmt.Sprintf("Vote Reminder for Proposal %v", voteReminder.GetProposalId()),
				Description: fmt.Sprintf("%v did not vote yet", wallet),
				Emoji:       "🗳",
				CreatedAt:   walletEvent.Timestamp,
				EventType:   eventpb.EventType_GOVERNANCE,
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
		Name:     chain.PrettyName,
		ImageUrl: chain.Image,
	}
	pbEvent.Read = entEvent.IsRead
	return pbEvent, nil
}
