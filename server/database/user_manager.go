package database

import (
	"context"
	"fmt"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/chain"
	"github.com/loomi-labs/star-scope/ent/commchannel"
	"github.com/loomi-labs/star-scope/ent/eventlistener"
	"github.com/loomi-labs/star-scope/ent/predicate"
	"github.com/loomi-labs/star-scope/ent/user"
	"github.com/loomi-labs/star-scope/grpc/settings/settingspb"
	"github.com/loomi-labs/star-scope/grpc/types"
	"github.com/loomi-labs/star-scope/kafka_internal"
	"github.com/shifty11/go-logger/log"
	"golang.org/x/exp/maps"
	"sort"
	"strings"
	"time"
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

func (m *UserManager) QueryByTelegram(ctx context.Context, tgUserId int64) (*ent.User, error) {
	return m.client.User.
		Query().
		Where(user.TelegramUserIDEQ(tgUserId)).
		Only(ctx)
}

func (m *UserManager) QueryByDiscord(ctx context.Context, discordUserId int64) (*ent.User, error) {
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

func (m *UserManager) QueryCommChannels(ctx context.Context, u *ent.User, t *commchannel.Type) ([]*ent.CommChannel, error) {
	var predicates []predicate.CommChannel
	if *t == commchannel.TypeTelegram {
		predicates = append(predicates, commchannel.TelegramChatIDNotNil())
	}
	if *t == commchannel.TypeDiscord {
		predicates = append(predicates, commchannel.DiscordChannelIDNotNil())
	}
	if *t == commchannel.TypeWebpush {
		return nil, fmt.Errorf("webpush not implemented")
	}
	return u.
		QueryCommChannels().
		Where(predicates...).
		Order(ent.Asc(commchannel.FieldName)).
		All(ctx)
}

func (m *UserManager) QuerySetup(ctx context.Context, u *ent.User) (*ent.UserSetup, error) {
	return u.
		QuerySetup().
		Only(ctx)
}

func (m *UserManager) QueryWallets(ctx context.Context, u *ent.User) ([]*settingspb.Wallet, error) {
	els, err := u.
		QueryEventListeners().
		Where(eventlistener.DataTypeIn(
			eventlistener.DataTypeWalletEvent_VoteReminder,
			eventlistener.DataTypeWalletEvent_Voted,
			eventlistener.DataTypeWalletEvent_CoinReceived,
			eventlistener.DataTypeWalletEvent_NeutronTokenVesting,
			eventlistener.DataTypeWalletEvent_OsmosisPoolUnlock,
			eventlistener.DataTypeWalletEvent_Unstake,
		)).
		Select(eventlistener.FieldWalletAddress, eventlistener.FieldDataType).
		WithChain(func(q *ent.ChainQuery) {
			q.Select(chain.FieldImage, chain.FieldIsQuerying, chain.FieldIsIndexing, chain.FieldIsEnabled)
		}).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var walletsMap = make(map[string]*settingspb.Wallet)
	for _, el := range els {
		if _, ok := walletsMap[el.WalletAddress]; !ok {
			walletsMap[el.WalletAddress] = &settingspb.Wallet{
				Address:                            el.WalletAddress,
				LogoUrl:                            el.Edges.Chain.Image,
				IsNotifyFundingSupported:           el.Edges.Chain.IsEnabled && el.Edges.Chain.IsIndexing,
				IsNotifyStakingSupported:           el.Edges.Chain.IsEnabled && el.Edges.Chain.IsIndexing,
				IsNotifyGovVotingReminderSupported: el.Edges.Chain.IsEnabled && el.Edges.Chain.IsQuerying,
			}
		}
		if el.DataType == eventlistener.DataTypeWalletEvent_VoteReminder {
			walletsMap[el.WalletAddress].NotifyGovVotingReminder = true
		}
		if el.DataType == eventlistener.DataTypeWalletEvent_CoinReceived {
			walletsMap[el.WalletAddress].NotifyFunding = true
		}
		if el.DataType == eventlistener.DataTypeWalletEvent_Unstake {
			walletsMap[el.WalletAddress].NotifyStaking = true
		}
	}
	var wallets = maps.Values(walletsMap)
	sort.Slice(wallets, func(i, j int) bool {
		return wallets[i].Address < wallets[j].Address
	})
	return wallets, nil
}

func (m *UserManager) UpdateRole(ctx context.Context, name string, role user.Role) (*ent.User, error) {
	entUser, err := m.client.User.
		Query().
		Where(user.Or(
			user.TelegramUsernameEQ(name),
			user.DiscordUsernameEQ(name),
			user.WalletAddressEQ(name),
		)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return entUser.
		Update().
		SetRole(role).
		Save(ctx)
}

func (m *UserManager) UpdateSetLoginDate(ctx context.Context, userId int) error {
	return m.client.User.
		UpdateOneID(userId).
		SetLastLoginTime(time.Now()).
		Exec(ctx)
}

func (m *UserManager) CreateByWalletAddress(ctx context.Context, walletAddress string) (*ent.User, error) {
	log.Sugar.Debugf("CreateByWalletAddress: %s", walletAddress)
	return withTxResult(m.client, ctx, func(tx *ent.Tx) (*ent.User, error) {
		u, err := tx.User.
			Create().
			SetWalletAddress(walletAddress).
			Save(ctx)
		if err != nil {
			return nil, err
		}
		err = tx.UserSetup.
			Create().
			SetUser(u).
			SetWalletAddresses([]string{walletAddress}).
			Exec(ctx)
		if err != nil {
			return nil, err
		}
		return u, nil
	})
}

func (m *UserManager) createOrAddTelegramCommChannel(ctx context.Context, tx *ent.Tx, u *ent.User, chatId int64, chatName string, isGroup bool) error {
	var eventListenerIds []int
	for _, eventListener := range u.Edges.EventListeners {
		eventListenerIds = append(eventListenerIds, eventListener.ID)
	}
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
			AddEventListenerIDs(eventListenerIds...).
			Exec(ctx)
	} else {
		return tx.CommChannel.
			Create().
			SetType(commchannel.TypeTelegram).
			SetName(chatName).
			SetTelegramChatID(chatId).
			SetIsGroup(isGroup).
			AddUsers(u).
			AddEventListenerIDs(eventListenerIds...).
			Exec(ctx)
	}
}

func (m *UserManager) CreateOrUpdateByTelegramUser(ctx context.Context, userId int64, userName string, chatId *int64, chatName *string, isGroup *bool) (*ent.User, error) {
	log.Sugar.Debugf("CreateOrUpdateByTelegramUser: %v (%v)", userName, userId)
	return withTxResult(m.client, ctx, func(tx *ent.Tx) (*ent.User, error) {
		u, err := tx.User.
			Query().
			Where(user.TelegramUserIDEQ(userId)).
			Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return u, err
		}
		if u == nil {
			u, err = tx.User.
				Create().
				SetTelegramUsername(userName).
				SetTelegramUserID(userId).
				Save(ctx)
			if err != nil {
				return u, err
			}
			if err := tx.UserSetup.Create().SetUser(u).Exec(ctx); err != nil {
				return nil, err
			}
		}
		if chatId == nil {
			return u, nil
		}
		return u, m.createOrAddTelegramCommChannel(ctx, tx, u, *chatId, *chatName, *isGroup)
	})
}

func (m *UserManager) createOrAddDiscordCommChannel(ctx context.Context, tx *ent.Tx, u *ent.User, channelId int64, channelName string, isGroup bool) error {
	var eventListenerIds []int
	for _, eventListener := range u.Edges.EventListeners {
		eventListenerIds = append(eventListenerIds, eventListener.ID)
	}
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
			AddEventListenerIDs(eventListenerIds...).
			Exec(ctx)
	} else {
		return tx.CommChannel.
			Create().
			SetType(commchannel.TypeDiscord).
			SetName(channelName).
			SetDiscordChannelID(channelId).
			SetIsGroup(isGroup).
			AddUsers(u).
			AddEventListenerIDs(eventListenerIds...).
			Exec(ctx)
	}
}

func (m *UserManager) CreateOrUpdateByDiscordUser(ctx context.Context, userId int64, userName string, channelId *int64, channelName *string, isGroup *bool) (*ent.User, error) {
	log.Sugar.Debugf("CreateOrUpdateByDiscordUser: %v (%v) in %v (%v)", userName, userId, channelName, channelId)
	return withTxResult(m.client, ctx, func(tx *ent.Tx) (*ent.User, error) {
		u, err := tx.User.
			Query().
			Where(user.DiscordUserIDEQ(userId)).
			WithEventListeners().
			Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return u, err
		}
		if u == nil {
			u, err = tx.User.
				Create().
				SetDiscordUsername(userName).
				SetDiscordUserID(userId).
				Save(ctx)
			if err != nil {
				return u, err
			}
			if err := tx.UserSetup.Create().SetUser(u).Exec(ctx); err != nil {
				return nil, err
			}
		}
		if channelId == nil {
			return u, nil
		}
		return u, m.createOrAddDiscordCommChannel(ctx, tx, u, *channelId, *channelName, *isGroup)
	})
}

func (m *UserManager) delete(ctx context.Context, tx *ent.Tx, u *ent.User) error {
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
}

func (m *UserManager) Delete(ctx context.Context, u *ent.User) error {
	err := withTx(m.client, ctx, func(tx *ent.Tx) error {
		return m.delete(ctx, tx, u)
	})
	if err != nil {
		return err
	}
	m.kafkaInternal.ProduceDbChangeMsg(kafka_internal.EventListenerDeleted)
	return nil
}

func (m *UserManager) DeleteTelegramCommChannel(ctx context.Context, userId int64, chatId int64, mustDeleteEmptyUsers bool) error {
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
	if cnt == 0 && mustDeleteEmptyUsers {
		u, err := m.client.User.Query().Where(user.TelegramUserID(userId)).Only(ctx)
		if err != nil {
			return err
		}
		return m.Delete(ctx, u)
	}
	return nil
}

func (m *UserManager) DeleteDiscordCommChannel(ctx context.Context, discordUserId int64, channelId int64, mustDeleteEmptyUsers bool) error {
	log.Sugar.Debugf("DeleteDiscordCommChannel: %d for %d", channelId, discordUserId)
	err := withTx(m.client, ctx, func(tx *ent.Tx) error {
		commChannel, err := tx.CommChannel.
			Query().
			Where(commchannel.And(
				commchannel.DiscordChannelIDEQ(channelId),
				commchannel.HasUsersWith(user.DiscordUserID(discordUserId)),
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
		Where(commchannel.HasUsersWith(user.DiscordUserID(discordUserId))).
		CountX(ctx)
	if cnt == 0 && mustDeleteEmptyUsers {
		u, err := m.client.User.Query().Where(user.DiscordUserID(discordUserId)).Only(ctx)
		if err != nil {
			return err
		}
		return m.Delete(ctx, u)
	}
	return nil
}

func (m *UserManager) moveEventListeners(ctx context.Context, tx *ent.Tx, newUser *ent.User, oldUser *ent.User) (bool, error) {
	oldEls, err := oldUser.
		QueryEventListeners().
		All(ctx)
	if err != nil {
		return false, err
	}
	els, err := tx.User.
		QueryEventListeners(newUser).
		WithUsers().
		All(ctx)
	if err != nil {
		return false, err
	}
	var isEventListenerDeleted = false
	for _, oldEl := range oldEls {
		var foundNewEl *ent.EventListener
		for _, el := range els {
			if oldEl.WalletAddress == el.WalletAddress && oldEl.DataType == el.DataType {
				foundNewEl = el
				break
			}
		}
		if foundNewEl != nil {
			if len(oldEl.Edges.Users) > 1 {
				err = tx.EventListener.
					UpdateOne(oldEl).
					RemoveUsers(oldUser).
					Exec(ctx)
				if err != nil {
					return false, err
				}
			} else {
				oldEvents := oldEl.
					QueryEvents().
					AllX(ctx)
				for _, event := range oldEvents {
					tx.Event.
						UpdateOne(event).
						SetEventListenerID(foundNewEl.ID).
						ExecX(ctx)
				}
				err = tx.EventListener.
					DeleteOne(oldEl).
					Exec(ctx)
				if err != nil {
					return false, err
				}
				isEventListenerDeleted = true
				// TODO: what happens to comm channels? Should we remove the deleted event listener from them?
			}
		} else {
			err = tx.EventListener.
				UpdateOne(oldEl).
				RemoveUsers(oldUser).
				AddUsers(newUser).
				Exec(ctx)
			if err != nil {
				return false, err
			}
		}
	}
	return isEventListenerDeleted, nil
}

func (m *UserManager) moveCommChannels(ctx context.Context, tx *ent.Tx, newUser *ent.User, oldUser *ent.User) error {
	oldCcs, err := oldUser.
		QueryCommChannels().
		All(ctx)
	if err != nil {
		return err
	}
	for _, oldCc := range oldCcs {
		err = tx.CommChannel.
			UpdateOne(oldCc).
			RemoveUsers(oldUser).
			AddUsers(newUser).
			Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *UserManager) UpdateConnectDiscord(ctx context.Context, u *ent.User, discord *types.DiscordIdentity) error {
	if u.DiscordUserID == discord.Id {
		return nil
	}
	if u.DiscordUserID != 0 {
		return fmt.Errorf("user %d already connected to discord", u.ID)
	}
	var isEventListenerDeleted = false
	err := withTx(m.client, ctx, func(tx *ent.Tx) error {
		oldDiscordUser, err := m.QueryByDiscord(ctx, discord.Id)
		if err != nil && !ent.IsNotFound(err) {
			return err
		} else if err == nil { // user already exists and has to be merged
			isEventListenerDeleted, err = m.moveEventListeners(ctx, tx, u, oldDiscordUser)
			if err != nil {
				return err
			}

			if err := m.moveCommChannels(ctx, tx, u, oldDiscordUser); err != nil {
				return err
			}
			err = m.delete(ctx, tx, oldDiscordUser)
			if err != nil {
				return err
			}
			updateQuery := tx.User.UpdateOne(u)
			if oldDiscordUser.TelegramUserID != 0 && u.TelegramUserID == 0 {
				updateQuery = updateQuery.SetTelegramUserID(u.TelegramUserID)
			}
			if oldDiscordUser.WalletAddress != "" && u.WalletAddress == "" {
				updateQuery = updateQuery.SetWalletAddress(u.WalletAddress)
			}
			return updateQuery.
				SetDiscordUserID(discord.Id).
				SetDiscordUsername(discord.Username).
				Exec(ctx)
		}

		return tx.User.
			UpdateOne(u).
			SetDiscordUserID(discord.Id).
			SetDiscordUsername(discord.Username).
			Exec(ctx)
	})
	if err != nil {
		return err
	}
	if isEventListenerDeleted {
		m.kafkaInternal.ProduceDbChangeMsg(kafka_internal.EventListenerDeleted)
	}
	return nil
}

func (m *UserManager) UpdateConnectTelegram(ctx context.Context, u *ent.User, data *types.TelegramLoginData) error {
	if u.TelegramUserID == data.UserId {
		return nil
	}
	if u.TelegramUserID != 0 {
		return fmt.Errorf("user %d already connected to telegram", u.ID)
	}
	var isEventListenerDeleted = false
	err := withTx(m.client, ctx, func(tx *ent.Tx) error {
		oldDiscordUser, err := m.QueryByTelegram(ctx, data.UserId)
		if err != nil && !ent.IsNotFound(err) {
			return err
		} else if err == nil { // user already exists and has to be merged
			isEventListenerDeleted, err = m.moveEventListeners(ctx, tx, u, oldDiscordUser)
			if err != nil {
				return err
			}

			if err := m.moveCommChannels(ctx, tx, u, oldDiscordUser); err != nil {
				return err
			}
			err = m.delete(ctx, tx, oldDiscordUser)
			if err != nil {
				return err
			}
			updateQuery := tx.User.UpdateOne(u)
			if oldDiscordUser.DiscordUserID != 0 && u.DiscordUserID == 0 {
				updateQuery = updateQuery.SetDiscordUserID(u.DiscordUserID)
			}
			if oldDiscordUser.WalletAddress != "" && u.WalletAddress == "" {
				updateQuery = updateQuery.SetWalletAddress(u.WalletAddress)
			}
			return updateQuery.
				SetTelegramUserID(data.UserId).
				SetTelegramUsername(data.Username).
				Exec(ctx)
		}

		return tx.User.
			UpdateOne(u).
			SetTelegramUserID(data.UserId).
			SetTelegramUsername(data.Username).
			Exec(ctx)
	})
	if err != nil {
		return err
	}
	if isEventListenerDeleted {
		m.kafkaInternal.ProduceDbChangeMsg(kafka_internal.EventListenerDeleted)
	}
	return nil
}

func getWalletEvents(entChain *ent.Chain, setup *ent.UserSetup) []eventlistener.DataType {
	var dt []eventlistener.DataType
	if setup.NotifyFunding {
		dt = append(dt, eventlistener.DataTypeWalletEvent_CoinReceived)
		if strings.Contains(entChain.Path, "neutron") {
			dt = append(dt, eventlistener.DataTypeWalletEvent_NeutronTokenVesting)
		}
		if strings.Contains(entChain.Path, "osmosis") {
			dt = append(dt, eventlistener.DataTypeWalletEvent_OsmosisPoolUnlock)
		}
	}
	if setup.NotifyStaking {
		if !strings.Contains(entChain.Path, "neutron") {
			dt = append(dt, eventlistener.DataTypeWalletEvent_Unstake)
		}
	}
	if setup.NotifyGovVotingReminder {
		dt = append(dt, eventlistener.DataTypeWalletEvent_Voted)
		dt = append(dt, eventlistener.DataTypeWalletEvent_VoteReminder)
	}
	return dt
}

func getChainEvents(entChain *ent.Chain, setup *ent.UserSetup) []eventlistener.DataType {
	var dt []eventlistener.DataType
	if strings.Contains(entChain.Path, "neutron") {
		return dt
	}
	if setup.NotifyGovNewProposal {
		dt = append(dt, eventlistener.DataTypeChainEvent_GovernanceProposal_Ongoing)
	}
	if setup.NotifyGovVotingEnd {
		dt = append(dt, eventlistener.DataTypeChainEvent_GovernanceProposal_Finished)
	}
	if setup.NotifyStaking {
		dt = append(dt, eventlistener.DataTypeChainEvent_ValidatorOutOfActiveSet)
	}
	return dt
}

func getContractEvents(entChain *ent.Chain) []eventlistener.DataType {
	var dt []eventlistener.DataType
	if strings.Contains(entChain.Path, "neutron") {
		dt = append(dt, eventlistener.DataTypeContractEvent_ContractGovernanceProposal_Ongoing)
		dt = append(dt, eventlistener.DataTypeContractEvent_ContractGovernanceProposal_Finished)
	}
	return dt
}

func (m *UserManager) createEventListeners(
	ctx context.Context,
	tx *ent.Tx,
	setup *ent.UserSetup,
	entUser *ent.User,
	chains []*ent.Chain,
	selectedChains []*ent.Chain,
	commChannels []*ent.CommChannel,
	selectedValidators []*ent.Validator,
) error {
	var bulk []*ent.EventListenerCreate
	for _, entChain := range chains {
		for _, address := range setup.WalletAddresses {
			if common.IsBech32AddressFromChain(address, entChain.Bech32Prefix) {
				for _, dt := range getWalletEvents(entChain, setup) {
					bulk = append(bulk, tx.EventListener.
						Create().
						SetChain(entChain).
						AddUsers(entUser).
						AddCommChannels(commChannels...).
						SetWalletAddress(address).
						SetDataType(dt))
				}
				break
			}
		}
	}
	for _, entChain := range selectedChains {
		for _, dt := range getChainEvents(entChain, setup) {
			bulk = append(bulk, tx.EventListener.
				Create().
				SetChain(entChain).
				AddUsers(entUser).
				AddCommChannels(commChannels...).
				SetDataType(dt))
		}
		for _, dt := range getContractEvents(entChain) {
			bulk = append(bulk, tx.EventListener.
				Create().
				SetChain(entChain).
				AddUsers(entUser).
				AddCommChannels(commChannels...).
				SetDataType(dt))
		}
	}
	valEventListeners := []eventlistener.DataType{
		eventlistener.DataTypeWalletEvent_Voted,
		eventlistener.DataTypeWalletEvent_VoteReminder,
	}
	for _, entValidator := range selectedValidators {
		for _, dt := range valEventListeners {
			bulk = append(bulk, tx.EventListener.
				Create().
				SetChain(entValidator.Edges.Chain).
				AddUsers(entUser).
				AddCommChannels(commChannels...).
				SetWalletAddress(entValidator.Address).
				SetDataType(dt))
		}
	}
	return tx.EventListener.
		CreateBulk(bulk...).
		Exec(ctx)
}

func (m *UserManager) UpdateSetup(ctx context.Context, u *ent.User, query *ent.UserSetupUpdateOne, isCompleted bool, availableChains []*ent.Chain) (*ent.UserSetup, error) {
	if isCompleted {
		result, err := withTxResult(m.client, ctx, func(tx *ent.Tx) (*ent.UserSetup, error) {
			err := tx.User.
				UpdateOne(u).
				SetIsSetupComplete(true).
				Exec(ctx)
			if err != nil {
				return nil, err
			}

			setup, err := query.Save(ctx)
			if err != nil {
				return nil, err
			}

			selectedChains, err := setup.
				QuerySelectedChains().
				All(ctx)
			if err != nil {
				return nil, err
			}

			selectedValidators, err := setup.
				QuerySelectedValidators().
				WithChain().
				All(ctx)
			if err != nil {
				return nil, err
			}

			commChannels, err := tx.CommChannel.
				Query().
				Where(commchannel.HasUsersWith(user.ID(u.ID))).
				All(ctx)
			if err != nil {
				return nil, err
			}

			err = m.createEventListeners(ctx, tx, setup, u, availableChains, selectedChains, commChannels, selectedValidators)
			if err != nil {
				return nil, err
			}
			return setup, nil
		})
		if err != nil {
			return nil, err
		}
		m.kafkaInternal.ProduceDbChangeMsg(kafka_internal.EventListenerCreated)
		return result, nil
	}
	return query.Save(ctx)
}
