package cache

import (
	"context"
	"sync"
	"time"
	"unsafe"

	"github.com/Kyrbanali/API-GateWay/internal/metrics"
	"github.com/Kyrbanali/API-GateWay/internal/models"
	"github.com/Kyrbanali/API-GateWay/internal/repository"
	"github.com/pkg/errors"
)

type Decorator struct {
	ttl      time.Duration
	userRepo repository.UserProvider

	rwmu  sync.RWMutex
	users map[string]WrapUser
}

type WrapUser struct {
	user      models.UserDTO
	updatedAt time.Time
}

func (d *Decorator) get(id string) (*models.UserDTO, bool) {
	d.rwmu.RLock()
	defer d.rwmu.RUnlock()
	u, ok := d.users[id]
	if ok {
		return &u.user, true
	}
	return nil, false
}

func (d *Decorator) set(id string, user models.UserDTO) {
	d.rwmu.Lock()
	d.users[id] = WrapUser{user: user, updatedAt: time.Now()}
	items, bytes := d.sizeLocked()
	d.rwmu.Unlock()
	metrics.SetCacheStats(items, bytes)
}

func (d *Decorator) delete(id string) {
	d.rwmu.Lock()
	delete(d.users, id)
	items, bytes := d.sizeLocked()
	d.rwmu.Unlock()
	metrics.SetCacheStats(items, bytes)
}

func New(repo repository.UserProvider, ttl time.Duration, cleanupInterval time.Duration) *Decorator {
	d := &Decorator{
		ttl:      ttl,
		users:    make(map[string]WrapUser),
		userRepo: repo,
	}
	d.startCleanup(cleanupInterval)
	metrics.SetCacheStats(0, 0)
	return d
}

func (d *Decorator) startCleanup(cleanupInterval time.Duration) {
	go func() {
		t := time.NewTicker(cleanupInterval)
		defer t.Stop()
		for range t.C {
			d.rwmu.Lock()
			for id, entry := range d.users {
				if time.Since(entry.updatedAt) > d.ttl {
					delete(d.users, id)
				}
			}
			items, bytes := d.sizeLocked()
			d.rwmu.Unlock()
			metrics.SetCacheStats(items, bytes)
		}
	}()
}

func (d *Decorator) sizeLocked() (items int, bytes int64) {
	items = len(d.users)
	var total int64
	for k, w := range d.users {
		total += int64(len(k))
		total += int64(unsafe.Sizeof(w))
		total += int64(len(w.user.ID))
		total += int64(len(w.user.Name))
		total += 8 // Age
		total += int64(unsafe.Sizeof(w.updatedAt))
	}
	return items, total
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
	if err := d.userRepo.DeleteUserByID(ctx, id); err != nil {
		return errors.Wrap(err, "DeleteUser")
	}

	d.delete(id)
	return nil
}

func (d *Decorator) UpdateUser(ctx context.Context, user models.UserDTO) error {
	if err := d.userRepo.UpdateUser(ctx, user); err != nil {
		return errors.Wrap(err, "UpdateUser")
	}

	d.set(user.ID, user)
	return nil
}
