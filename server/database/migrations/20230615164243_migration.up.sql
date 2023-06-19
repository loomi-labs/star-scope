-- modify "comm_channels" table
ALTER TABLE "comm_channels" ADD COLUMN "wallet_address" character varying NOT NULL;
-- modify "users" table
ALTER TABLE "users" DROP COLUMN "wallet_address";
