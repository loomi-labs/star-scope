package database

import (
	"context"
	"github.com/shifty11/blocklog-backend/ent"
	"github.com/shifty11/blocklog-backend/ent/chain"
	"github.com/shifty11/blocklog-backend/ent/channel"
	"github.com/shifty11/go-logger/log"
)

type ProjectManager struct {
	client *ent.Client
}

func NewProjectManager(client *ent.Client) *ProjectManager {
	return &ProjectManager{client: client}
}

func (m *ProjectManager) QueryById(ctx context.Context, id int) (*ent.User, error) {
	return m.client.User.Get(ctx, id)
}

func (m *ProjectManager) CreateCosmosProject(ctx context.Context, user *ent.User, walletAddress string) (*ent.Project, error) {
	log.Sugar.Debugf("CreateCosmosProject: %s %s", user.Name, walletAddress)
	osmosisChain, err := m.client.Chain.
		Query().
		Where(chain.NameEQ("Osmosis")).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	proj, err := m.client.Project.
		Create().
		SetName("Cosmos").
		SetUser(user).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	ch, err := m.client.Channel.
		Create().
		SetName("Funding").
		SetType(channel.TypeFunding).
		SetProject(proj).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	_, err = m.client.EventListener.
		Create().
		SetChannel(ch).
		SetChain(osmosisChain).
		SetWalletAddress(walletAddress).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return proj, nil
}
