-- reverse: modify "events" table
ALTER TABLE "events" DROP COLUMN "tx_event", DROP COLUMN "type", ADD COLUMN "description" character varying NOT NULL, ADD COLUMN "title" character varying NOT NULL;
