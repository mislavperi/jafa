package psql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Pool is the database handle a service is constructed with. It can run queries
// directly (DBTX) and open transactions (Begin). Both *pgxpool.Pool and the
// test mock satisfy it, so services depend on this interface rather than the
// concrete pool.
type Pool interface {
	DBTX
	Begin(ctx context.Context) (pgx.Tx, error)
}

// RunInTx runs fn inside a single transaction and returns its result. The
// transaction is committed when fn succeeds and rolled back on any error (or
// panic), so a multi-statement operation either fully applies or leaves no
// trace. fn must use the transaction-bound *Queries it is given, not the
// service's pool-bound queries, or its writes will not be part of the tx.
func RunInTx[T any](ctx context.Context, pool Pool, fn func(q *Queries) (T, error)) (T, error) {
	var zero T
	tx, err := pool.Begin(ctx)
	if err != nil {
		return zero, err
	}
	defer tx.Rollback(ctx)

	result, err := fn(New(tx))
	if err != nil {
		return zero, err
	}
	if err := tx.Commit(ctx); err != nil {
		return zero, err
	}
	return result, nil
}
