-- reverse: modify "chains" table
ALTER TABLE "chains" DROP COLUMN "unhandled_message_types", DROP COLUMN "has_custom_indexer", DROP COLUMN "path";
