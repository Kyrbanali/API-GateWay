package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
	"unsafe"

	"github.com/Kyrbanali/API-GateWay/internal/metrics"
	"github.com/Kyrbanali/API-GateWay/internal/models"
	"github.com/Kyrbanali/API-GateWay/internal/repository"
	"github.com/pkg/errors"
)

type Config struct {
	TTL             time.Duration
	CleanupInterval time.Duration
}

type Decorator struct {
	ttl      time.Duration
	userRepo repository.UserProvider

	mu    sync.RWMutex
	users map[string]WrapUser

	metrics metrics.CacheMetrics
}

type WrapUser struct {
	user      models.UserDTO
	updatedAt time.Time
}

func (d *Decorator) get(id string) (*models.UserDTO, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	u, ok := d.users[id]
	if ok {
		return &u.user, true
	}
	return nil, false
}

func (d *Decorator) set(id string, user models.UserDTO) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.users[id] = WrapUser{user: user, updatedAt: time.Now()}
	d.updateMetrics()
}

func (d *Decorator) delete(id string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.users, id)
	d.updateMetrics()
}

func (d *Decorator) updateMetrics() {
	items, bytes := d.sizeLocked()
	d.metrics.Items.Set(float64(items))
	d.metrics.Bytes.Set(float64(bytes))
}

func New(repo repository.UserProvider, cfg Config) *Decorator {
	d := &Decorator{
		ttl:      cfg.TTL,
		users:    make(map[string]WrapUser),
		userRepo: repo,
	}
	d.metrics = metrics.NewCache(fmt.Sprintf("users_cache_%p", d))
	d.updateMetrics()

	d.startCleanup(cfg.CleanupInterval)
	return d
}

func (d *Decorator) startCleanup(cleanupInterval time.Duration) {
	go func() {
		t := time.NewTicker(cleanupInterval)
		defer t.Stop()
		for range t.C {
			d.mu.Lock()
			for id, entry := range d.users {
				if time.Since(entry.updatedAt) > d.ttl {
					delete(d.users, id)
				}
			}
			d.updateMetrics()
			d.mu.Unlock()
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
		total += 8
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
