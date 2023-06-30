-- reverse: modify "validators" table
ALTER TABLE "validators" ADD COLUMN "logo_url" character varying NULL, ADD COLUMN "identity" character varying NOT NULL DEFAULT '';
