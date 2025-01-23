BEGIN;

CREATE TABLE IF NOT EXISTS events (
  event_id BIGSERIAL PRIMARY KEY,
  fulldisplayname TEXT NOT NULL,
  displayname TEXT NOT NULL,
  format TEXT NOT NULL,
  iconcode TEXT NOT NULL,
  scramblingcode TEXT NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT event_unique UNIQUE (fulldisplayname, displayname, format, iconcode, scramblingcode)
);

COMMIT;
