package repository

import (
	"Notification_Preferences/internal/entities"
	"context"
	"errors"
	"fmt"

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
	client            *mongo.Client
}

func NewMongoFeedRepository(db *mongo.Database) FeedRepository {
	return &MongoFeedRepository{
		postCollection:    db.Collection("posts"),
		likeCollection:    db.Collection("likes"),
		commentCollection: db.Collection("comments"),
		feedCollection:    db.Collection("feeds"),
		client:            db.Client(),
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

	session, err := r.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start mongo session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {

		filter := bson.M{"post_id": like.PostID, "user_id": like.UserID}
		// Use upsert to ensure idempotency: only insert if not exists
		update := bson.M{"$setOnInsert": like}
		opts := options.Update().SetUpsert(true)
		_, err := r.likeCollection.UpdateOne(sessCtx, filter, update, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to save like: %w", err)
		}

		likeFilter := bson.M{"_id": like.PostID}
		updateLike := bson.M{"$inc": bson.M{"like_count": 1}}
		_, dbErr := r.postCollection.UpdateOne(sessCtx, likeFilter, updateLike)
		if dbErr != nil {
			return nil, fmt.Errorf("failed to update like count: %w", err)
		}

		return nil, nil
	})
	return err
}

func (r *MongoFeedRepository) RemoveLike(ctx context.Context, postID, userID uuid.UUID) error {
	session, err := r.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start mongo session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		filter := bson.M{"post_id": postID, "user_id": userID}
		result, err := r.likeCollection.DeleteOne(ctx, filter)
		if err != nil {
			return nil, fmt.Errorf("failed to remove like: %w", err)
		}
		if result.DeletedCount == 0 {
			return nil, errors.New("like not found")
		}

		unlikeFilter := bson.M{"_id": postID}
		update := bson.M{"$inc": bson.M{"like_count": -1}}
		_, dbErr := r.postCollection.UpdateOne(ctx, unlikeFilter, update)
		if dbErr != nil {
			return nil, fmt.Errorf("failed to update like count: %w", err)
		}
		return nil, nil
	})
	return err
}

func (r *MongoFeedRepository) IsPostLikedByUser(ctx context.Context, postID, userID uuid.UUID) (bool, error) {
	filter := bson.M{"post_id": postID, "user_id": userID}
	err := r.likeCollection.FindOne(ctx, filter).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *MongoFeedRepository) SaveComment(ctx context.Context, comment *entities.Comment) error {
	if comment.ID == uuid.Nil {
		comment.ID = uuid.New()
	}

	session, err := r.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start mongo session: %w", err)
	}
	defer session.EndSession(ctx)

	//execute with transaction
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		_, err := r.commentCollection.InsertOne(sessCtx, comment)
		if err != nil {
			return nil, fmt.Errorf("failed to insert comment: %w", err)
		}

		commentFilter := bson.M{"_id": comment.PostID}
		update := bson.M{"$inc": bson.M{"comment_count": 1}}

		_, dbErr := r.postCollection.UpdateOne(sessCtx, commentFilter, update)
		if dbErr != nil {
			return nil, fmt.Errorf("failed to update post comment count: %w", err)
		}

		return nil, nil
	})

	return err
}

func (r *MongoFeedRepository) GetCommentsByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*entities.Comment, error) {
	filter := bson.M{"post_id": postID}
	findOpts := options.Find()
	if limit > 0 {
		findOpts.SetLimit(int64(limit))
	}
	if offset > 0 {
		findOpts.SetSkip(int64(offset))
	}
	findOpts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.commentCollection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []*entities.Comment
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, err
	}

	return comments, nil
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
