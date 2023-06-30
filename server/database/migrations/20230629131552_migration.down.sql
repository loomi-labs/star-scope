-- reverse: modify "user_setups" table
ALTER TABLE "user_setups" ALTER COLUMN "notify_gov_new_proposal" SET DEFAULT false, ALTER COLUMN "notify_staking" SET DEFAULT false, ALTER COLUMN "notify_funding" SET DEFAULT false;
