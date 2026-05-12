package repository

import (
	"context"

	"Notification_Preferences/internal/entities"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PreferenceRepository interface {
	GetPreferenceByUserID(ctx context.Context, userID uuid.UUID) (entities.NotificationPreferences, error)
	UpsertPreference(ctx context.Context, userID uuid.UUID, prefs entities.NotificationPreferences) error
}

type MongoPreferenceRepository struct {
	collection *mongo.Collection
}

func NewMongoPreferenceRepository(db *mongo.Database) PreferenceRepository {
	return &MongoPreferenceRepository{
		collection: db.Collection("preferences"),
	}
}

func (r *MongoPreferenceRepository) GetPreferenceByUserID(ctx context.Context, userID uuid.UUID) (entities.NotificationPreferences, error) {
	var prefs entities.NotificationPreferences
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&prefs)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return default preferences if none found
			return entities.NotificationPreferences{
				Likes:    entities.ChannelConfig{InApp: entities.PrefAll, Push: entities.PrefAll, Email: entities.PrefNone},
				Comments: entities.ChannelConfig{InApp: entities.PrefAll, Push: entities.PrefAll, Email: entities.PrefNone},
				Follows:  entities.ChannelConfig{InApp: entities.PrefAll, Push: entities.PrefAll, Email: entities.PrefNone},
				Posts:    entities.ChannelConfig{InApp: entities.PrefAll, Push: entities.PrefAll, Email: entities.PrefNone},
			}, nil
		}
		return prefs, err
	}
	return prefs, nil
}

func (r *MongoPreferenceRepository) UpsertPreference(ctx context.Context, userID uuid.UUID, prefs entities.NotificationPreferences) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"user_id": userID,
			"likes":   prefs.Likes,
			"comments": prefs.Comments,
			"follows":  prefs.Follows,
			"posts":    prefs.Posts,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}
