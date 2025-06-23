package main

import (
	"log/slog"
	"os"

	"github.com/Kyrbanali/API-GateWay/internal/app"
)

func main() {
	appInstance := app.GetRouter()

	if err := appInstance.Listen(":3000"); err != nil {
		slog.Error("app run", slog.Any("error", err))
		os.Exit(1)
	}
}
