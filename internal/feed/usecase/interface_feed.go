package usecase

import (
	"Notification_Preferences/internal/entities"
	"context"

	"github.com/google/uuid"
)

type FeedUseCase interface {
	CreatePost(ctx context.Context, post *entities.Post) error
	LikePost(ctx context.Context, like *entities.Like) error
	UnlikePost(ctx context.Context, postID, userID uuid.UUID) error
	IsPostLikedByUser(ctx context.Context, postID, userID uuid.UUID) (bool, error)
	CommentOnPost(ctx context.Context, comment *entities.Comment) error
	GetPostComments(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*entities.Comment, error)
	GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Post, error)
}
