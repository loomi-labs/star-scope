-- reverse: modify "validators" table
ALTER TABLE "validators" DROP CONSTRAINT "validators_user_setups_selected_validators", DROP COLUMN "user_setup_selected_validators", ADD COLUMN "user_setup_validators" bigint NULL;
-- reverse: modify "chains" table
ALTER TABLE "chains" DROP CONSTRAINT "chains_user_setups_selected_chains", DROP COLUMN "user_setup_selected_chains";
-- reverse: modify "users" table
ALTER TABLE "users" DROP COLUMN "is_setup_complete";
-- reverse: modify "user_setups" table
ALTER TABLE "user_setups" ADD COLUMN "is_finished" boolean NOT NULL DEFAULT false;
