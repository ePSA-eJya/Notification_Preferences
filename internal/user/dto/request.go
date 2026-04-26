package dto

type RegisterRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=6"`
	UserHandle string `json:"user_handle" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type PatchUserRequest struct {
	UserHandle string `json:"user_handle" validate:"required"`
}
