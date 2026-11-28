BEGIN;

CREATE INDEX idx_scrambles_profile_query ON scrambles (competition_id, event_id, "order");

COMMIT;
