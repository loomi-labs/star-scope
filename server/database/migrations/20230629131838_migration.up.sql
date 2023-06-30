-- drop index "user_setups_user_setup_key" from table: "user_setups"
DROP INDEX "user_setups_user_setup_key";
-- modify "user_setups" table
ALTER TABLE "user_setups" DROP CONSTRAINT "user_setups_users_setup", ADD CONSTRAINT "user_setups_users_setup" FOREIGN KEY ("user_setup") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
