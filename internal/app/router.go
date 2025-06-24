package app

import (
	"github.com/Kyrbanali/API-GateWay/internal/models"
	"github.com/Kyrbanali/API-GateWay/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func GetRouter() *fiber.App {
	app := fiber.New()

	app.Post("/user", createUser)
	app.Get("/user/:id", getUserById)
	// app.Delete("/user/:id", deleteUserById)
	// app.Put("/user", updateUser)

	return app
}

func createUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse json",
		})
	}

	user.ID = uuid.NewString()

	if err := storage.InsertUser(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error inserting into DB",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": user.ID,
	})
}

func getUserById(c *fiber.Ctx) error {
	id := c.Params("id")

	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid uuid",
		})
	}

	user, err := storage.GetUserById(c.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error getting user",
		})
	}

	return c.JSON(user)
}

// func deleteUserById(c *fiber.Ctx) error {
// 	id := c.Params("id")

// 	if _, err := uuid.Parse(id); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "invalid uuid",
// 		})
// 	}

// 	_, exists := users[id]
// 	if !exists {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"error": "user not found",
// 		})
// 	}

// 	delete(users, id)

// 	return c.JSON(fiber.Map{
// 		"message": "user deleted",
// 	})
// }

// func updateUser(c *fiber.Ctx) error {
// 	var updatedUser User

// 	if err := c.BodyParser(&updatedUser); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "failed to parse Json",
// 		})
// 	}

// 	_, esists := users[updatedUser.ID]
// 	if !esists {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"error": "user not found",
// 		})
// 	}

// 	users[updatedUser.ID] = updatedUser

// 	return c.JSON(fiber.Map{
// 		"id": updatedUser.ID,
// 	})

// }
