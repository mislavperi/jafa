-- migrate:up
CREATE TABLE user_preferences (
    user_id    BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    accent_id  TEXT NOT NULL DEFAULT 'amber',
    font_size  TEXT NOT NULL DEFAULT 'normal',
    dark_mode  BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- migrate:down
DROP TABLE IF EXISTS user_preferences;
