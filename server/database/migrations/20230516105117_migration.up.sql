-- modify "chains" table
ALTER TABLE "chains" ADD COLUMN "path" character varying NOT NULL, ADD COLUMN "has_custom_indexer" boolean NOT NULL DEFAULT false, ADD COLUMN "unhandled_message_types" character varying NOT NULL DEFAULT '';
