package service

import (
	"Notification_Preferences/internal/entities"
	"context"
)

type NotificationService interface {
	ProcessEvent(ctx context.Context, event *entities.Event) error
}
