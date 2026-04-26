package repository

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetEmailByUserID(ctx context.Context, userID uuid.UUID) (string, error)
	GetDeviceTokenByUserID(ctx context.Context, userID uuid.UUID) (string, error)
}
