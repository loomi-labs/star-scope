-- reverse: create "validator_selected_by_setups" table
DROP TABLE "validator_selected_by_setups";
-- reverse: modify "validators" table
ALTER TABLE "validators" ADD COLUMN "user_setup_selected_validators" bigint NULL;
-- reverse: create "chain_selected_by_setups" table
DROP TABLE "chain_selected_by_setups";
-- reverse: modify "chains" table
ALTER TABLE "chains" ADD COLUMN "user_setup_selected_chains" bigint NULL;
