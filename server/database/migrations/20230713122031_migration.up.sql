-- modify "chains" table
ALTER TABLE "chains" ADD COLUMN "is_querying" boolean NOT NULL DEFAULT false, ADD COLUMN "is_indexing" boolean NOT NULL DEFAULT false;
