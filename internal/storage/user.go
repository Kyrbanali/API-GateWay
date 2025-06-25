package storage

import (
	"context"

	"github.com/Kyrbanali/API-GateWay/internal/models"
)

func GetAllUsers(ctx context.Context) ([]models.User, error) {
	rows, err := Conn.Query(ctx, "SELECT id, name, age FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		user := models.User{}
		rows.Scan(&user.ID, &user.Name, &user.Age)
		users = append(users, user)
	}
	return users, nil
}

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
