-- reverse: modify "event_listeners" table
ALTER TABLE "event_listeners" DROP COLUMN "data_type", ALTER COLUMN "wallet_address" SET NOT NULL;
