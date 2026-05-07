package repository

import (
	"context"

	"github.com/google/uuid"
)

type FollowRepository interface {
	IsFollowing(ctx context.Context, recipientID *uuid.UUID, actorID *uuid.UUID) (bool, error)
}
