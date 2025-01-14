BEGIN;

ALTER TABLE wca_competitions_announcements_subscriptions DROP CONSTRAINT wca_competitions_announcements_subscript_country_id_user_id_key;
ALTER TABLE wca_competitions_announcements_subscriptions ADD CONSTRAINT wca_competitions_announcements_subscript_country_id_state_user_id_key UNIQUE (country_id, user_id);

COMMIT;
