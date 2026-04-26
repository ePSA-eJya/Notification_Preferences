package repository

import (
	"Notification_Preferences/internal/entities"
	"context"

	"github.com/google/uuid"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *entities.Notification) error
	UpdateStatusByID(ctx context.Context, notificationId *uuid.UUID, status entities.DeliveryStatus, providerId string) error
	// GetNotificationByUser(userId *uuid.UUID)
}
