package database

import (
	"context"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/grpc/types"
	"github.com/loomi-labs/star-scope/kafka_internal"
	"testing"
)

func newTestUserManager(t *testing.T) *UserManager {
	manager := NewUserManager(testClient(t), kafka_internal.NewKafkaInternalDummy())
	t.Cleanup(func() { closeTestClient(manager.client) })
	return manager
}

func TestUserManager_UpdateConnectDiscord(t *testing.T) {
	m := newTestUserManager(t)

	ctx := context.Background()
	user1, err := withTxResult(m.client, ctx, func(tx *ent.Tx) (*ent.User, error) {
		return m.CreateByWalletAddress(ctx, "cosmos15tp8np2adn47620394wm0jt4sjpw95t8um08xe")
	})
	if err != nil {
		t.Fatal(err)
	}
	user2, err := m.CreateOrUpdateByDiscordUser(ctx, int64(1), "discordUser", nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	dIdent := types.DiscordIdentity{
		Id:       user2.DiscordUserID,
		Username: user2.DiscordUsername,
	}
	err = m.UpdateConnectDiscord(ctx, user1, &dIdent)
	if err != nil {
		t.Fatal(err)
	}

	if m.client.User.Query().CountX(ctx) != 1 {
		t.Fatal("expected 1 userByDiscord")
	}

	userByWallet, err := m.QueryByWalletAddress(ctx, "cosmos15tp8np2adn47620394wm0jt4sjpw95t8um08xe")
	if err != nil {
		t.Fatal(err)
	}
	if userByWallet.WalletAddress != "cosmos15tp8np2adn47620394wm0jt4sjpw95t8um08xe" &&
		userByWallet.DiscordUserID != user2.DiscordUserID &&
		userByWallet.DiscordUsername != user2.DiscordUsername {
		t.Fatal("expected userByDiscord to be updated")
	}

	_, err = m.QueryByDiscord(ctx, user2.DiscordUserID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserManager_UpdateConnectTelegram(t *testing.T) {
	m := newTestUserManager(t)

	ctx := context.Background()
	user1, err := withTxResult(m.client, ctx, func(tx *ent.Tx) (*ent.User, error) {
		return m.CreateByWalletAddress(ctx, "cosmos15tp8np2adn47620394wm0jt4sjpw95t8um08xe")
	})
	if err != nil {
		t.Fatal(err)
	}
	user2, err := m.CreateOrUpdateByTelegramUser(ctx, int64(1), "telegramUser", nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	tIdent := types.TelegramLoginData{
		UserId:   user2.TelegramUserID,
		Username: user2.TelegramUsername,
	}
	err = m.UpdateConnectTelegram(ctx, user1, &tIdent)
	if err != nil {
		t.Fatal(err)
	}

	if m.client.User.Query().CountX(ctx) != 1 {
		t.Fatal("expected 1 userByTelegram")
	}

	userByWallet, err := m.QueryByWalletAddress(ctx, "cosmos15tp8np2adn47620394wm0jt4sjpw95t8um08xe")
	if err != nil {
		t.Fatal(err)
	}
	if userByWallet.WalletAddress != "cosmos15tp8np2adn47620394wm0jt4sjpw95t8um08xe" &&
		userByWallet.TelegramUserID != user2.TelegramUserID &&
		userByWallet.TelegramUsername != user2.TelegramUsername {
		t.Fatal("expected userByTelegram to be updated")
	}

	_, err = m.QueryByTelegram(ctx, user2.TelegramUserID)
	if err != nil {
		t.Fatal(err)
	}
}
