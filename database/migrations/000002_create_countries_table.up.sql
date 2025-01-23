BEGIN;

CREATE TABLE IF NOT EXISTS countries (
  country_id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  continent_id TEXT REFERENCES continents (continent_id) NOT NULL,
  iso2 TEXT NOT NULL,
  CONSTRAINT countries_unique UNIQUE (country_id, name, continent_id, iso2)
);

COMMIT;
