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

func (h *Handle) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid uuid",
		})
	}

	row := h.conn.QueryRow(c.Context(), `
		SELECT id, name, age FROM users WHERE id = $1
	`, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Age); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.JSON(user)
}

func (h *Handle) GetAllUsers(c *fiber.Ctx) error {
	rows, err := h.conn.Query(c.Context(), "SELECT id, name, age FROM users")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		user := models.User{}
		rows.Scan(&user.ID, &user.Name, &user.Age)
		users = append(users, user)
	}
	return c.JSON(users)
}

func (h *Handle) DeleteUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid uuid",
		})
	}

	res, err := h.conn.Exec(c.Context(), `
		DELETE FROM users WHERE id = $1
	`, id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if res.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "user deleted",
	})
}

func (h *Handle) UpdateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse Json",
		})
	}

	if _, err := uuid.Parse(user.ID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid UUID",
		})
	}

	res, err := h.conn.Exec(c.Context(), `
		UPDATE users SET name = $1, age = $2 WHERE id = $3
	`, user.Name, user.Age, user.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if res.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.JSON(fiber.Map{
		"id": user.ID,
	})
}
