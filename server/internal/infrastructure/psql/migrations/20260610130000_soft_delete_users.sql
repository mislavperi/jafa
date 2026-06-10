-- migrate:up
ALTER TABLE users
    ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT FALSE;

-- Free a soft-deleted account's username for re-registration: uniqueness only
-- applies to active users.
ALTER TABLE users DROP CONSTRAINT users_username_key;
CREATE UNIQUE INDEX idx_users_username_active ON users(username) WHERE is_deleted = false;

-- migrate:down
-- Soft-deleted rows must go before the full unique constraint can be restored,
-- mirroring the hard-delete semantics this migration replaced.
DELETE FROM users WHERE is_deleted = true;
DROP INDEX idx_users_username_active;
ALTER TABLE users ADD CONSTRAINT users_username_key UNIQUE (username);
ALTER TABLE users DROP COLUMN is_deleted;
