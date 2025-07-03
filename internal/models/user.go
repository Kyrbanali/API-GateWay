package models

type CreateUserRequest struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"gte=18"`
}

type UpdateUserRequest struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"gte=18"`
}

type CreateUserResponse struct {
	ID string `json:"id" validate:"required"`
}

type UserDTO struct {
	ID   string `validate:"required"`
	Name string `validate:"required"`
	Age  int    `validate:"gte=18"`
}
