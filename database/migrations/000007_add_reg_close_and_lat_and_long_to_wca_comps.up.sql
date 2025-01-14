BEGIN;

ALTER TABLE upcoming_wca_competitions 
  ADD COLUMN IF NOT EXISTS registration_close TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  ADD COLUMN IF NOT EXISTS latitude_degrees   NUMERIC   DEFAULT 0.0               NOT NULL,
  ADD COLUMN IF NOT EXISTS longitude_degrees  NUMERIC   DEFAULT 0.0               NOT NULL;

COMMIT;
