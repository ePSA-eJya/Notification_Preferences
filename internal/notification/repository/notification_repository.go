package repository

import (
	"context"

	"Notification_Preferences/internal/entities"

	"github.com/google/uuid"
)

type NotificationRepository interface {
	UpdateStatusByID(ctx context.Context, notifID uuid.UUID, status entities.DeliveryStatus, providerID string) error
}
