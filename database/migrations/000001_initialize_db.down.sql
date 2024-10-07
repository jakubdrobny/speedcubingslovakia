BEGIN;

DROP TABLE IF EXISTS continents;
DROP TABLE IF EXISTS countries;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS competitions;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS competition_events;
DROP TABLE IF EXISTS results_status;
DROP TABLE IF EXISTS results;
DROP TABLE IF EXISTS scrambles;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS announcements;
DROP TABLE IF EXISTS announcement_tags;
DROP TABLE IF EXISTS announcement_read;
DROP TABLE IF EXISTS announcement_reaction;

COMMIT;
