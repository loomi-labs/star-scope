-- reverse: create index "users_wallet_address_key" to table: "users"
DROP INDEX "users_wallet_address_key";
-- reverse: create index "users_telegram_user_id_key" to table: "users"
DROP INDEX "users_telegram_user_id_key";
-- reverse: create index "users_discord_user_id_key" to table: "users"
DROP INDEX "users_discord_user_id_key";
-- reverse: create index "user_wallet_address" to table: "users"
DROP INDEX "user_wallet_address";
-- reverse: create index "user_telegram_user_id" to table: "users"
DROP INDEX "user_telegram_user_id";
-- reverse: create index "user_discord_user_id" to table: "users"
DROP INDEX "user_discord_user_id";
-- reverse: modify "users" table
ALTER TABLE "users" DROP COLUMN "wallet_address", DROP COLUMN "discord_user_id", DROP COLUMN "telegram_user_id";
-- reverse: create index "commchannel_wallet_address" to table: "comm_channels"
DROP INDEX "commchannel_wallet_address";
-- reverse: create index "commchannel_telegram_chat_id" to table: "comm_channels"
DROP INDEX "commchannel_telegram_chat_id";
-- reverse: create index "commchannel_discord_channel_id" to table: "comm_channels"
DROP INDEX "commchannel_discord_channel_id";
