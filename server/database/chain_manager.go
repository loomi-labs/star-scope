package database

import (
	"context"
	"errors"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/chain"
	"github.com/loomi-labs/star-scope/ent/contractproposal"
	"github.com/loomi-labs/star-scope/ent/eventlistener"
	"github.com/loomi-labs/star-scope/ent/proposal"
	kafkaevent "github.com/loomi-labs/star-scope/event"
	"github.com/loomi-labs/star-scope/types"
	"github.com/shifty11/go-logger/log"
	"golang.org/x/exp/slices"
	"strings"
	"time"
)

type ChainManager struct {
	client *ent.Client
}

func NewChainManager(client *ent.Client) *ChainManager {
	return &ChainManager{client: client}
}

func (m *ChainManager) QueryAll(ctx context.Context) []*ent.Chain {
	return m.client.Chain.
		Query().
		AllX(ctx)
}

func (m *ChainManager) QueryEnabled(ctx context.Context) []*ent.Chain {
	return m.client.Chain.
		Query().
		Where(chain.IsEnabledEQ(true)).
		AllX(ctx)
}

func (m *ChainManager) QueryByBech32Prefix(ctx context.Context, bech32Prefix string) (*ent.Chain, error) {
	return m.client.Chain.
		Query().
		Where(chain.Bech32Prefix(bech32Prefix)).
		First(ctx)
}

func (m *ChainManager) QueryByName(ctx context.Context, name string) (*ent.Chain, error) {
	return m.client.Chain.
		Query().
		Where(chain.Or(
			chain.NameEQ(name),
			chain.PrettyNameEQ(name),
		)).
		Only(ctx)
}

func (m *ChainManager) QueryById(background context.Context, id int) (*ent.Chain, error) {
	return m.client.Chain.
		Get(background, id)
}

func (m *ChainManager) QueryProposals(ctx context.Context, entChain *ent.Chain) []*ent.Proposal {
	return entChain.
		QueryProposals().
		AllX(ctx)
}

func (m *ChainManager) QueryContractProposals(ctx context.Context, entChain *ent.Chain) []*ent.ContractProposal {
	return entChain.
		QueryContractProposals().
		AllX(ctx)
}

func (m *ChainManager) QueryNewAccounts(ctx context.Context, entChain *ent.Chain) []*ent.EventListener {
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	return entChain.
		QueryEventListeners().
		Where(
			eventlistener.CreateTimeGTE(oneHourAgo),
		).
		Select(eventlistener.FieldWalletAddress).
		AllX(ctx)
}

func (m *ChainManager) Create(ctx context.Context, chainData *types.ChainData) (*ent.Chain, error) {
	log.Sugar.Debugf("Creating chain: %v", chainData.PrettyName)
	return m.client.Chain.
		Create().
		SetChainID(chainData.ChainId).
		SetName(chainData.Name).
		SetPrettyName(chainData.PrettyName).
		SetBech32Prefix(chainData.Bech32Prefix).
		SetPath(chainData.Path).
		SetImage(chainData.Image).
		Save(ctx)
}

func (m *ChainManager) UpdateChainInfo(ctx context.Context, entChain *ent.Chain, chainData *types.ChainData) (*ent.Chain, error) {
	log.Sugar.Debugf("Creating chain: %v", chainData.PrettyName)
	return entChain.
		Update().
		SetChainID(chainData.ChainId).
		SetName(chainData.Name).
		SetPrettyName(chainData.PrettyName).
		SetBech32Prefix(chainData.Bech32Prefix).
		SetImage(chainData.Image).
		Save(ctx)
}

func (m *ChainManager) UpdateSetEnabled(ctx context.Context, entChain *ent.Chain, isEnabled bool) (*ent.Chain, error) {
	return entChain.
		Update().
		SetIsEnabled(isEnabled).
		Save(ctx)
}

func getUniqueMessageTypes(messageTypes []string, forbiddenMessageTypes []string) []string {
	uniqueMessageTypes := make(map[string]bool)
	for _, element := range messageTypes {
		uniqueMessageTypes[element] = true
	}
	var umt []string
	for element := range uniqueMessageTypes {
		if element != "" && !slices.Contains(forbiddenMessageTypes, element) {
			umt = append(umt, element)
		}
	}
	return umt
}

func (m *ChainManager) UpdateIndexStatus(
	ctx context.Context,
	id int,
	indexingHeight uint64,
	handledMessageTypes []string,
	unhandledMessageTypes []string,
) error {
	if unhandledMessageTypes == nil {
		return m.client.Chain.
			UpdateOneID(id).
			SetIndexingHeight(indexingHeight).
			Exec(context.Background())
	}
	c, err := m.client.Chain.Get(ctx, id)
	if err != nil {
		return err
	}

	handledMessageTypes = append(handledMessageTypes, strings.Split(c.HandledMessageTypes, ",")...)
	unhandledMessageTypes = append(unhandledMessageTypes, strings.Split(c.UnhandledMessageTypes, ",")...)
	handledMessageTypes = getUniqueMessageTypes(handledMessageTypes, nil)
	unhandledMessageTypes = getUniqueMessageTypes(unhandledMessageTypes, handledMessageTypes)

	return c.Update().
		SetIndexingHeight(indexingHeight).
		SetHandledMessageTypes(strings.Join(handledMessageTypes, ",")).
		SetUnhandledMessageTypes(strings.Join(unhandledMessageTypes, ",")).
		Exec(ctx)
}

func (m *ChainManager) createProposal(ctx context.Context, entChain *ent.Chain, prop *kafkaevent.GovernanceProposalEvent) (*ent.Proposal, error) {
	return m.client.Proposal.
		Create().
		SetChain(entChain).
		SetProposalID(prop.ProposalId).
		SetStatus(proposal.Status(prop.GetProposalStatus().String())).
		SetTitle(prop.GetTitle()).
		SetDescription(prop.GetDescription()).
		SetVotingStartTime(prop.GetVotingStartTime().AsTime()).
		SetVotingEndTime(prop.GetVotingEndTime().AsTime()).
		Save(ctx)
}

func (m *ChainManager) UpdateProposal(ctx context.Context, entChain *ent.Chain, govProp *kafkaevent.GovernanceProposalEvent) (*ent.Proposal, error) {
	if govProp == nil {
		return nil, errors.New("governance prop is nil")
	}
	prop, err := entChain.
		QueryProposals().
		Where(proposal.ProposalIDEQ(govProp.ProposalId)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return m.createProposal(ctx, entChain, govProp)
	} else if err != nil {
		return nil, err
	}
	return prop.
		Update().
		SetStatus(proposal.Status(govProp.GetProposalStatus().String())).
		Save(ctx)
}

func (m *ChainManager) createContractProposal(ctx context.Context, entChain *ent.Chain, prop *kafkaevent.ContractGovernanceProposalEvent) (*ent.ContractProposal, error) {
	return m.client.ContractProposal.
		Create().
		SetChain(entChain).
		SetProposalID(prop.ProposalId).
		SetStatus(contractproposal.Status(prop.GetProposalStatus().String())).
		SetTitle(prop.GetTitle()).
		SetDescription(prop.GetDescription()).
		SetContractAddress(prop.GetContractAddress()).
		SetFirstSeenTime(prop.GetFirstSeenTime().AsTime()).
		SetVotingEndTime(prop.GetVotingEndTime().AsTime()).
		Save(ctx)
}

func (m *ChainManager) UpdateContractProposal(ctx context.Context, entChain *ent.Chain, govProp *kafkaevent.ContractGovernanceProposalEvent) (*ent.ContractProposal, error) {
	if govProp == nil {
		return nil, errors.New("governance prop is nil")
	}
	prop, err := entChain.
		QueryContractProposals().
		Where(contractproposal.And(
			contractproposal.ProposalIDEQ(govProp.ProposalId),
			contractproposal.ContractAddressEQ(govProp.GetContractAddress()),
		)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return m.createContractProposal(ctx, entChain, govProp)
	} else if err != nil {
		return nil, err
	}
	return prop.
		Update().
		SetStatus(contractproposal.Status(govProp.GetProposalStatus().String())).
		Save(ctx)
}
