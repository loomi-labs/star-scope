-- modify "comm_channels" table
ALTER TABLE "comm_channels" ALTER COLUMN "telegram_chat_id" DROP NOT NULL, ALTER COLUMN "discord_channel_id" DROP NOT NULL, ALTER COLUMN "is_group" SET DEFAULT false, ALTER COLUMN "wallet_address" DROP NOT NULL;
