package psql

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabaseConnection() (*pgxpool.Pool, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}
	connPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return connPool, nil
}
