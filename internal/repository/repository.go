package repository

import (
	"context"

	"github.com/Kyrbanali/API-GateWay/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *UserRepo {
	return &UserRepo{conn: conn}
}

func (r *UserRepo) CreateUser(ctx context.Context, user models.User) (string, error) {
	user.ID = uuid.NewString()

	_, err := r.conn.Exec(ctx, `
		INSERT INTO users (id, name, age) VALUES ($1, $2, $3)
	`, user.ID, user.Name, user.Age)

	if err != nil {
		return "", err
	}

	return user.ID, nil
}
