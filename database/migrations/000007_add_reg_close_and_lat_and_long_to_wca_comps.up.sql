BEGIN;

ALTER TABLE upcoming_wca_competitions 
  ADD COLUMN registration_close TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  ADD COLUMN latitude_degrees   NUMERIC   DEFAULT 0.0               NOT NULL,
  ADD COLUMN longitude_degrees  NUMERIC   DEFAULT 0.0               NOT NULL;

COMMIT;
