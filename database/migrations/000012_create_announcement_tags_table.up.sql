BEGIN;

CREATE TABLE IF NOT EXISTS announcement_tags (
    announcement_tags_id BIGSERIAL PRIMARY KEY,
    announcement_id INTEGER REFERENCES announcements (announcement_id) NOT NULL,
    tag_id INTEGER REFERENCES tags (tag_id) NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
