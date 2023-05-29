-- reverse: modify "events" table
ALTER TABLE "events" DROP CONSTRAINT "events_event_listeners_events", ADD CONSTRAINT "events_event_listeners_events" FOREIGN KEY ("event_listener_events") REFERENCES "event_listeners" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
-- reverse: modify "event_listeners" table
ALTER TABLE "event_listeners" DROP CONSTRAINT "event_listeners_users_event_listeners", ADD CONSTRAINT "event_listeners_users_event_listeners" FOREIGN KEY ("user_event_listeners") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
