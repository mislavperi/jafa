-- name: GetExpenseById :one
SELECT * FROM expenses 
WHERE id=$1 LIMIT 1;