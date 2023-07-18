-- reverse: modify "chains" table
ALTER TABLE "chains" DROP COLUMN "last_successful_validator_query", DROP COLUMN "last_successful_proposal_query";
