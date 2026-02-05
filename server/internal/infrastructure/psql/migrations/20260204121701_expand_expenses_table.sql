-- migrate:up
ALTER TABLE expenses
    ADD COLUMN IF NOT EXISTS amount DECIMAL(10, 3) NOT NULL;

-- migrate:down
ALTER TABLE expenses
    DROP COLUMN IF EXISTS amount;
