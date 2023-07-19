package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/chain"
	"github.com/loomi-labs/star-scope/ent/contractproposal"
	"github.com/loomi-labs/star-scope/ent/eventlistener"
	"github.com/loomi-labs/star-scope/ent/proposal"
	"github.com/loomi-labs/star-scope/grpc/settings/settingspb"
	"github.com/loomi-labs/star-scope/kafka_internal"
	"github.com/loomi-labs/star-scope/types"
	"github.com/shifty11/go-logger/log"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"sort"
	"strings"
	"time"
)

type ChainManager struct {
	client        *ent.Client
	kafkaInternal kafka_internal.KafkaInternal
}

func NewChainManager(client *ent.Client, kafkaInternal kafka_internal.KafkaInternal) *ChainManager {
	return &ChainManager{client: client, kafkaInternal: kafkaInternal}
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

func (m *ChainManager) QueryIsQuerying(ctx context.Context) []*ent.Chain {
	return m.client.Chain.
		Query().
		Where(chain.IsQuerying(true)).
		AllX(ctx)
}

func (m *ChainManager) QueryIsIndexing(ctx context.Context) []*ent.Chain {
	return m.client.Chain.
		Query().
		Where(chain.IsIndexing(true)).
		AllX(ctx)
}

func (m *ChainManager) QueryIsQueryingWithProposals(ctx context.Context) []*ent.Chain {
	return m.client.Chain.
		Query().
		Where(chain.IsQuerying(true)).
		WithProposals().
		AllX(ctx)
}

func (m *ChainManager) QueryByBech32Prefix(ctx context.Context, bech32Prefix string) (*ent.Chain, error) {
	return m.client.Chain.
		Query().
		Where(chain.Bech32Prefix(bech32Prefix)).
		First(ctx)
}

func (m *ChainManager) QueryByNewAddress(ctx context.Context, address string) (*ent.Chain, error) {
	for _, c := range m.QueryEnabled(ctx) {
		if common.IsBech32AddressFromChain(address, c.Bech32Prefix) {
			return c, nil
		}
	}
	return nil, errors.New("no chain found for address")
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

func (m *ChainManager) QuerySubscribedChains(ctx context.Context, entUser *ent.User) ([]*settingspb.Chain, error) {
	els, err := entUser.
		QueryEventListeners().
		Where(eventlistener.DataTypeIn(
			eventlistener.DataTypeChainEvent_GovernanceProposal_Ongoing,
			eventlistener.DataTypeChainEvent_GovernanceProposal_Finished,
			//eventlistener.DataTypeChainEvent_ValidatorOutOfActiveSet,
			//eventlistener.DataTypeChainEvent_ValidatorSlash,
		)).
		Select(eventlistener.FieldDataType).
		WithChain().
		All(ctx)
	if err != nil {
		return nil, err
	}

	var chainMap = make(map[int]*settingspb.Chain)
	for _, el := range els {
		if _, ok := chainMap[el.Edges.Chain.ID]; !ok {
			chainMap[el.Edges.Chain.ID] = &settingspb.Chain{
				Id:                                uint64(el.Edges.Chain.ID),
				Name:                              el.Edges.Chain.PrettyName,
				LogoUrl:                           el.Edges.Chain.Image,
				NotifyNewProposals:                false,
				NotifyProposalFinished:            false,
				IsNotifyNewProposalsSupported:     el.Edges.Chain.IsEnabled && el.Edges.Chain.IsQuerying,
				IsNotifyProposalFinishedSupported: el.Edges.Chain.IsEnabled && el.Edges.Chain.IsQuerying,
			}
		}
		if el.DataType == eventlistener.DataTypeChainEvent_GovernanceProposal_Ongoing {
			chainMap[el.Edges.Chain.ID].NotifyNewProposals = true
		}
		if el.DataType == eventlistener.DataTypeChainEvent_GovernanceProposal_Finished {
			chainMap[el.Edges.Chain.ID].NotifyProposalFinished = true
		}
	}
	var chains = maps.Values(chainMap)
	sort.Slice(chains, func(i, j int) bool {
		return chains[i].Name < chains[j].Name
	})
	return chains, nil
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
		SetRestEndpoint(fmt.Sprintf("https://rest.cosmos.directory/%v", chainData.Path)).
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

func (m *ChainManager) UpdateSetEnabled(ctx context.Context, entChain *ent.Chain, isEnabled bool, isQuerying bool, isIndexing bool, height *uint64) (*ent.Chain, error) {
	query := entChain.
		Update().
		SetIsEnabled(isEnabled).
		SetIsQuerying(isQuerying).
		SetIsIndexing(isIndexing)
	if height != nil {
		query = query.SetIndexingHeight(*height)
	}
	updated, err := query.Save(ctx)
	if err != nil {
		return updated, err
	}
	if isEnabled || isQuerying || isIndexing {
		m.kafkaInternal.ProduceDbChangeMsg(kafka_internal.ChainEnabled)
	} else {
		m.kafkaInternal.ProduceDbChangeMsg(kafka_internal.ChainDisabled)
	}
	return updated, nil
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
) (*ent.Chain, error) {
	if unhandledMessageTypes == nil {
		return m.client.Chain.
			UpdateOneID(id).
			SetIndexingHeight(indexingHeight).
			Save(context.Background())
	}
	c, err := m.client.Chain.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	handledMessageTypes = append(handledMessageTypes, strings.Split(c.HandledMessageTypes, ",")...)
	unhandledMessageTypes = append(unhandledMessageTypes, strings.Split(c.UnhandledMessageTypes, ",")...)
	handledMessageTypes = getUniqueMessageTypes(handledMessageTypes, nil)
	unhandledMessageTypes = getUniqueMessageTypes(unhandledMessageTypes, handledMessageTypes)

	return c.Update().
		SetIndexingHeight(indexingHeight).
		SetHandledMessageTypes(strings.Join(handledMessageTypes, ",")).
		SetUnhandledMessageTypes(strings.Join(unhandledMessageTypes, ",")).
		Save(ctx)
}

func (m *ChainManager) UpdateSetLastSuccessfulProposalQuery(ctx context.Context, entChain *ent.Chain) {
	err := entChain.
		Update().
		SetLastSuccessfulProposalQuery(time.Now()).
		Exec(ctx)
	if err != nil {
		log.Sugar.Errorf("Error updating last successful proposal query: %v", err)
	}
}

func (m *ChainManager) UpdateSetLastSuccessfulValidatorQuery(ctx context.Context, entChain *ent.Chain) {
	err := entChain.
		Update().
		SetLastSuccessfulValidatorQuery(time.Now()).
		Exec(ctx)
	if err != nil {
		log.Sugar.Errorf("Error updating last successful validator query: %v", err)
	}
}

func (m *ChainManager) createProposal(ctx context.Context, entChain *ent.Chain, prop *types.Proposal) (*ent.Proposal, error) {
	return m.client.Proposal.
		Create().
		SetChain(entChain).
		SetProposalID(uint64(prop.ProposalId)).
		SetStatus(proposal.Status(prop.Status.String())).
		SetTitle(prop.Content.Title).
		SetDescription(prop.Content.Description).
		SetVotingStartTime(prop.VotingStartTime).
		SetVotingEndTime(prop.VotingEndTime).
		Save(ctx)
}

func (m *ChainManager) CreateOrUpdateProposal(ctx context.Context, entChain *ent.Chain, govProp *types.Proposal) (*ent.Proposal, error) {
	if govProp == nil {
		return nil, errors.New("governance prop is nil")
	}
	prop, err := entChain.
		QueryProposals().
		Where(proposal.ProposalIDEQ(uint64(govProp.ProposalId))).
		Only(ctx)
	if ent.IsNotFound(err) {
		return m.createProposal(ctx, entChain, govProp)
	} else if err != nil {
		return nil, err
	}
	return prop.
		Update().
		SetStatus(proposal.Status(govProp.Status.String())).
		Save(ctx)
}

func (m *ChainManager) createContractProposal(ctx context.Context, entChain *ent.Chain, propId uint64, contractAddress string, govProp *types.ContractProposal) (*ent.ContractProposal, error) {
	return m.client.ContractProposal.
		Create().
		SetChain(entChain).
		SetProposalID(propId).
		SetStatus(contractproposal.Status(govProp.Status.String())).
		SetTitle(govProp.Title).
		SetDescription(govProp.Description).
		SetContractAddress(contractAddress).
		SetFirstSeenTime(time.Now()).
		SetVotingEndTime(time.Time(govProp.Expiration.AtTime)).
		Save(ctx)
}

func (m *ChainManager) CreateOrUpdateContractProposal(ctx context.Context, entChain *ent.Chain, propId uint64, contractAddress string, govProp *types.ContractProposal) (*ent.ContractProposal, error) {
	if govProp == nil {
		return nil, errors.New("governance prop is nil")
	}
	prop, err := entChain.
		QueryContractProposals().
		Where(contractproposal.And(
			contractproposal.ProposalIDEQ(propId),
			contractproposal.ContractAddressEQ(contractAddress),
		)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return m.createContractProposal(ctx, entChain, propId, contractAddress, govProp)
	} else if err != nil {
		return nil, err
	}
	return prop.
		Update().
		SetStatus(contractproposal.Status(govProp.Status.String())).
		Save(ctx)
}
