package app

import (
	middleware "github.com/Kyrbanali/API-GateWay/internal/app/middlware"
	"github.com/Kyrbanali/API-GateWay/internal/handler"
	"github.com/Kyrbanali/API-GateWay/internal/metrics"
	"github.com/gofiber/fiber/v2"
)

func GetRouter(h *handler.Handle) *fiber.App {
	app := fiber.New()

	app.Use(middleware.Metrics())
	app.Get("/metrics", metrics.Handler())

	app.Post("/user", h.CreateUser)
	app.Get("/user/:id", h.GetUserByID)
	app.Get("/users", h.GetAllUsers)
	app.Delete("/user/:id", h.DeleteUserByID)
	app.Put("/user", h.UpdateUser)

	return app
}
