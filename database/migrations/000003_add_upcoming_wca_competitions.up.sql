BEGIN;

CREATE TABLE IF NOT EXISTS upcoming_wca_competitions(
  upcoming_wca_competition_id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  startdate TIMESTAMP NOT NULL,
  enddate TIMESTAMP NOT NULL,
  registered INTEGER NOT NULL,
  competitor_limit INTEGER NOT NULL,
  venue_address TEXT NOT NULL,
  url TEXT NOT NULL,
  country_id TEXT REFERENCES countries (country_id) NOT NULL,
  registration_open TIMESTAMP NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS upcoming_wca_competition_events(
  upcoming_wca_competition_event_id BIGSERIAL PRIMARY KEY,
  upcoming_wca_competition_id TEXT REFERENCES upcoming_wca_competitions (upcoming_wca_competition_id) ON DELETE CASCADE NOT NULL,
  event_id INTEGER REFERENCES events (event_id) NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (upcoming_wca_competition_id, event_id)
);

COMMIT;
