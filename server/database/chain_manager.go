package database

import (
	"context"
	"github.com/loomi-labs/star-scope/ent"
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

func (m *ChainManager) Update(
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
