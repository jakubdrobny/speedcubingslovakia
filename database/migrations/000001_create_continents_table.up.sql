BEGIN;

CREATE TABLE IF NOT EXISTS continents (
    continent_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    recordName TEXT NOT NULL,
    CONSTRAINT continents_unique UNIQUE (continent_id, name, recordName)
);

COMMIT;
