-- migrate:up
ALTER TABLE expenses ADD COLUMN user_id BIGINT REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE tags ADD COLUMN user_id BIGINT REFERENCES users(id) ON DELETE CASCADE;

-- Backfilling pre-existing rows is only safe when ownership is unambiguous,
-- i.e. there is exactly one user. With multiple users we cannot infer who owns
-- each legacy expense/tag, and blindly assigning them all to the oldest user
-- would leak every record into one account. Abort instead and require a manual
-- backfill in that case.
DO $$
DECLARE
    user_count BIGINT;
    orphan_count BIGINT;
BEGIN
    SELECT count(*) INTO user_count FROM users;
    SELECT count(*) INTO orphan_count
    FROM (
        SELECT 1 FROM expenses WHERE user_id IS NULL
        UNION ALL
        SELECT 1 FROM tags WHERE user_id IS NULL
    ) o;

    IF orphan_count > 0 AND user_count <> 1 THEN
        RAISE EXCEPTION
            'Cannot backfill user_id: % users exist with % orphan expense/tag rows; owner is ambiguous. Assign user_id manually, then re-run this migration.',
            user_count, orphan_count;
    END IF;
END $$;

UPDATE expenses SET user_id = (SELECT id FROM users ORDER BY id LIMIT 1) WHERE user_id IS NULL;
UPDATE tags SET user_id = (SELECT id FROM users ORDER BY id LIMIT 1) WHERE user_id IS NULL;

DELETE FROM expenses WHERE user_id IS NULL;
DELETE FROM tags WHERE user_id IS NULL;

-- Deduplicate tags per (user_id, name). Repoint expenses_tags to surviving tag id, drop dupes.
WITH ranked AS (
    SELECT id,
           user_id,
           name,
           MIN(id) OVER (PARTITION BY user_id, name) AS keeper_id
    FROM tags
),
remap AS (
    SELECT id AS dup_id, keeper_id
    FROM ranked
    WHERE id <> keeper_id
)
UPDATE expenses_tags et
SET tag_id = r.keeper_id
FROM remap r
WHERE et.tag_id = r.dup_id
  AND NOT EXISTS (
      SELECT 1 FROM expenses_tags et2
      WHERE et2.expense_id = et.expense_id AND et2.tag_id = r.keeper_id
  );

WITH dupes AS (
    SELECT id FROM (
        SELECT id, ROW_NUMBER() OVER (PARTITION BY user_id, name ORDER BY id) AS rn
        FROM tags
    ) ranked
    WHERE rn > 1
)
DELETE FROM expenses_tags WHERE tag_id IN (SELECT id FROM dupes);

WITH dupes AS (
    SELECT id FROM (
        SELECT id, ROW_NUMBER() OVER (PARTITION BY user_id, name ORDER BY id) AS rn
        FROM tags
    ) ranked
    WHERE rn > 1
)
DELETE FROM tags WHERE id IN (SELECT id FROM dupes);

ALTER TABLE expenses ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE tags ALTER COLUMN user_id SET NOT NULL;

CREATE INDEX idx_expenses_user_id ON expenses(user_id);
CREATE INDEX idx_tags_user_id ON tags(user_id);

ALTER TABLE tags DROP CONSTRAINT IF EXISTS tags_name_key;
CREATE UNIQUE INDEX IF NOT EXISTS idx_tags_user_name ON tags(user_id, name);

-- migrate:down
DROP INDEX IF EXISTS idx_tags_user_name;
DROP INDEX IF EXISTS idx_expenses_user_id;
DROP INDEX IF EXISTS idx_tags_user_id;
ALTER TABLE expenses DROP COLUMN IF EXISTS user_id;
ALTER TABLE tags DROP COLUMN IF EXISTS user_id;
