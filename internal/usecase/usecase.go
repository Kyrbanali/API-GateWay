package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/Kyrbanali/API-GateWay/internal/models"
	"github.com/Kyrbanali/API-GateWay/internal/repository"
	"github.com/Kyrbanali/API-GateWay/internal/worker"
)

type UseCase struct {
	userRepo repository.UserProvider
	worker   *worker.Worker
}

func New(userRepo repository.UserProvider, worker *worker.Worker) *UseCase {
	return &UseCase{userRepo: userRepo, worker: worker}
}

func (u *UseCase) CreateUser(ctx context.Context, user models.UserDTO) (string, error) {
	id, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	u.worker.Push(ctx, u.someJob(id))

	return id, nil
}

func (u *UseCase) someJob(userID string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		if err := u.doSome(ctx, userID); err != nil {
			return err
		}
		return nil
	}
}

func (u *UseCase) doSome(ctx context.Context, userID string) error {
	user, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("%s %d", user.Name, user.Age)
	links, err := u.worker.FetchLinks(query, 3)
	if err != nil {
		return err
	}
	log.Printf("links for user %s (%s, %d): %v", user.ID, user.Name, user.Age, links)
	return err
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
