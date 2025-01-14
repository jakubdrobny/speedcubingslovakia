BEGIN;

ALTER TABLE upcoming_wca_competitions DROP COLUMN IF EXISTS state;

COMMIT;
