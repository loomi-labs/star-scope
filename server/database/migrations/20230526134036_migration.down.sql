-- reverse: modify "events" table
ALTER TABLE "events" DROP COLUMN "data_type", DROP COLUMN "event_type", ADD COLUMN "type" character varying NOT NULL;
