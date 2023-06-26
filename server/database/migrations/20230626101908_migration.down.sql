-- reverse: modify "validators" table
ALTER TABLE "validators" DROP CONSTRAINT "validators_user_setups_validators", DROP COLUMN "user_setup_validators", DROP COLUMN "logo_url", DROP COLUMN "identity";
-- reverse: create index "user_setups_user_setup_key" to table: "user_setups"
DROP INDEX "user_setups_user_setup_key";
-- reverse: create "user_setups" table
DROP TABLE "user_setups";
