BEGIN;

CREATE TABLE IF NOT EXISTS wca_competitions_announcements_subscriptions(
  wca_competitions_announcements_subscription_id BIGSERIAL PRIMARY KEY,
  country_id TEXT REFERENCES countries (country_id) NOT NULL,
  user_id INTEGER REFERENCES users (user_id) NOT NULL, 
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (country_id, user_id)
);

COMMIT;
