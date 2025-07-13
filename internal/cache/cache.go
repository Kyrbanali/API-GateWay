package cache

import (
	"context"
	"sync"

	"github.com/Kyrbanali/API-GateWay/internal/models"
	"github.com/Kyrbanali/API-GateWay/internal/repository"
	"github.com/pkg/errors"
)

type Decorator struct {
	mutex    sync.Mutex
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
	d.mutex.Lock()
	user, ok := d.users[id]
	if ok {
		d.mutex.Unlock()
		return &user, nil
	}
	d.mutex.Unlock()

	userPtr, err := d.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "GetUser")
	}

	d.mutex.Lock()
	d.users[id] = *userPtr
	d.mutex.Unlock()

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

	d.mutex.Lock()
	d.users[id] = user
	d.mutex.Unlock()

	return id, nil
}

func (d *Decorator) DeleteUserByID(ctx context.Context, id string) error {
	err := d.userRepo.DeleteUserByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "DeleteUser")
	}

	d.mutex.Lock()
	delete(d.users, id)
	d.mutex.Unlock()

	return nil
}

func (d *Decorator) UpdateUser(ctx context.Context, user models.UserDTO) error {
	err := d.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return errors.Wrap(err, "UpdateUser")
	}

	d.mutex.Lock()
	d.users[user.ID] = user
	d.mutex.Unlock()

	return nil
}
