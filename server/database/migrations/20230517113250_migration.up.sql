-- modify "chains" table
ALTER TABLE "chains" ADD COLUMN "handled_message_types" character varying NOT NULL DEFAULT '';
