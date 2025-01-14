BEGIN;

ALTER TABLE upcoming_wca_competitions DROP COLUMN registration_close, DROP COLUMN latitude_degrees, DROP COLUMN longitude_degrees;

COMMIT;
