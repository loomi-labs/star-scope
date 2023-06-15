-- reverse: create "user_comm_channels" table
DROP TABLE "user_comm_channels";
-- reverse: create "user_event_listeners" table
DROP TABLE "user_event_listeners";
-- reverse: create "comm_channel_event_listeners" table
DROP TABLE "comm_channel_event_listeners";
-- reverse: modify "event_listeners" table
ALTER TABLE "event_listeners" ADD COLUMN "user_event_listeners" bigint NULL;
-- reverse: create index "comm_channels_telegram_chat_id_key" to table: "comm_channels"
DROP INDEX "comm_channels_telegram_chat_id_key";
-- reverse: create index "comm_channels_discord_channel_id_key" to table: "comm_channels"
DROP INDEX "comm_channels_discord_channel_id_key";
-- reverse: create "comm_channels" table
DROP TABLE "comm_channels";
