package usecase

import (
	"context"

	"github.com/Kyrbanali/API-GateWay/internal/models"
	// "github.com/Kyrbanali/API-GateWay/internal/repository"
)

type UserProvider interface {
	// New(userRepo *repository.UserRepo) *UseCase
	CreateUser(ctx context.Context, user models.UserDTO) (string, error)
	GetUserByID(ctx context.Context, id string) (*models.UserDTO, error)
	GetAllUsers(ctx context.Context) ([]models.UserDTO, error)
	DeleteUserByID(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, user models.UserDTO) error
}
