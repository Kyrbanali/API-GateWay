package models

type CreateUserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UpdateUserRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type CreateUserResponse struct {
	ID string `json:"id"`
}

type UserDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}
