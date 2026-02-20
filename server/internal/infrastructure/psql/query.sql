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