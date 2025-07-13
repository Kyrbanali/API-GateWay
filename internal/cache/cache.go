package cache

import (
	"context"
	"sync"
	"time"

	"github.com/Kyrbanali/API-GateWay/internal/models"
	"github.com/Kyrbanali/API-GateWay/internal/repository"
	"github.com/pkg/errors"
)

type Decorator struct {
	ttl      time.Duration
	mutex    sync.Mutex
	users    map[string]WrapUser
	userRepo repository.UserProvider
}

type WrapUser struct {
	user      models.UserDTO
	updatedAt time.Time
}

func New(repo repository.UserProvider, ttl time.Duration) *Decorator {
	d := &Decorator{
		ttl:      ttl,
		users:    make(map[string]WrapUser),
		userRepo: repo,
	}
	d.startEvictionLoop()
	return d
}

func (d *Decorator) startEvictionLoop() {
	go func() {
		for {
			time.Sleep(10 * time.Second)

			d.mutex.Lock()
			for id, entry := range d.users {
				if time.Now().After(entry.updatedAt.Add(d.ttl)) {
					delete(d.users, id)
				}
			}
			d.mutex.Unlock()
		}
	}()
}

func (d *Decorator) GetUserByID(ctx context.Context, id string) (*models.UserDTO, error) {
	d.mutex.Lock()
	user, ok := d.users[id]
	if ok {
		d.mutex.Unlock()
		return &user.user, nil
	}
	d.mutex.Unlock()

	userPtr, err := d.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "GetUser")
	}

	d.mutex.Lock()
	d.users[id] = WrapUser{
		user:      *userPtr,
		updatedAt: time.Now(),
	}
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
	d.users[id] = WrapUser{
		user:      user,
		updatedAt: time.Now(),
	}
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
	d.users[user.ID] = WrapUser{
		user:      user,
		updatedAt: time.Now(),
	}
	d.mutex.Unlock()

	return nil
}
