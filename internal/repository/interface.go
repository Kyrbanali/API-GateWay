package repository

import (
	"context"

	"github.com/Kyrbanali/API-GateWay/internal/models"
	// "github.com/jackc/pgx/v5/pgxpool"
)

type UserProvider interface {
	// New(conn *pgxpool.Pool) *UserRepo
	CreateUser(ctx context.Context, user models.UserDTO) (string, error)
	GetUserByID(ctx context.Context, id string) (*models.UserDTO, error)
	GetAllUsers(ctx context.Context) ([]models.UserDTO, error)
	DeleteUserByID(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, user models.UserDTO) error
}
