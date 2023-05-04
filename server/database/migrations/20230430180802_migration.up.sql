-- modify "chains" table
ALTER TABLE "chains" ADD COLUMN "indexing_height" bigint NOT NULL DEFAULT 0;
