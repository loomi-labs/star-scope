-- modify "event_listeners" table
ALTER TABLE "event_listeners" DROP CONSTRAINT "event_listeners_chains_event_listeners", ADD CONSTRAINT "event_listeners_chains_event_listeners" FOREIGN KEY ("chain_event_listeners") REFERENCES "chains" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- create "proposals" table
CREATE TABLE "proposals" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "create_time" timestamptz NOT NULL, "update_time" timestamptz NOT NULL, "proposal_id" bigint NOT NULL, "title" character varying NOT NULL, "description" character varying NOT NULL, "voting_start_time" timestamptz NOT NULL, "voting_end_time" timestamptz NOT NULL, "status" character varying NOT NULL, "chain_proposals" bigint NULL, PRIMARY KEY ("id"), CONSTRAINT "proposals_chains_proposals" FOREIGN KEY ("chain_proposals") REFERENCES "chains" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
