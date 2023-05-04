-- modify "events" table
ALTER TABLE "events" ADD COLUMN "notify_time" timestamptz NOT NULL;
