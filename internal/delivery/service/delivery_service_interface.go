package service

import (
	"context"

	"github.com/google/uuid"
)

type DeliveryService interface {
	SendGmail(ctx context.Context, notifID *uuid.UUID, recipientEmail []string, subject, body string) error
	SendPush(ctx context.Context, notifID *uuid.UUID, deviceToken string, message string) error
	// HandleInApp(ctx context.Context, notifID string) error
	// Send(ctx context.Context, to, subject, body string) error //sendgrid
}
