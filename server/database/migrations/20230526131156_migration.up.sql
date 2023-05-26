-- modify "events" table
ALTER TABLE "events" DROP COLUMN "event_data", ADD COLUMN "data" bytea NOT NULL, ADD COLUMN "is_tx_event" boolean NOT NULL, ADD COLUMN "is_read" boolean NOT NULL DEFAULT false;
