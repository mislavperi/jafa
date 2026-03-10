-- migrate:up
ALTER TABLE users ADD COLUMN avatar_url TEXT NOT NULL DEFAULT '';

-- migrate:down
ALTER TABLE users DROP COLUMN avatar_url;
