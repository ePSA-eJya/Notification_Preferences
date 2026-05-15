package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoFollowRepository struct {
	collection *mongo.Collection
}

func NewMongoFollowRepository(db *mongo.Database) FollowRepository {
	return &MongoFollowRepository{
		collection: db.Collection("follows"),
	}
}

func (r *MongoFollowRepository) IsFollowing(ctx context.Context, follower uuid.UUID, followee uuid.UUID) (bool, error) {
	filter := bson.M{
		"follower_id": follower,
		"followee_id": followee,
	}

	err := r.collection.FindOne(ctx, filter).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
