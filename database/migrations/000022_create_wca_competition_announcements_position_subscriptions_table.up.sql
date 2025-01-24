BEGIN;

CREATE TABLE IF NOT EXISTS wca_competition_announcements_position_subscriptions (
  wca_competition_announcements_position_subscription_id BIGSERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users (user_id) NOT NULL,
  latitude_degrees NUMERIC DEFAULT 0.0 NOT NULL,
  longitude_degrees NUMERIC DEFAULT 0.0 NOT NULL,
  radius INTEGER DEFAULT 50 NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT user_location_unique UNIQUE(user_id, latitude_degrees, longitude_degrees, radius),
  CONSTRAINT radius_earth CHECK (radius <= 100000)
);

COMMIT;
