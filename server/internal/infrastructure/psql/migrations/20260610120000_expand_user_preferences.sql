-- migrate:up
ALTER TABLE user_preferences
    ADD COLUMN week_start TEXT NOT NULL DEFAULT 'Monday',
    ADD COLUMN monthly_budget DECIMAL(10,2) NOT NULL DEFAULT 0,
    ADD COLUMN notify_weekly_summary BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN notify_budget_alerts BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN notify_product_updates BOOLEAN NOT NULL DEFAULT FALSE;

-- migrate:down
ALTER TABLE user_preferences
    DROP COLUMN week_start,
    DROP COLUMN monthly_budget,
    DROP COLUMN notify_weekly_summary,
    DROP COLUMN notify_budget_alerts,
    DROP COLUMN notify_product_updates;
