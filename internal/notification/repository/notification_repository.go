package repository

import (
	"Notification_Preferences/internal/entities"
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRepositoryImpl struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewNotificationRepository(db *mongo.Database) NotificationRepository {
	return &NotificationRepositoryImpl{db: db, collection: db.Collection("notifications")}
}

func (r *NotificationRepositoryImpl) Create(ctx context.Context, notification *entities.Notification) error {
	_, err := r.collection.InsertOne(ctx, notification)
	return err
}

func (r *NotificationRepositoryImpl) UpdateStatusByID(ctx context.Context, notificationId *uuid.UUID, status entities.DeliveryStatus, providerId string) error {
	filter := bson.M{"_id": notificationId}
	update := bson.M{
		"$set": bson.M{
			"status":      status,
			"provider_id": providerId,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err

}
