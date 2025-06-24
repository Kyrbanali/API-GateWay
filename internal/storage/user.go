package storage

import (
	"context"

	"github.com/Kyrbanali/API-GateWay/internal/models"
)

func InsertUser(ctx context.Context, user models.User) error {
	_, err := Conn.Exec(ctx, `
		INSERT INTO USERS (id, name, age) VALUES ($1, $2, $3)
	`, user.ID, user.Name, user.Age)
	return err
}

// func GetAllUsers(ctx context.Context)

func GetUserByID(ctx context.Context, id string) (models.User, error) {
	var user models.User
	err := Conn.QueryRow(ctx, `
		SELECT id, name, age FROM users WHERE id = $1
	`, id).Scan(&user.ID, &user.Name, &user.Age)
	return user, err
}

func UpdateUser(ctx context.Context, user models.User) error {
	_, err := Conn.Exec(ctx, `
		UPDATE users SET name = $1, age = $2 WHERE id = $3
	`, user.Name, user.Age, user.ID)
	return err
}

func DeleteUser(ctx context.Context, id string) error {
	_, err := Conn.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}
