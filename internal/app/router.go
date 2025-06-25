package app

import (
	"errors"

	"github.com/Kyrbanali/API-GateWay/internal/handler"
	"github.com/Kyrbanali/API-GateWay/internal/models"
	"github.com/Kyrbanali/API-GateWay/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func GetRouter(h *handler.Handle) *fiber.App {
	app := fiber.New()

	app.Post("/user", h.CreateUser)
	app.Get("/user/:id", h.GetUserByID)
	app.Get("/users", getAllUsers)
	app.Delete("/user/:id", deleteUserById)
	app.Put("/user", updateUser)

	return app
}

func getAllUsers(c *fiber.Ctx) error {
	users, err := storage.GetAllUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(users)
}

func deleteUserById(c *fiber.Ctx) error {
	id := c.Params("id")

	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid uuid",
		})
	}

	_, exists := storage.GetUserByID(c.Context(), id)
	if exists != nil {
		if exists == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
	}

	if err := storage.DeleteUser(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error deleting user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "user deleted",
	})
}

func updateUser(c *fiber.Ctx) error {
	var updatedUser models.User

	if err := c.BodyParser(&updatedUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse Json",
		})
	}

	if _, err := uuid.Parse(updatedUser.ID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid UUID",
		})
	}

	_, err := storage.GetUserByID(c.Context(), updatedUser.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "error checking user existence",
			"details": err.Error(),
		})
	}

	if err := storage.UpdateUser(c.Context(), updatedUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "error updating user",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"id": updatedUser.ID,
	})

}
