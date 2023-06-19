-- reverse: modify "comm_channels" table
ALTER TABLE "comm_channels" ALTER COLUMN "wallet_address" SET NOT NULL, ALTER COLUMN "is_group" DROP DEFAULT, ALTER COLUMN "discord_channel_id" SET NOT NULL, ALTER COLUMN "telegram_chat_id" SET NOT NULL;
