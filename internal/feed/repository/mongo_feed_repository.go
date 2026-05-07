package repository

import (
	"Notification_Preferences/internal/entities"
	"context"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoFeedRepository struct {
	postCollection    *mongo.Collection
	likeCollection    *mongo.Collection
	commentCollection *mongo.Collection
	feedCollection    *mongo.Collection
}

func NewMongoFeedRepository(db *mongo.Database) FeedRepository {
	return &MongoFeedRepository{
		postCollection:    db.Collection("posts"),
		likeCollection:    db.Collection("likes"),
		commentCollection: db.Collection("comments"),
		feedCollection:    db.Collection("feeds"),
	}
}

func (r *MongoFeedRepository) CreatePost(ctx context.Context, post *entities.Post) error {
	if post.ID == uuid.Nil {
		post.ID = uuid.New()
	}
	_, err := r.postCollection.InsertOne(ctx, post)
	return err
}

func (r *MongoFeedRepository) GetPostByID(ctx context.Context, id string) (*entities.Post, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var post entities.Post
	filter := bson.M{"_id": uuidID}
	err = r.postCollection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return &post, nil
}

func (r *MongoFeedRepository) SaveLike(ctx context.Context, like *entities.Like) error {
	if like.ID == uuid.Nil {
		like.ID = uuid.New()
	}
	filter := bson.M{"post_id": like.PostID, "user_id": like.UserID}
	// Use upsert to ensure idempotency: only insert if not exists
	update := bson.M{"$setOnInsert": like}
	opts := options.Update().SetUpsert(true)

	_, err := r.likeCollection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *MongoFeedRepository) RemoveLike(ctx context.Context, postID, userID uuid.UUID) error {
	filter := bson.M{"post_id": postID, "user_id": userID}
	result, err := r.likeCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("like not found")
	}
	return nil
}

func (r *MongoFeedRepository) SaveComment(ctx context.Context, comment *entities.Comment) error {
	if comment.ID == uuid.Nil {
		comment.ID = uuid.New()
	}
	_, err := r.commentCollection.InsertOne(ctx, comment)
	return err
}

func (r *MongoFeedRepository) AddFeedItems(ctx context.Context, items []*entities.FeedItem) error {
	if len(items) == 0 {
		return nil
	}
	docs := make([]interface{}, 0, len(items))
	for _, it := range items {
		if it.ID == uuid.Nil {
			it.ID = uuid.New()
		}
		docs = append(docs, it)
	}
	_, err := r.feedCollection.InsertMany(ctx, docs)
	return err
}

func (r *MongoFeedRepository) GetUserTimeline(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.FeedItem, error) {
	filter := bson.M{"user_id": userID}
	findOpts := options.Find()
	if limit > 0 {
		findOpts.SetLimit(int64(limit))
	}
	if offset > 0 {
		findOpts.SetSkip(int64(offset))
	}
	// return newest first
	findOpts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.feedCollection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*entities.FeedItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}
