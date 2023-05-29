package kafka

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/grpc/event/eventpb"
	"github.com/loomi-labs/star-scope/indexevent"
	"github.com/loomi-labs/star-scope/queryevent"
	"github.com/shifty11/go-logger/log"
)

func txEventToProto(data []byte) (uint64, *eventpb.Event, error) {
	var txEvent indexevent.TxEvent
	err := proto.Unmarshal(data, &txEvent)
	if err != nil {
		return 0, nil, err
	}
	switch txEvent.GetEvent().(type) {
	case *indexevent.TxEvent_CoinReceived:
		var coinReceived = txEvent.GetCoinReceived()
		return txEvent.ChainId, &eventpb.Event{
			Title:       "Token Received",
			Description: fmt.Sprintf("%v received %v%v from %v", txEvent.GetWalletAddress(), coinReceived.GetCoin().Amount, coinReceived.GetCoin().Denom, coinReceived.Sender),
			CreatedAt:   txEvent.Timestamp,
			NotifyAt:    txEvent.NotifyTime,
			EventType:   eventpb.EventType_FUNDING,
		}, nil
	case *indexevent.TxEvent_OsmosisPoolUnlock:
		return txEvent.ChainId, &eventpb.Event{
			Title:       "Pool Unlock",
			Description: fmt.Sprintf("Unlock period of Osmosis pool for %v is over", txEvent.GetWalletAddress()),
			CreatedAt:   txEvent.Timestamp,
			NotifyAt:    txEvent.NotifyTime,
			EventType:   eventpb.EventType_DEX,
		}, nil
	case *indexevent.TxEvent_Unstake:
		var unstake = txEvent.GetUnstake()
		return txEvent.ChainId, &eventpb.Event{
			Title:       "Unstake",
			Description: fmt.Sprintf("Unbonding period for %v is over. %v %v Available", txEvent.GetWalletAddress(), unstake.GetCoin().Amount, unstake.GetCoin().Denom),
			CreatedAt:   txEvent.Timestamp,
			NotifyAt:    txEvent.NotifyTime,
			EventType:   eventpb.EventType_STAKING,
		}, nil
	case *indexevent.TxEvent_NeutronTokenVesting:
		var neutronTokenVesting = txEvent.GetNeutronTokenVesting()
		return txEvent.ChainId, &eventpb.Event{
			Title:       "Vesting Unlock",
			Description: fmt.Sprintf("Vesting period for %v is over. %v Neutron available.", txEvent.GetWalletAddress(), neutronTokenVesting.GetAmount()/1000000),
			CreatedAt:   txEvent.Timestamp,
			NotifyAt:    txEvent.NotifyTime,
			EventType:   eventpb.EventType_FUNDING,
		}, nil
	}
	return 0, nil, errors.New(fmt.Sprintf("No type defined for event %v", txEvent.GetEvent()))
}

func queryEventToProto(data []byte) (uint64, *eventpb.Event, error) {
	var queryEvent queryevent.QueryEvent
	err := proto.Unmarshal(data, &queryEvent)
	if err != nil {
		return 0, nil, err
	}
	switch queryEvent.GetEvent().(type) {
	case *queryevent.QueryEvent_GovernanceProposal:
		var statusText = "Proposal %v"
		switch queryEvent.GetGovernanceProposal().GetProposalStatus() {
		case queryevent.ProposalStatus_PROPOSAL_STATUS_VOTING_PERIOD:
			statusText = "New Proposal - %v"
		case queryevent.ProposalStatus_PROPOSAL_STATUS_PASSED:
			statusText = "Proposal %v Passed"
		case queryevent.ProposalStatus_PROPOSAL_STATUS_REJECTED:
			statusText = "Proposal %v Rejected"
		case queryevent.ProposalStatus_PROPOSAL_STATUS_FAILED:
			statusText = "Proposal %v Failed"
		default:
			log.Sugar.Errorf("Unknown proposal status %v", queryEvent.GetGovernanceProposal().GetProposalStatus())
		}
		return queryEvent.ChainId, &eventpb.Event{
			Title:       fmt.Sprintf(statusText, queryEvent.GetGovernanceProposal().GetProposalId()),
			Subtitle:    queryEvent.GetGovernanceProposal().GetTitle(),
			Description: queryEvent.GetGovernanceProposal().GetDescription(),
			CreatedAt:   queryEvent.Timestamp,
			NotifyAt:    queryEvent.NotifyTime,
			EventType:   eventpb.EventType_GOVERNANCE,
		}, nil
	case *queryevent.QueryEvent_ContractGovernanceProposal:
		var statusText = "Proposal %v"
		switch queryEvent.GetContractGovernanceProposal().GetProposalStatus() {
		case queryevent.ContractProposalStatus_OPEN:
			statusText = "New Proposal - %v"
		case queryevent.ContractProposalStatus_PASSED:
			statusText = "Proposal %v Passed"
		case queryevent.ContractProposalStatus_REJECTED:
			statusText = "Proposal %v Rejected"
		case queryevent.ContractProposalStatus_EXECUTION_FAILED:
			statusText = "Proposal %v Failed"
		case queryevent.ContractProposalStatus_CLOSED:
			statusText = "Proposal %v Closed"
		}
		return queryEvent.ChainId, &eventpb.Event{
			Title:       fmt.Sprintf(statusText, queryEvent.GetContractGovernanceProposal().GetProposalId()),
			Subtitle:    queryEvent.GetContractGovernanceProposal().GetTitle(),
			Description: queryEvent.GetContractGovernanceProposal().GetDescription(),
			CreatedAt:   queryEvent.Timestamp,
			NotifyAt:    queryEvent.NotifyTime,
			EventType:   eventpb.EventType_GOVERNANCE,
		}, nil
	}
	return 0, nil, errors.New(fmt.Sprintf("No type defined for event %v", queryEvent.GetEvent()))
}

func kafkaMsgToProto(data []byte, chains []*ent.Chain) (*eventpb.Event, error) {
	var chainId, pbEvent, err = txEventToProto(data)
	if err != nil {
		chainId, pbEvent, err = queryEventToProto(data)
		if err != nil {
			return nil, err
		}
	}
	for _, chain := range chains {
		if uint64(chain.ID) == chainId {
			pbEvent.Chain = &eventpb.ChainData{
				Id:       int64(chain.ID),
				Name:     chain.Name,
				ImageUrl: chain.Image,
			}
			break
		}
	}
	return pbEvent, nil
}

func EntEventToProto(entEvent *ent.Event, chain *ent.Chain) (*eventpb.Event, error) {
	var pbEvent *eventpb.Event
	var err error
	if entEvent.IsTxEvent {
		_, pbEvent, err = txEventToProto(entEvent.Data)
		if err != nil {
			return nil, err
		}
	} else {
		_, pbEvent, err = queryEventToProto(entEvent.Data)
		if err != nil {
			return nil, err
		}
	}
	pbEvent.Id = int64(entEvent.ID)
	pbEvent.Chain = &eventpb.ChainData{
		Id:       int64(chain.ID),
		Name:     chain.Name,
		ImageUrl: chain.Image,
	}
	return pbEvent, nil
}
