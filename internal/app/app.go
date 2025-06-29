package app

import (
	"fmt"

	"github.com/Kyrbanali/API-GateWay/config"
	"github.com/Kyrbanali/API-GateWay/internal/handler"
	"github.com/Kyrbanali/API-GateWay/internal/repository"
	"github.com/Kyrbanali/API-GateWay/internal/storage"
	"github.com/Kyrbanali/API-GateWay/internal/usecase"
	"github.com/pkg/errors"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return errors.Wrap(err, "config load failed")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Postgres.User,
		cfg.Postgres.Pass,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DB,
	)

	conn, err := storage.GetConnect(dsn)
	if err != nil {
		return errors.Wrap(err, "failed to connect to DB")
	}
	defer conn.Close()

	repo := repository.New(conn)
	uc := usecase.New(repo)
	h := handler.New(uc)

	router := GetRouter(h)

	if err := router.Listen(":3000"); err != nil {
		return errors.Wrap(err, "app listen")
	}

	return nil
}
