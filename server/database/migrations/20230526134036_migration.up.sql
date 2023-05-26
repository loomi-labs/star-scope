-- modify "events" table
ALTER TABLE "events" DROP COLUMN "type", ADD COLUMN "event_type" character varying NOT NULL, ADD COLUMN "data_type" character varying NOT NULL;
