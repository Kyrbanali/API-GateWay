package handler

import (
	"strings"

	"github.com/Kyrbanali/API-GateWay/internal/models"
	"github.com/Kyrbanali/API-GateWay/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handle struct {
	uc *usecase.UseCase
}

func New(uc *usecase.UseCase) *Handle {
	return &Handle{uc: uc}
}

func (h *Handle) CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse json",
		})
	}

	id, err := h.uc.CreateUser(c.Context(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}

func (h *Handle) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid uuid",
		})
	}

	user, err := h.uc.GetUserByID(c.Context(), id)
	if err != nil {
		if err.Error() == "no rows in result set" || strings.Contains(err.Error(), "no rows") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "error querying user",
			"details": err.Error(),
		})
	}

	return c.JSON(user)
}

func (h *Handle) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.uc.GetAllUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to fetch users",
			"details": err.Error(),
		})
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

	err := h.uc.DeleteUserByID(c.Context(), id)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to delete user",
			"details": err.Error(),
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

	err := h.uc.UpdateUser(c.Context(), user)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to update user",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"id": user.ID,
	})
}
