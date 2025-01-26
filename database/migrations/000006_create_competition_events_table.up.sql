BEGIN;

CREATE TABLE IF NOT EXISTS competition_events (
  competition_events_id BIGSERIAL PRIMARY KEY,
  competition_id TEXT REFERENCES competitions (competition_id) NOT NULL,
  event_id INTEGER REFERENCES events (event_id) NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
