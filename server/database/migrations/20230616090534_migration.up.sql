-- create index "commchannel_discord_channel_id" to table: "comm_channels"
CREATE INDEX "commchannel_discord_channel_id" ON "comm_channels" ("discord_channel_id");
-- create index "commchannel_telegram_chat_id" to table: "comm_channels"
CREATE INDEX "commchannel_telegram_chat_id" ON "comm_channels" ("telegram_chat_id");
-- create index "commchannel_wallet_address" to table: "comm_channels"
CREATE INDEX "commchannel_wallet_address" ON "comm_channels" ("wallet_address");
-- modify "users" table
ALTER TABLE "users" ADD COLUMN "telegram_user_id" bigint NULL, ADD COLUMN "discord_user_id" bigint NULL, ADD COLUMN "wallet_address" character varying NULL;
-- create index "user_discord_user_id" to table: "users"
CREATE INDEX "user_discord_user_id" ON "users" ("discord_user_id");
-- create index "user_telegram_user_id" to table: "users"
CREATE INDEX "user_telegram_user_id" ON "users" ("telegram_user_id");
-- create index "user_wallet_address" to table: "users"
CREATE INDEX "user_wallet_address" ON "users" ("wallet_address");
-- create index "users_discord_user_id_key" to table: "users"
CREATE UNIQUE INDEX "users_discord_user_id_key" ON "users" ("discord_user_id");
-- create index "users_telegram_user_id_key" to table: "users"
CREATE UNIQUE INDEX "users_telegram_user_id_key" ON "users" ("telegram_user_id");
-- create index "users_wallet_address_key" to table: "users"
CREATE UNIQUE INDEX "users_wallet_address_key" ON "users" ("wallet_address");
