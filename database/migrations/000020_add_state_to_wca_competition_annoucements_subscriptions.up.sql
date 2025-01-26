BEGIN;

ALTER TABLE wca_competitions_announcements_subscriptions ADD COLUMN IF NOT EXISTS state TEXT DEFAULT '' NOT NULL;

COMMIT;
