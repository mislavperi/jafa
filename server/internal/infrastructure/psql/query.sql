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