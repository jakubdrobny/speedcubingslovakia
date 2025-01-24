BEGIN;

CREATE TABLE IF NOT EXISTS announcements (
    announcement_id BIGSERIAL PRIMARY KEY,
    author_id INTEGER REFERENCES users (user_id) NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
