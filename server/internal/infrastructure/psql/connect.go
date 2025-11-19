package psql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabaseConnection() (*pgxpool.Pool, error) {
	connPool, err := pgxpool.New(context.Background(), "postgres://postgres:secretpass@localhost:5432/jafa")
	if err != nil {
		return nil, err
	}
	return connPool, nil
}
