-- modify "user_setups" table
ALTER TABLE "user_setups" DROP COLUMN "is_finished";
-- modify "users" table
ALTER TABLE "users" ADD COLUMN "is_setup_complete" boolean NOT NULL DEFAULT false;
-- modify "chains" table
ALTER TABLE "chains" ADD COLUMN "user_setup_selected_chains" bigint NULL, ADD CONSTRAINT "chains_user_setups_selected_chains" FOREIGN KEY ("user_setup_selected_chains") REFERENCES "user_setups" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
-- modify "validators" table
ALTER TABLE "validators" DROP COLUMN "user_setup_validators", ADD COLUMN "user_setup_selected_validators" bigint NULL, ADD CONSTRAINT "validators_user_setups_selected_validators" FOREIGN KEY ("user_setup_selected_validators") REFERENCES "user_setups" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
