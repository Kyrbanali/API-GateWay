package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Kyrbanali/API-GateWay/internal/link"
	"github.com/Kyrbanali/API-GateWay/internal/models"
	"github.com/Kyrbanali/API-GateWay/internal/repository"
	"github.com/Kyrbanali/API-GateWay/internal/worker"
	"github.com/pkg/errors"
)

type UseCase struct {
	userRepo repository.UserProvider
	worker   worker.WorkerProvider
	fetcher  *link.Fetcher
}

func New(userRepo repository.UserProvider, worker worker.WorkerProvider, fetcher *link.Fetcher) *UseCase {
	return &UseCase{userRepo: userRepo, worker: worker, fetcher: fetcher}
}

func (u *UseCase) CreateUser(ctx context.Context, user models.UserDTO) (string, error) {
	id, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	u.worker.Push(ctx, u.userLinks(user))

	return id, nil
}

func (u *UseCase) userLinks(user models.UserDTO) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		query := fmt.Sprintf("%s %d", user.Name, user.Age)
		links, err := u.fetcher.FetchLinks(query, 3)
		if err != nil {
			return errors.Wrap(err, "")
		}
		slog.Debug("links for user",
			slog.String("id", user.ID),
			slog.String("name", user.Name),
			slog.Int("age", user.Age),
			slog.Any("links", links),
		)
		return nil
	}
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
