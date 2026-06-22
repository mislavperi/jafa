-- migrate:up
ALTER TABLE expenses
ADD COLUMN installment_count INTEGER;

-- migrate:down
ALTER TABLE expenses
DROP COLUMN IF EXISTS installment_count;
