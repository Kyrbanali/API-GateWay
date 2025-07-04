package cache

import (
	"context"

	"github.com/Kyrbanali/API-GateWay/internal/models"
	"github.com/Kyrbanali/API-GateWay/internal/repository"
	"github.com/pkg/errors"
)

type Decorator struct {
	users    map[string]models.UserDTO
	userRepo repository.UserProvider
}

func New() *Decorator {
	return &Decorator{}
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

func (d *Decorator) GetAllUsers(ctx context.Context) error {}

func CreateUser(ctx context.Context, user models.UserDTO) (string, error) {}

func DeleteUserByID(ctx context.Context, id string) error {}

func UpdateUser(ctx context.Context, user models.UserDTO) error {}
