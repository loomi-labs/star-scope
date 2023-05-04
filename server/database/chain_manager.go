package database

import (
	"context"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/chain"
)

type ChainManager struct {
	client *ent.Client
}

func NewChainManager(client *ent.Client) *ChainManager {
	return &ChainManager{client: client}
}

func (m *ChainManager) QueryByName(ctx context.Context, name string) (*ent.Chain, error) {
	return m.client.Chain.
		Query().
		Where(chain.NameEQ(name)).
		Only(ctx)
}

func (m *ChainManager) UpdateIndexingHeight(ctx context.Context, name string, indexingHeight int64) (int, error) {
	return m.client.Chain.
		Update().
		Where(chain.NameEQ(name)).
		SetIndexingHeight(indexingHeight).
		Save(ctx)
}
