-- name: GetExpenseById :one
SELECT * FROM expenses 
WHERE id=$1 LIMIT 1;

-- name: GetAllExpenses :many
SELECT * from expenses;

-- name: GetTotalSpendThisMonth :one
SELECT COALESCE(SUM(amount), 0)::DECIMAL(10,3) AS total
FROM expenses
WHERE is_deleted = false
  AND created_at >= date_trunc('month', CURRENT_TIMESTAMP);

-- name: GetDailySpend :many
SELECT created_at::date AS day, COALESCE(SUM(amount), 0)::DECIMAL(10,3) AS total
FROM expenses
WHERE is_deleted = false
  AND created_at >= (date_trunc('month', CURRENT_TIMESTAMP) - (sqlc.arg(months)::int || ' months')::interval)
GROUP BY created_at::date
ORDER BY day;

-- name: GetExpensesByMonth :many
SELECT id, name, amount, cost, item_id, is_deleted, created_at, updated_at
FROM expenses
WHERE is_deleted = false
  AND EXTRACT(YEAR FROM created_at) = sqlc.arg(year)::int
  AND EXTRACT(MONTH FROM created_at) = sqlc.arg(month)::int
ORDER BY created_at;

-- name: CreateExpense :one
INSERT INTO expenses (name, amount, cost)
VALUES ($1, $2, $3)
RETURNING id, name, amount, cost, item_id, is_deleted, created_at, updated_at;

-- name: GetTagsForExpense :many
SELECT t.id, t.name, t.created_at, t.updated_at, t.is_deleted
FROM tags t
JOIN expenses_tags et ON t.id = et.tag_id
WHERE et.expense_id = $1 AND t.is_deleted = false
ORDER BY t.name;

-- name: GetAllTags :many
SELECT id, name, created_at, updated_at, is_deleted
FROM tags
WHERE is_deleted = false
ORDER BY name;

-- name: CreateTag :one
INSERT INTO tags (name)
VALUES ($1)
RETURNING id, name, created_at, updated_at, is_deleted;

-- name: AddTagToExpense :exec
INSERT INTO expenses_tags (expense_id, tag_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveTagFromExpense :exec
DELETE FROM expenses_tags
WHERE expense_id = $1 AND tag_id = $2;

-- name: GetFirstExpenseDate :one
SELECT COALESCE(TO_CHAR(MIN(created_at::date), 'YYYY-MM-DD'), '') AS first_date
FROM expenses
WHERE is_deleted = false;

-- name: GetDailySpendForMonth :many
SELECT created_at::date AS day, COALESCE(SUM(amount), 0)::DECIMAL(10,3) AS total
FROM expenses
WHERE is_deleted = false
  AND EXTRACT(YEAR FROM created_at) = sqlc.arg(year)::int
  AND EXTRACT(MONTH FROM created_at) = sqlc.arg(month)::int
GROUP BY created_at::date
ORDER BY day;

-- name: GetUserByUsername :one
SELECT id, username, password, avatar_url, first_name, last_name, email, created_at, updated_at FROM users WHERE username = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (username, password, first_name, last_name, email) VALUES ($1, $2, $3, $4, $5)
RETURNING id, username, password, avatar_url, first_name, last_name, email, created_at, updated_at;