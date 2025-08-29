package psql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func NewDatabaseConnection() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:secretpass@localhost:5432/jafa")
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())
	return conn, nil

}
