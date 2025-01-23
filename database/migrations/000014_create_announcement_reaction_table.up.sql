BEGIN;

CREATE TABLE IF NOT EXISTS announcement_reaction (
    announcement_reaction_id BIGSERIAL PRIMARY KEY,
    announcement_id INTEGER REFERENCES announcements (announcement_id) ON DELETE CASCADE NOT NULL,
    user_id INTEGER REFERENCES users (user_id) ON DELETE CASCADE NOT NULl,
    emoji TEXT NOT NULL,
    "by" TEXT NOT NULL,
    "set" BOOLEAN NOT NULL
);

COMMIT;
