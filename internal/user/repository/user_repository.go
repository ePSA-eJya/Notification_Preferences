package repository

import (
	"Notification_Preferences/internal/entities"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryImpl struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &UserRepositoryImpl{db: db, collection: db.Collection("users")}
}

func (r *UserRepositoryImpl) GetEmailByUserID(ctx context.Context, userID uuid.UUID) (string, error) {
	var user entities.User
	err := r.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("no records found")
			return "", fmt.Errorf("no records found")
		}
		return "", err
	}
	return user.Email, nil
}

func (r *UserRepositoryImpl) GetDeviceTokenByUserID(ctx context.Context, userID uuid.UUID) (string, error) {
	var user entities.User
	err := r.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return "", err
	}
	return user.DeviceToken, nil
}
