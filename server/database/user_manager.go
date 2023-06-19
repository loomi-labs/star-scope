package database

import (
	"context"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/commchannel"
	"github.com/loomi-labs/star-scope/ent/user"
	"github.com/loomi-labs/star-scope/kafka_internal"
	"github.com/shifty11/go-logger/log"
)

type UserManager struct {
	client        *ent.Client
	kafkaInternal kafka_internal.KafkaInternal
}

func NewUserManager(client *ent.Client, kafkaInternal kafka_internal.KafkaInternal) *UserManager {
	return &UserManager{client: client, kafkaInternal: kafkaInternal}
}

func (m *UserManager) StartTx(ctx context.Context) (*ent.Tx, error) {
	return m.client.Tx(ctx)
}

func (m *UserManager) WithTx(ctx context.Context, fn func(tx *ent.Tx) error) error {
	return withTx(m.client, ctx, func(tx *ent.Tx) error {
		return fn(tx)
	})
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

func (m *UserManager) QueryByTelegram(ctx context.Context, tgChatId int64) (*ent.User, error) {
	return m.client.User.
		Query().
		Where(user.TelegramUserIDEQ(tgChatId)).
		Only(ctx)
}

func (m *UserManager) QueryByDiscordChannel(ctx context.Context, discordUserId int64) (*ent.User, error) {
	return m.client.User.
		Query().
		Where(user.DiscordUserIDEQ(discordUserId)).
		Only(ctx)
}

func (m *UserManager) QueryAdmins(ctx context.Context) ([]*ent.User, error) {
	return m.client.User.
		Query().
		Where(user.RoleEQ(user.RoleAdmin)).
		All(ctx)
}

func (m *UserManager) QueryUsersForTelegramChat(ctx context.Context, chatId int64) []*ent.User {
	users, err := m.client.User.
		Query().
		Where(user.HasCommChannelsWith(commchannel.TelegramChatID(chatId))).
		All(ctx)
	if err != nil {
		log.Sugar.Errorf("Could not get users for telegram chat: %v", err)
	}
	return users
}

func (m *UserManager) QueryUsersForDiscordChannel(ctx context.Context, channelId int64) []*ent.User {
	users, err := m.client.User.
		Query().
		Where(user.HasCommChannelsWith(commchannel.DiscordChannelIDEQ(channelId))).
		All(ctx)
	if err != nil {
		log.Sugar.Errorf("Could not get users for discord channel: %v", err)
	}
	return users
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

func (m *UserManager) CreateByWalletAddress(ctx context.Context, tx *ent.Tx, walletAddress string) (*ent.User, error) {
	log.Sugar.Debugf("CreateByWalletAddress: %s", walletAddress)
	return tx.User.
		Create().
		SetName(walletAddress).
		SetWalletAddress(walletAddress).
		Save(ctx)
}

func (m *UserManager) createOrAddTelegramCommChannel(ctx context.Context, tx *ent.Tx, u *ent.User, chatId int64, chatName string, isGroup bool) error {
	commChannel, err := tx.CommChannel.
		Query().
		Where(commchannel.TelegramChatIDEQ(chatId)).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return err
	}
	if commChannel != nil {
		return commChannel.Update().
			AddUsers(u).
			Exec(ctx)
	} else {
		return tx.CommChannel.
			Create().
			SetType(commchannel.TypeTelegram).
			SetName(chatName).
			SetTelegramChatID(chatId).
			SetIsGroup(isGroup).
			AddUsers(u).
			Exec(ctx)
	}
}

func (m *UserManager) CreateOrUpdateForTelegramUser(ctx context.Context, userId int64, userName string, chatId int64, chatName string, isGroup bool) error {
	log.Sugar.Debugf("CreateOrUpdateForTelegramUser: %v (%v)", userName, userId)
	return withTx(m.client, ctx, func(tx *ent.Tx) error {
		u, err := tx.User.
			Query().
			Where(user.HasCommChannelsWith(commchannel.TelegramChatIDEQ(userId))).
			Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return err
		}
		if u == nil {
			u, err = tx.User.
				Create().
				SetName(userName).
				SetTelegramUserID(userId).
				Save(ctx)
			if err != nil {
				return err
			}
		}
		return m.createOrAddTelegramCommChannel(ctx, tx, u, chatId, chatName, isGroup)
	})
}

func (m *UserManager) createOrAddDiscordCommChannel(ctx context.Context, tx *ent.Tx, u *ent.User, channelId int64, channelName string, isGroup bool) error {
	commChannel, err := tx.CommChannel.
		Query().
		Where(commchannel.DiscordChannelIDEQ(channelId)).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return err
	}
	if commChannel != nil {
		return commChannel.Update().
			AddUsers(u).
			Exec(ctx)
	} else {
		return tx.CommChannel.
			Create().
			SetType(commchannel.TypeDiscord).
			SetName(channelName).
			SetDiscordChannelID(channelId).
			SetIsGroup(isGroup).
			AddUsers(u).
			Exec(ctx)
	}
}

func (m *UserManager) CreateOrUpdateForDiscordUser(ctx context.Context, userId int64, userName string, channelId int64, channelName string, isGroup bool) error {
	log.Sugar.Debugf("CreateOrUpdateForDiscordUser: %s (%d) in %s (%d)", userName, userId, channelName, channelId)
	return withTx(m.client, ctx, func(tx *ent.Tx) error {
		u, err := tx.User.
			Query().
			Where(user.DiscordUserIDEQ(userId)).
			Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return err
		}
		if u == nil {
			u, err = tx.User.
				Create().
				SetName(userName).
				SetDiscordUserID(userId).
				Save(ctx)
			if err != nil {
				return err
			}
		}
		return m.createOrAddDiscordCommChannel(ctx, tx, u, channelId, channelName, isGroup)
	})
}

func (m *UserManager) Delete(ctx context.Context, u *ent.User) error {
	err := withTx(m.client, ctx, func(tx *ent.Tx) error {
		els, err := tx.User.QueryEventListeners(u).All(ctx)
		if err != nil {
			return err
		}
		for _, el := range els {
			cntUsers, err := tx.EventListener.QueryUsers(el).Count(ctx)
			if err != nil {
				return err
			}
			if cntUsers <= 1 {
				err = tx.EventListener.DeleteOne(el).Exec(ctx)
				if err != nil {
					return err
				}
			}
		}
		ccs, err := tx.User.QueryCommChannels(u).All(ctx)
		if err != nil {
			return err
		}
		for _, cc := range ccs {
			cntUsers, err := tx.CommChannel.QueryUsers(cc).Count(ctx)
			if err != nil {
				return err
			}
			if cntUsers <= 1 {
				err = tx.CommChannel.DeleteOne(cc).Exec(ctx)
				if err != nil {
					return err
				}
			}
		}
		return tx.User.DeleteOne(u).Exec(ctx)
	})
	if err != nil {
		return err
	}
	m.kafkaInternal.ProduceDbChangeMsg(kafka_internal.EventListenerDeleted)
	return nil
}

func (m *UserManager) DeleteTelegramCommChannel(ctx context.Context, userId int64, chatId int64) error {
	log.Sugar.Debugf("DeleteTelegramCommChannel: %d for %d", chatId, userId)
	err := withTx(m.client, ctx, func(tx *ent.Tx) error {
		commChannel, err := tx.CommChannel.
			Query().
			Where(commchannel.And(
				commchannel.TelegramChatIDEQ(chatId),
				commchannel.HasUsersWith(user.TelegramUserID(userId)),
			)).
			Only(ctx)
		if err != nil {
			return err
		}
		err = commChannel.Update().
			ClearUsers().
			Exec(ctx)
		if err != nil {
			return err
		}
		return tx.CommChannel.DeleteOne(commChannel).Exec(ctx)
	})
	if err != nil {
		return err
	}
	cnt := m.client.CommChannel.Query().
		Where(commchannel.HasUsersWith(user.TelegramUserID(userId))).
		CountX(ctx)
	if cnt == 0 {
		u, err := m.client.User.Query().Where(user.TelegramUserID(userId)).Only(ctx)
		if err != nil {
			return err
		}
		return m.Delete(ctx, u)
	}
	return nil
}

func (m *UserManager) DeleteDiscordCommChannel(ctx context.Context, userId int64, channelId int64) error {
	log.Sugar.Debugf("DeleteDiscordCommChannel: %d for %d", channelId, userId)
	err := withTx(m.client, ctx, func(tx *ent.Tx) error {
		commChannel, err := tx.CommChannel.
			Query().
			Where(commchannel.And(
				commchannel.DiscordChannelIDEQ(channelId),
				commchannel.HasUsersWith(user.DiscordUserID(userId)),
			)).
			Only(ctx)
		if err != nil {
			return err
		}
		err = commChannel.Update().
			ClearUsers().
			Exec(ctx)
		if err != nil {
			return err
		}
		return tx.CommChannel.DeleteOne(commChannel).Exec(ctx)
	})
	if err != nil {
		return err
	}
	cnt := m.client.CommChannel.Query().
		Where(commchannel.HasUsersWith(user.DiscordUserID(userId))).
		CountX(ctx)
	if cnt == 0 {
		u, err := m.client.User.Query().Where(user.DiscordUserID(userId)).Only(ctx)
		if err != nil {
			return err
		}
		return m.Delete(ctx, u)
	}
	return nil
}
