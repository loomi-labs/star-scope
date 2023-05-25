-- reverse: create "proposals" table
DROP TABLE "proposals";
-- reverse: modify "event_listeners" table
ALTER TABLE "event_listeners" DROP CONSTRAINT "event_listeners_chains_event_listeners", ADD CONSTRAINT "event_listeners_chains_event_listeners" FOREIGN KEY ("chain_event_listeners") REFERENCES "chains" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
