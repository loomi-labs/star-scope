-- modify "user_setups" table
ALTER TABLE "user_setups" ALTER COLUMN "notify_funding" SET DEFAULT true, ALTER COLUMN "notify_staking" SET DEFAULT true, ALTER COLUMN "notify_gov_new_proposal" SET DEFAULT true;
