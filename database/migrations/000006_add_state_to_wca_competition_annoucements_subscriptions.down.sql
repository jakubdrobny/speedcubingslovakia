BEGIN;

ALTER TABLE wca_competitions_announcements_subscriptions DROP COLUMN state;

COMMIT;
