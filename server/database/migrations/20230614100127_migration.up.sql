-- modify "events" table
ALTER TABLE "events" ADD COLUMN "is_background" boolean NOT NULL DEFAULT false;
