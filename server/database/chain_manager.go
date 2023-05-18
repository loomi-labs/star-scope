package database

import (
	"context"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/chain"
	"github.com/loomi-labs/star-scope/types"
	"github.com/shifty11/go-logger/log"
	"golang.org/x/exp/slices"
	"strings"
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

func (m *ChainManager) QueryByName(ctx context.Context, name string) []*ent.Chain {
	return m.client.Chain.
		Query().
		Where(chain.Or(
			chain.NameEQ(name),
			chain.PrettyNameEQ(name),
		)).
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
