package app

import (
	"github.com/Kyrbanali/API-GateWay/internal/handler"
	"github.com/Kyrbanali/API-GateWay/internal/storage"
	"github.com/pkg/errors"
)

func Run() error {
	dsn := "postgres://postgres:postgres@postgres:5432/postgres"
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
