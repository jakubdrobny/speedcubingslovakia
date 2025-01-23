BEGIN;

CREATE TABLE IF NOT EXISTS results (
  result_id BIGSERIAL PRIMARY KEY,
  competition_id TEXT REFERENCES competitions (competition_id) NOT NULL,
  user_id INTEGER REFERENCES users (user_id) NOT NULL,
  event_id INTEGER REFERENCES events (event_id) NOT NULL,
  solve1 TEXT NOT NULL,
  solve2 TEXT NOT NULL,
  solve3 TEXT NOT NULL,
  solve4 TEXT NOT NULL,
  solve5 TEXT NOT NULL,
  comment TEXT NOT NULL,
  status_id INTEGER REFERENCES results_status (results_status_id) NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
