package entities

import (
	"time"

	"github.com/google/uuid"
)

// --- Custom Types for Delivery Status ---

type DeliveryStatus string

const (
	StatusDelivered DeliveryStatus = "DELIVERED"
	StatusSent      DeliveryStatus = "SENT"
	StatusSkipped   DeliveryStatus = "SKIPPED"
	StatusFailed    DeliveryStatus = "FAILED"
	StatusPending   DeliveryStatus = "PENDING"
)

type ChannelDelivery struct {
	Status     DeliveryStatus `bson:"status" json:"status"`
	ReadAt     *time.Time     `bson:"read_at,omitempty" json:"read_at,omitempty"`
	ProviderID *string        `bson:"provider_id,omitempty" json:"provider_id,omitempty"`
}

type NotificationChannels struct {
	InApp ChannelDelivery `bson:"in_app" json:"in_app"`
	Push  ChannelDelivery `bson:"push" json:"push"`
	Email ChannelDelivery `bson:"email" json:"email"`
}

// --- Main Entities ---

type Notification struct {
	ID          uuid.UUID            `bson:"_id" json:"id"`
	RecipientID uuid.UUID            `bson:"recipient_id" json:"recipient_id"`
	EventID     uuid.UUID            `bson:"event_id" json:"event_id"`
	Message     string               `bson:"message" json:"message"`
	Channels    NotificationChannels `bson:"channels" json:"channels"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
}
