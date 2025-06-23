package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users = make(map[string]User)

func GetRouter() *fiber.App {
	app := fiber.New()

	app.Post("/user", createUser)
	app.Get("/user/:id", getUserById)
	app.Delete("/user/:id", deleteUserById)
	app.Put("/user", updateUser)

	return app
}

func createUser(c *fiber.Ctx) error {
	var user User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "не удалось распарсить json",
		})
	}

	user.ID = uuid.NewString()

	users[user.ID] = user

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": user.ID,
	})
}

func getUserById(c *fiber.Ctx) error {
	id := c.Params("id")

	user, exists := users[id]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Пользователь не найден",
		})
	}

	return c.JSON(user)
}

func deleteUserById(c *fiber.Ctx) error {
	id := c.Params("id")

	_, exists := users[id]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Пользователь не найден",
		})
	}

	delete(users, id)

	return c.JSON(fiber.Map{
		"message": "Пользователь удалён",
	})
}

func updateUser(c *fiber.Ctx) error {
	var updatedUser User

	if err := c.BodyParser(&updatedUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "не удалось распарсить Json",
		})
	}

	_, esists := users[updatedUser.ID]
	if !esists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Пользователь не найден",
		})
	}

	users[updatedUser.ID] = updatedUser

	return c.JSON(fiber.Map{
		"id": updatedUser.ID,
	})

}
