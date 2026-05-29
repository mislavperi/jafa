-- migrate:up
ALTER TABLE user_preferences
    ADD COLUMN currency TEXT NOT NULL DEFAULT 'EUR';

-- migrate:down
ALTER TABLE user_preferences
    DROP COLUMN currency;
