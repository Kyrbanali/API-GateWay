package app

import (
	"time"

	"github.com/Kyrbanali/API-GateWay/config"
	"github.com/Kyrbanali/API-GateWay/database"
	"github.com/Kyrbanali/API-GateWay/internal/cache"
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

	if err := database.Migrate(cfg.Postgres.BuildDSN()); err != nil {
		return errors.Wrap(err, "migration")
	}

	conn, err := storage.GetConnect(cfg.Postgres.BuildDSN())
	if err != nil {
		return errors.Wrap(err, "failed to connect to DB")
	}
	defer conn.Close()

	repo := repository.New(conn)

	ttl, err := time.ParseDuration(cfg.Cache.TTL)
	if err != nil {
		return errors.Wrap(err, "invalid cache ttl")
	}

	cleanupInterval, err := time.ParseDuration(cfg.Cache.CleanupInterval)
	if err != nil {
		return errors.Wrap(err, "invalid cache cleanupInterval")
	}

	cacheDecorator := cache.New(repo, ttl, cleanupInterval)
	uc := usecase.New(cacheDecorator)
	handle := handler.New(uc)

	router := GetRouter(handle)

	if err := router.Listen(":" + cfg.App.Port); err != nil {
		return errors.Wrap(err, "app listen")
	}

	return nil
}
