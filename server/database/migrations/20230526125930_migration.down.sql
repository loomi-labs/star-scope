-- reverse: modify "events" table
ALTER TABLE "events" DROP COLUMN "event_data", ADD COLUMN "tx_event" bytea NOT NULL;
