package database

import (
	"context"
	"github.com/shifty11/blocklog-backend/ent"
	"github.com/shifty11/blocklog-backend/ent/user"
	"github.com/shifty11/go-logger/log"
)

type UserManager struct {
	client *ent.Client
}

func NewUserManager(client *ent.Client) *UserManager {
	return &UserManager{client: client}
}

func (m *UserManager) QueryById(ctx context.Context, id int) (*ent.User, error) {
	return m.client.User.Get(ctx, id)
}

func (m *UserManager) QueryByWalletAddress(ctx context.Context, walletAddress string) (*ent.User, error) {
	return m.client.User.
		Query().
		Where(user.WalletAddressEQ(walletAddress)).
		Only(ctx)
}

func (m *UserManager) QueryAdmins(ctx context.Context) ([]*ent.User, error) {
	return m.client.User.
		Query().
		Where(user.RoleEQ(user.RoleAdmin)).
		All(ctx)
}

func (m *UserManager) UpdateRole(ctx context.Context, name string, role user.Role) (*ent.User, error) {
	entUser, err := m.client.User.
		Query().
		Where(user.NameEQ(name)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return entUser.
		Update().
		SetRole(role).
		Save(ctx)
}

func (m *UserManager) CreateOrUpdate(ctx context.Context, userName string, walletAddress string) *ent.User {
	log.Sugar.Debugf("CreateOrUpdate: %s %s", userName, walletAddress)
	entUser, err := m.client.User.
		Query().
		Where(user.WalletAddressEQ(walletAddress)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			entUser, err = m.client.User.
				Create().
				SetName(userName).
				SetWalletAddress(walletAddress).
				Save(ctx)
			if err != nil {
				log.Sugar.Panicf("Could not create user: %v", err)
			}
		} else {
			log.Sugar.Panicf("Could not find user: %v", err)
		}
	} else if entUser.Name != userName {
		entUser, err = m.client.User.
			UpdateOne(entUser).
			SetName(userName).
			Save(ctx)
		if err != nil {
			log.Sugar.Panicf("Could not update user: %v", err)
		}
	}
	return entUser
}
