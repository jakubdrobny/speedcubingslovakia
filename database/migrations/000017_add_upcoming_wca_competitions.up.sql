BEGIN;

CREATE TABLE IF NOT EXISTS upcoming_wca_competitions(
  upcoming_wca_competition_id TEXT,
  name TEXT NOT NULL,
  startdate TIMESTAMP NOT NULL,
  enddate TIMESTAMP NOT NULL,
  registered INTEGER NOT NULL,
  competitor_limit INTEGER NOT NULL,
  venue_address TEXT NOT NULL,
  url TEXT NOT NULL,
  country_id TEXT REFERENCES countries (country_id) NOT NULL,
  registration_open TIMESTAMP NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (upcoming_wca_competition_id, country_id)
);

CREATE TABLE IF NOT EXISTS upcoming_wca_competition_events(
  upcoming_wca_competition_event_id BIGSERIAL PRIMARY KEY,
  upcoming_wca_competition_id TEXT NOT NULL,
  country_id TEXT NOT NULL,
  event_id INTEGER REFERENCES events (event_id) NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (upcoming_wca_competition_id, country_id, event_id),
  FOREIGN KEY (upcoming_wca_competition_id, country_id) REFERENCES upcoming_wca_competitions (upcoming_wca_competition_id, country_id) ON DELETE CASCADE
);

COMMIT;
