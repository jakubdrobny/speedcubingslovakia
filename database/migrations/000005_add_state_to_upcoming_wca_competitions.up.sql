BEGIN;

ALTER TABLE upcoming_wca_competitions ADD COLUMN state TEXT DEFAULT '' NOT NULL;

COMMIT;
