-- name: GetExpenseById :one
SELECT * FROM expenses
WHERE id=$1 AND user_id=$2 AND is_deleted = false
LIMIT 1;

-- name: GetAllExpenses :many
SELECT * FROM expenses
WHERE user_id=$1 AND is_deleted = false;

-- name: GetTotalSpendThisMonth :one
SELECT COALESCE(SUM(amount), 0)::DECIMAL(10,3) AS total
FROM expenses
WHERE is_deleted = false
  AND user_id = $1
  AND created_at >= date_trunc('month', CURRENT_TIMESTAMP);

-- name: GetDailySpend :many
SELECT created_at::date AS day, COALESCE(SUM(amount), 0)::DECIMAL(10,3) AS total
FROM expenses
WHERE is_deleted = false
  AND user_id = sqlc.arg(user_id)::bigint
  AND created_at >= (date_trunc('month', CURRENT_TIMESTAMP) - (sqlc.arg(months)::int || ' months')::interval)
GROUP BY created_at::date
ORDER BY day;

-- name: GetExpensesByMonth :many
SELECT *
FROM expenses
WHERE is_deleted = false
  AND user_id = sqlc.arg(user_id)::bigint
  AND EXTRACT(YEAR FROM created_at) = sqlc.arg(year)::int
  AND EXTRACT(MONTH FROM created_at) = sqlc.arg(month)::int
ORDER BY created_at;

-- name: CreateExpense :one
INSERT INTO expenses (name, amount, cost, recurrence_interval, recurrence_day, recurrence_start_date, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateExpense :one
UPDATE expenses
SET name = $3,
    amount = $4,
    cost = $5,
    recurrence_interval = $6,
    recurrence_day = $7,
    recurrence_start_date = $8,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $2 AND is_deleted = false
RETURNING *;

-- name: SoftDeleteExpense :exec
UPDATE expenses
SET is_deleted = true,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $2;

-- name: GetTagsForExpense :many
SELECT t.*
FROM tags t
JOIN expenses_tags et ON t.id = et.tag_id
JOIN expenses e ON e.id = et.expense_id
WHERE et.expense_id = $1
  AND e.user_id = $2
  AND t.user_id = $2
  AND t.is_deleted = false
ORDER BY t.name;

-- name: GetAllTags :many
SELECT *
FROM tags
WHERE is_deleted = false AND user_id = $1
ORDER BY name;

-- name: CreateTag :one
INSERT INTO tags (name, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: AddTagToExpense :exec
INSERT INTO expenses_tags (expense_id, tag_id)
SELECT sqlc.arg(expense_id)::bigint, sqlc.arg(tag_id)::bigint
WHERE EXISTS (SELECT 1 FROM expenses WHERE id = sqlc.arg(expense_id)::bigint AND user_id = sqlc.arg(user_id)::bigint AND is_deleted = false)
  AND EXISTS (SELECT 1 FROM tags WHERE id = sqlc.arg(tag_id)::bigint AND user_id = sqlc.arg(user_id)::bigint AND is_deleted = false)
ON CONFLICT DO NOTHING;

-- name: RemoveTagFromExpense :exec
DELETE FROM expenses_tags
WHERE expense_id = sqlc.arg(expense_id)::bigint AND tag_id = sqlc.arg(tag_id)::bigint
  AND EXISTS (SELECT 1 FROM expenses WHERE id = sqlc.arg(expense_id)::bigint AND user_id = sqlc.arg(user_id)::bigint);

-- name: GetFirstExpenseDate :one
SELECT COALESCE(TO_CHAR(MIN(created_at::date), 'YYYY-MM-DD'), '') AS first_date
FROM expenses
WHERE is_deleted = false AND user_id = $1;

-- name: GetDailySpendForMonth :many
SELECT created_at::date AS day, COALESCE(SUM(amount), 0)::DECIMAL(10,3) AS total
FROM expenses
WHERE is_deleted = false
  AND user_id = sqlc.arg(user_id)::bigint
  AND EXTRACT(YEAR FROM created_at) = sqlc.arg(year)::int
  AND EXTRACT(MONTH FROM created_at) = sqlc.arg(month)::int
GROUP BY created_at::date
ORDER BY day;

-- name: GetUserPreferences :one
SELECT * FROM user_preferences WHERE user_id = $1 LIMIT 1;

-- name: UpsertUserPreferences :one
INSERT INTO user_preferences (user_id, accent_id, font_size, dark_mode, currency)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id) DO UPDATE
  SET accent_id = EXCLUDED.accent_id,
      font_size = EXCLUDED.font_size,
      dark_mode = EXCLUDED.dark_mode,
      currency = EXCLUDED.currency,
      updated_at = NOW()
RETURNING *;

-- name: GetUserByUsername :one
SELECT id, username, password, avatar_url, first_name, last_name, email, created_at, updated_at FROM users WHERE username = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (username, password, first_name, last_name, email) VALUES ($1, $2, $3, $4, $5)
RETURNING id, username, password, avatar_url, first_name, last_name, email, created_at, updated_at;

-- name: ListCategories :many
SELECT * FROM categories ORDER BY sort_order;

-- name: GetMonthlySpend :many
SELECT to_char(date_trunc('month', created_at), 'YYYY-MM') AS month,
       COALESCE(SUM(amount), 0)::DECIMAL(10,3) AS total
FROM expenses
WHERE is_deleted = false AND user_id = $1
GROUP BY date_trunc('month', created_at)
ORDER BY date_trunc('month', created_at);
