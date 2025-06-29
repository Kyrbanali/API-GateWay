package repository

import (
	"context"
	"fmt"

	"github.com/Kyrbanali/API-GateWay/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *UserRepo {
	return &UserRepo{conn: conn}
}

func (r *UserRepo) CreateUser(ctx context.Context, user models.UserDTO) (string, error) {
	_, err := r.conn.Exec(ctx, `
		INSERT INTO users (id, name, age) VALUES ($1, $2, $3)
	`, user.ID, user.Name, user.Age)

	if err != nil {
		return "", err
	}

	return user.ID, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id string) (models.UserDTO, error) {
	var user models.UserDTO

	row := r.conn.QueryRow(ctx, `
	    SELECT id, name, age FROM users WHERE id = $1
	`, id)

	err := row.Scan(&user.ID, &user.Name, &user.Age)

	return user, err
}

func (r *UserRepo) GetAllUsers(ctx context.Context) ([]models.UserDTO, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT id, name, age FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserDTO

	for rows.Next() {
		var user models.UserDTO
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepo) DeleteUserByID(ctx context.Context, id string) error {
	cmd, err := r.conn.Exec(ctx, `
		DELETE FROM users WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, user models.UserDTO) error {
	cmd, err := r.conn.Exec(ctx, `
		UPDATE users SET name = $1, age = $2 WHERE id = $3
	`, user.Name, user.Age, user.ID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
