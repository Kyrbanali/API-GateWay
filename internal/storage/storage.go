package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

func GetConnect(connStr string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect pgx pool")
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, errors.Wrap(err, "failed to ping to DB")
	}

	return pool, nil
}
