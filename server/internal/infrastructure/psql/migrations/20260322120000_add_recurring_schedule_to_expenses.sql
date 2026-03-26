-- migrate:up
ALTER TABLE expenses
ADD COLUMN recurrence_interval TEXT,
ADD COLUMN recurrence_day INTEGER,
ADD COLUMN recurrence_start_date DATE;

-- migrate:down
ALTER TABLE expenses
DROP COLUMN IF EXISTS recurrence_interval,
DROP COLUMN IF EXISTS recurrence_day,
DROP COLUMN IF EXISTS recurrence_start_date;
