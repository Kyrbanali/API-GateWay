package usecase

import (
	"context"

	"github.com/Kyrbanali/API-GateWay/internal/models"
	"github.com/Kyrbanali/API-GateWay/internal/repository"
	"github.com/Kyrbanali/API-GateWay/internal/worker"
)

type UseCase struct {
	userRepo repository.UserProvider
	jobs     chan<- worker.Job
}

func New(userRepo repository.UserProvider, jobs chan<- worker.Job) *UseCase {
	return &UseCase{userRepo: userRepo, jobs: jobs}
}

func (u *UseCase) CreateUser(ctx context.Context, user models.UserDTO) (string, error) {
	id, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	if u.jobs != nil {
		u.jobs <- worker.Job{ID: id}
	}

	return id, nil
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
