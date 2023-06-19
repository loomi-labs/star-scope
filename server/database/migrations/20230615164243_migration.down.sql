-- reverse: modify "users" table
ALTER TABLE "users" ADD COLUMN "wallet_address" character varying NOT NULL;
-- reverse: modify "comm_channels" table
ALTER TABLE "comm_channels" DROP COLUMN "wallet_address";
