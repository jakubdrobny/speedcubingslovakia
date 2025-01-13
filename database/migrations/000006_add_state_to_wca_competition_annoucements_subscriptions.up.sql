BEGIN;

ALTER TABLE wca_competitions_announcements_subscriptions ADD COLUMN state TEXT DEFAULT '' NOT NULL;

COMMIT;
