package repository

import (
	"context"

	"Notification_Preferences/internal/entities"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 1. The Interface (From your HEAD)
type NotificationRepository interface {
	Create(ctx context.Context, notification *entities.Notification) error
	UpdateStatusByID(ctx context.Context, notifID uuid.UUID, status entities.DeliveryStatus, providerID string) error
}

// 2. The Struct (From Remote)
type NotificationRepositoryImpl struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// 3. Constructor
func NewNotificationRepository(db *mongo.Database) NotificationRepository {
	return &NotificationRepositoryImpl{db: db, collection: db.Collection("notifications")}
}

// 4. Implementations
func (r *NotificationRepositoryImpl) Create(ctx context.Context, notification *entities.Notification) error {
	_, err := r.collection.InsertOne(ctx, notification)
	return err
}

func (r *NotificationRepositoryImpl) UpdateStatusByID(ctx context.Context, notifID uuid.UUID, status entities.DeliveryStatus, providerID string) error {
	filter := bson.M{"_id": notifID}
	update := bson.M{
		"$set": bson.M{
			"status":      status,
			"provider_id": providerID,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
