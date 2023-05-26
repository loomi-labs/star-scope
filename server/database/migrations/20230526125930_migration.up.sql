-- modify "events" table
ALTER TABLE "events" DROP COLUMN "tx_event", ADD COLUMN "event_data" bytea NOT NULL;
