-- reverse: modify "events" table
ALTER TABLE "events" DROP COLUMN "is_read", DROP COLUMN "is_tx_event", DROP COLUMN "data", ADD COLUMN "event_data" bytea NOT NULL;
