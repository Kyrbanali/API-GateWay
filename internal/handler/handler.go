package handler

import (
	"context"

	"github.com/Kyrbanali/API-GateWay/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handle struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *Handle {
	return &Handle{conn: conn}
}

func (h *Handle) CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse json",
		})
	}

	user.ID = uuid.NewString()

	_, err := h.conn.Exec(context.Background(), `
		INSERT INTO USERS (id, name, age) VALUES ($1, $2, $3)
	`, user.ID, user.Name, user.Age)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "error inserting user",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": user.ID,
	})
}
