package app

import (
	"fmt"

	"github.com/Kyrbanali/API-GateWay/config"
	"github.com/Kyrbanali/API-GateWay/internal/handler"
	"github.com/Kyrbanali/API-GateWay/internal/storage"
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

	db, err := storage.GetConnect(dsn)
	if err != nil {
		return errors.Wrap(err, "failed to connect to DB")
	}
	defer db.Close()

	h := handler.New(db)
	router := GetRouter(h)

	if err := router.Listen(":3000"); err != nil {
		return errors.Wrap(err, "app listen")
	}

	return nil
}
