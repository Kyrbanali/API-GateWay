package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetConnect(connStr string) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), connStr)
}
