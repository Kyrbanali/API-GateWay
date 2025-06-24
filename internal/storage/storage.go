package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Conn *pgxpool.Pool

func GetConnect() (*pgxpool.Pool, error) {
	connStr := "postgres://postgres:postgres@postgres:5432/postgres"

	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return db, nil
}
