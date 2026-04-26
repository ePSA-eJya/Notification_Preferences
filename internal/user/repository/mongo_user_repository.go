// internal/repository/mongo_user_repository.go
package repository

import (
	"Notification_Preferences/internal/entities"
	"context"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	collection       *mongo.Collection
	followCollection *mongo.Collection
}

// We pass in the specific MongoDB collection this repo will manage
func NewMongoUserRepository(db *mongo.Database) UserRepository {
	return &MongoUserRepository{
		collection:       db.Collection("users"),
		followCollection: db.Collection("follows"),
	}
}

func (r *MongoUserRepository) Save(ctx context.Context, user *entities.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	filter := bson.M{"email": email}

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepository) FindByID(ctx context.Context, id string) (*entities.User, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var user entities.User
	filter := bson.M{"_id": uuidID}

	err = r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepository) Patch(ctx context.Context, id string, user *entities.User) error {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": uuidID}
	// We use $set to only update the provided fields
	update := bson.M{"$set": user}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *MongoUserRepository) Delete(ctx context.Context, id string) error {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": uuidID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *MongoUserRepository) FindAll(ctx context.Context) ([]*entities.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*entities.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *MongoUserRepository) IsFollowing(ctx context.Context, followerID, followeeID uuid.UUID) (bool, error) {
	filter := bson.M{
		"follower_id": followerID,
		"followee_id": followeeID,
	}

	err := r.followCollection.FindOne(ctx, filter).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *MongoUserRepository) CreateFollow(ctx context.Context, follow *entities.Follow) error {
	_, err := r.followCollection.InsertOne(ctx, follow)
	return err
}

func (r *MongoUserRepository) DeleteFollow(ctx context.Context, followerID, followeeID uuid.UUID) error {
	filter := bson.M{
		"follower_id": followerID,
		"followee_id": followeeID,
	}

	result, err := r.followCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("follow not found")
	}

	return nil
}
