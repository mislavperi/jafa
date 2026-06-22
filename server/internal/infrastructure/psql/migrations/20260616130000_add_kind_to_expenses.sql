-- migrate:up
ALTER TABLE expenses
ADD COLUMN kind TEXT NOT NULL DEFAULT 'expense';

ALTER TABLE expenses
ADD CONSTRAINT expenses_kind_check CHECK (kind IN ('expense', 'income'));

-- migrate:down
ALTER TABLE expenses
DROP CONSTRAINT IF EXISTS expenses_kind_check;

ALTER TABLE expenses
DROP COLUMN IF EXISTS kind;
