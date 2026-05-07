package usecase

import (
	"Notification_Preferences/internal/entities"
	"context"

	"github.com/google/uuid"
)

type FeedUseCase interface {
	CreatePost(ctx context.Context, post *entities.Post) error
	LikePost(ctx context.Context, like *entities.Like) error
	CommentOnPost(ctx context.Context, comment *entities.Comment) error
	GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Post, error)
}
