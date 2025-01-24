BEGIN;

CREATE TABLE IF NOT EXISTS scrambles (
  scramble_id BIGSERIAL PRIMARY KEY,
  scramble TEXT NOT NULL,
  event_id INTEGER REFERENCES events (event_id) NOT NULL,
  competition_id TEXT REFERENCES competitions (competition_id) NOT NULL,
  "order" INTEGER NOT NULL,
  img TEXT NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
