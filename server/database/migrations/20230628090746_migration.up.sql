-- modify "chains" table
ALTER TABLE "chains" DROP COLUMN "user_setup_selected_chains";
-- create "chain_selected_by_setups" table
CREATE TABLE "chain_selected_by_setups" ("chain_id" bigint NOT NULL, "user_setup_id" bigint NOT NULL, PRIMARY KEY ("chain_id", "user_setup_id"), CONSTRAINT "chain_selected_by_setups_chain_id" FOREIGN KEY ("chain_id") REFERENCES "chains" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "chain_selected_by_setups_user_setup_id" FOREIGN KEY ("user_setup_id") REFERENCES "user_setups" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- modify "validators" table
ALTER TABLE "validators" DROP COLUMN "user_setup_selected_validators";
-- create "validator_selected_by_setups" table
CREATE TABLE "validator_selected_by_setups" ("validator_id" bigint NOT NULL, "user_setup_id" bigint NOT NULL, PRIMARY KEY ("validator_id", "user_setup_id"), CONSTRAINT "validator_selected_by_setups_user_setup_id" FOREIGN KEY ("user_setup_id") REFERENCES "user_setups" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "validator_selected_by_setups_validator_id" FOREIGN KEY ("validator_id") REFERENCES "validators" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
