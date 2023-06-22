-- modify "users" table
ALTER TABLE "users" DROP COLUMN "name", ADD COLUMN "telegram_username" character varying NULL, ADD COLUMN "discord_username" character varying NULL;
