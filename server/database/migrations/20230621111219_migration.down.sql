-- reverse: modify "users" table
ALTER TABLE "users" DROP COLUMN "discord_username", DROP COLUMN "telegram_username", ADD COLUMN "name" character varying NOT NULL;
