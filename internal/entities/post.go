package entities

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID           uuid.UUID `bson:"_id" json:"post_id"`
	UserID       uuid.UUID `bson:"user_id" json:"user_id"`
	ImageURL     string    `bson:"image_url" json:"image_url"`
	LikeCount    int       `bson:"like_count" json:"like_count"`
	CommentCount int       `bson:"comment_count" json:"comment_count"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
}

type Comment struct {
	ID        uuid.UUID `bson:"_id" json:"comment_id"`
	PostID    uuid.UUID `bson:"post_id" json:"post_id"`
	UserID    uuid.UUID `bson:"user_id" json:"user_id"`
	Text      string    `bson:"comment_text" json:"comment"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type Like struct {
	PostID    uuid.UUID `bson:"post_id" json:"post_id"`
	UserID    uuid.UUID `bson:"user_id" json:"user_id"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
