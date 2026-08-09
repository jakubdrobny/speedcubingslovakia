BEGIN;

ALTER TABLE competition_events ADD COLUMN format TEXT;
UPDATE competition_events ce SET format = e.format FROM events e WHERE ce.event_id = e.event_id;

COMMIT;
