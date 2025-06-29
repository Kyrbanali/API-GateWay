package usecase

import (
	"context"

	"github.com/Kyrbanali/API-GateWay/internal/models"
	"github.com/Kyrbanali/API-GateWay/internal/repository"
)

type UseCase struct {
	userRepo *repository.UserRepo
}

func New(userRepo *repository.UserRepo) *UseCase {
	return &UseCase{userRepo: userRepo}
}

func (u *UseCase) CreateUser(ctx context.Context, user models.User) (string, error) {
	return u.userRepo.CreateUser(ctx, user)
}

func (u *UseCase) GetUserByID(ctx context.Context, id string) (models.User, error) {
	return u.userRepo.GetUserByID(ctx, id)
}
