package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Config struct {
	Host string
	User string
	Pass string
	DB   string
	Port string
}

func (c *Config) BuildDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.User,
		c.Pass,
		c.Host,
		c.Port,
		c.DB,
	)
}

func GetConnect(connStr string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, errors.Wrap(err, "connect pgx pool")
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, errors.Wrap(err, "ping to DB")
	}

	return pool, nil
}
