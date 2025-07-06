package cache

import (
	"context"
	"log"

	"github.com/Kyrbanali/API-GateWay/internal/models"
	"github.com/Kyrbanali/API-GateWay/internal/repository"
	"github.com/pkg/errors"
)

type Decorator struct {
	users    map[string]models.UserDTO
	userRepo repository.UserProvider
}

func New(repo repository.UserProvider) *Decorator {
	return &Decorator{
		users:    make(map[string]models.UserDTO),
		userRepo: repo,
	}
}

func (d *Decorator) GetUserByID(ctx context.Context, id string) (*models.UserDTO, error) {
	user, ok := d.users[id]
	if ok {
		return &user, nil
	}

	userPtr, err := d.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "GetUser")
	}

	d.users[id] = *userPtr

	return userPtr, nil
}

func (d *Decorator) GetAllUsers(ctx context.Context) ([]models.UserDTO, error) {
	users, err := d.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "GetAllUsers")
	}

	return users, nil
}

func (d *Decorator) CreateUser(ctx context.Context, user models.UserDTO) (string, error) {
	id, err := d.userRepo.CreateUser(ctx, user)
	if err != nil {
		return "", errors.Wrap(err, "CreateUser")
	}

	d.users[id] = user

	return id, nil
}

func (d *Decorator) DeleteUserByID(ctx context.Context, id string) error {
	err := d.userRepo.DeleteUserByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "DeleteUser")
	}

	delete(d.users, id)
	log.Println("users cache", d.users)

	return nil
}

func (d *Decorator) UpdateUser(ctx context.Context, user models.UserDTO) error {
	err := d.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return errors.Wrap(err, "UpdateUser")
	}

	d.users[user.ID] = user

	return nil
}
