package dto

import "github.com/google/uuid"

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	UserHandle string    `json:"user_handle"`
}
