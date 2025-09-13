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
	mu       sync.Mutex
	users    map[string]WrapUser
	userRepo repository.UserProvider
}

type WrapUser struct {
	user      models.UserDTO
	updatedAt time.Time
}

func (d *Decorator) get(id string) (*models.UserDTO, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	user, ok := d.users[id]
	if ok {
		return &user.user, true
	}
	return nil, false
}

func (d *Decorator) set(id string, user models.UserDTO) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.users[id] = WrapUser{
		user:      user,
		updatedAt: time.Now(),
	}
}

func (d *Decorator) delete(id string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.users, id)
}

func New(repo repository.UserProvider, ttl time.Duration, cleanupInterval time.Duration) *Decorator {
	d := &Decorator{
		ttl:      ttl,
		users:    make(map[string]WrapUser),
		userRepo: repo,
	}
	d.startEvictionLoop(cleanupInterval)
	return d
}

func (d *Decorator) startEvictionLoop(cleanupInterval time.Duration) {
	go func() {
		for {
			time.Sleep(cleanupInterval)

			d.mu.Lock()
			for id, entry := range d.users {
				if time.Now().After(entry.updatedAt.Add(d.ttl)) {
					delete(d.users, id)
				}
			}
			d.mu.Unlock()
		}
	}()
}

func (d *Decorator) GetUserByID(ctx context.Context, id string) (*models.UserDTO, error) {
	if user, ok := d.get(id); ok {
		return user, nil
	}

	userPtr, err := d.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "GetUser")
	}

	d.set(id, *userPtr)

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

	d.set(id, user)

	return id, nil
}

func (d *Decorator) DeleteUserByID(ctx context.Context, id string) error {
	err := d.userRepo.DeleteUserByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "DeleteUser")
	}

	d.delete(id)
	return nil
}

func (d *Decorator) UpdateUser(ctx context.Context, user models.UserDTO) error {
	err := d.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return errors.Wrap(err, "UpdateUser")
	}

	d.set(user.ID, user)
	return nil
}
