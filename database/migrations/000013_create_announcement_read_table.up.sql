BEGIN;

CREATE TABLE IF NOT EXISTS announcement_read (
    announcement_read_id BIGSERIAL PRIMARY KEY,
    announcement_id INTEGER REFERENCES announcements (announcement_id) NOT NULL,
    user_id INTEGER REFERENCES users (user_id) NOT NULL,
    read BOOLEAN NOT NULL DEFAULT FALSE,
    read_timestamp TIMESTAMP,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
