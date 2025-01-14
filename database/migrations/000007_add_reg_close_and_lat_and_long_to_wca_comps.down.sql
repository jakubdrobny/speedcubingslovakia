BEGIN;

ALTER TABLE upcoming_wca_competitions DROP COLUMN IF EXISTS registration_close, DROP COLUMN IF EXISTS latitude_degrees, DROP COLUMN IF EXISTS longitude_degrees;

COMMIT;
