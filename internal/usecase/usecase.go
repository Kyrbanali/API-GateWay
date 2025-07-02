package usecase

import (
	"context"

	"github.com/Kyrbanali/API-GateWay/internal/models"
	"github.com/Kyrbanali/API-GateWay/internal/repository"
)

type UseCase struct {
	userRepo repository.UserProvider
}

func New(userRepo *repository.UserRepo) *UseCase {
	return &UseCase{userRepo: userRepo}
}

func (u *UseCase) CreateUser(ctx context.Context, user models.UserDTO) (string, error) {
	return u.userRepo.CreateUser(ctx, user)
}

func (u *UseCase) GetUserByID(ctx context.Context, id string) (*models.UserDTO, error) {
	return u.userRepo.GetUserByID(ctx, id)
}

func (u *UseCase) GetAllUsers(ctx context.Context) ([]models.UserDTO, error) {
	return u.userRepo.GetAllUsers(ctx)
}

func (u *UseCase) DeleteUserByID(ctx context.Context, id string) error {
	return u.userRepo.DeleteUserByID(ctx, id)
}

func (u *UseCase) UpdateUser(ctx context.Context, user models.UserDTO) error {
	return u.userRepo.UpdateUser(ctx, user)
}
