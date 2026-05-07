package repository

import (
	"context"

	"Notification_Preferences/internal/entities"

	"github.com/google/uuid"
)

// FeedRepository defines storage operations for posts, likes, comments and feed items.
type FeedRepository interface {
	CreatePost(ctx context.Context, post *entities.Post) error
	GetPostByID(ctx context.Context, id string) (*entities.Post, error)

	// SaveLike must be idempotent: a user can only like a post once.
	SaveLike(ctx context.Context, like *entities.Like) error
	RemoveLike(ctx context.Context, postID, userID uuid.UUID) error

	SaveComment(ctx context.Context, comment *entities.Comment) error

	// AddFeedItems performs a bulk insert (fan-out) of feed items.
	AddFeedItems(ctx context.Context, items []*entities.FeedItem) error

	// GetUserTimeline returns paginated FeedItems for the given timeline owner.
	GetUserTimeline(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.FeedItem, error)
}
