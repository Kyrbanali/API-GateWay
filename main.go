package main

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

func main() {
	app := fiber.New()

	app.Post("/user", createUser)
	app.Get("/user/:id", getUserById)

	app.Listen(":3000")
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
