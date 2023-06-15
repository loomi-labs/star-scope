-- modify "event_listeners" table
ALTER TABLE "event_listeners" ALTER COLUMN "wallet_address" DROP NOT NULL, ADD COLUMN "data_type" character varying NOT NULL;
