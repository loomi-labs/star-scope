-- modify "chains" table
ALTER TABLE "chains" ADD COLUMN "last_successful_proposal_query" timestamptz NULL, ADD COLUMN "last_successful_validator_query" timestamptz NULL;
