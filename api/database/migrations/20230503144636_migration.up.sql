-- modify "events" table
ALTER TABLE "events" DROP COLUMN "title", DROP COLUMN "description", ADD COLUMN "type" character varying NOT NULL, ADD COLUMN "tx_event" bytea NOT NULL;
