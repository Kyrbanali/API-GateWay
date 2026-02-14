package app

import (
	"github.com/Kyrbanali/API-GateWay/config"
	"github.com/Kyrbanali/API-GateWay/database"
	"github.com/Kyrbanali/API-GateWay/internal/cache"
	"github.com/Kyrbanali/API-GateWay/internal/handler"
	"github.com/Kyrbanali/API-GateWay/internal/link"
	"github.com/Kyrbanali/API-GateWay/internal/repository"
	"github.com/Kyrbanali/API-GateWay/internal/storage"
	"github.com/Kyrbanali/API-GateWay/internal/usecase"
	"github.com/Kyrbanali/API-GateWay/internal/worker"
	"github.com/pkg/errors"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return errors.Wrap(err, "config load")
	}

	if err := database.Migrate(cfg.Postgres.BuildDSN()); err != nil {
		return errors.Wrap(err, "migration")
	}

	conn, err := storage.GetConnect(cfg.Postgres.BuildDSN())
	if err != nil {
		return errors.Wrap(err, "connect to DB")
	}
	defer conn.Close()

	repo := repository.New(conn)
	cacheDecorator := cache.New(repo, cfg.Cache)

	work := worker.New(cfg.Workers)
	fetcher := link.New()

	uc := usecase.New(cacheDecorator, work, fetcher)
	handle := handler.New(uc)
	router := GetRouter(handle)

	if err := router.Listen(":" + cfg.App.Port); err != nil {
		return errors.Wrap(err, "app listen")
	}

	return nil
}
