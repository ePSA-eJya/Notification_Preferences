package entities

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `bson:"_id" json:"id"`
	UserID    uuid.UUID `bson:"user_id" json:"user_id"`
	Content   string    `bson:"content" json:"content"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type Comment struct {
	ID        uuid.UUID `bson:"_id" json:"id"`
	PostID    uuid.UUID `bson:"post_id" json:"post_id"`
	UserID    uuid.UUID `bson:"user_id" json:"user_id"`
	Text      string    `bson:"text" json:"text"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type Like struct {
	ID        uuid.UUID `bson:"_id" json:"id"`
	PostID    uuid.UUID `bson:"post_id" json:"post_id"`
	UserID    uuid.UUID `bson:"user_id" json:"user_id"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type FeedItem struct {
	ID        uuid.UUID `bson:"_id" json:"id"`
	UserID    uuid.UUID `bson:"user_id" json:"user_id"`
	PostID    uuid.UUID `bson:"post_id" json:"post_id"`
	AuthorID  uuid.UUID `bson:"author_id" json:"author_id"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
