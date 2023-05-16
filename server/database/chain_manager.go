package database

import (
	"context"
	"github.com/loomi-labs/star-scope/ent"
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

func (m *ChainManager) Update(ctx context.Context, id int, indexingHeight uint64, unhandledMessageTypes []string) error {
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

	uniqueUnhandledMessageTypes := make(map[string]bool)
	for _, element := range strings.Split(c.UnhandledMessageTypes, ",") {
		uniqueUnhandledMessageTypes[element] = true
	}
	for _, element := range unhandledMessageTypes {
		uniqueUnhandledMessageTypes[element] = true
	}
	var umt []string
	for element := range uniqueUnhandledMessageTypes {
		umt = append(umt, element)
	}
	return c.Update().
		SetIndexingHeight(indexingHeight).
		SetUnhandledMessageTypes(strings.Join(umt, ",")).
		Exec(ctx)
}
