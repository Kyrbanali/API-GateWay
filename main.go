package main

import (
	"log/slog"
	"os"

	"github.com/Kyrbanali/API-GateWay/internal/app"
	"github.com/Kyrbanali/API-GateWay/internal/storage"
)

func main() {
	db, err := storage.GetConnect()
	if err != nil {
		slog.Error("failed to connect to DB", slog.Any("error", err))
		os.Exit(1)
	}

	storage.Conn = db

	router := app.GetRouter()
	if err := router.Listen(":3000"); err != nil {
		slog.Error("app run", slog.Any("error", err))
		os.Exit(1)
	}
}
