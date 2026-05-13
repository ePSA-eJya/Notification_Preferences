package usecase

import (
	"Notification_Preferences/internal/entities"
	"Notification_Preferences/internal/notification/repository"
	"Notification_Preferences/pkg/config"
	"context"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/google/uuid"

	fcm "firebase.google.com/go/v4/messaging"
)

type DeliveryServiceImpl struct {
	// repo      repository.DeliveryRepository
	notificationRepo repository.NotificationRepository
	smtpHost         string
	smtpPort         string
	email            string
	password         string
	fcmClient        *fcm.Client
}

func NewDeliveryService(notificationRepo repository.NotificationRepository, config config.SMTPConfig, fcmClient *fcm.Client) DeliveryService {
	return &DeliveryServiceImpl{
		// repo:      repo,
		notificationRepo: notificationRepo,
		smtpHost:         config.Host,
		smtpPort:         config.Port,
		email:            config.Email,
		password:         config.Password,
		fcmClient:        fcmClient,
	}
}

func (s *DeliveryServiceImpl) SendGmail(ctx context.Context, notifID *uuid.UUID, recipientEmail []string, subject, body string) error {
	auth := smtp.PlainAuth("", s.email, s.password, s.smtpHost)
	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
			"\r\n"+
			"%s\r\n",
		s.email,
		strings.Join(recipientEmail, ","),
		subject,
		body,
	))

	err := smtp.SendMail(s.smtpHost+":"+s.smtpPort, auth, s.email, recipientEmail, msg)
	if err != nil {
		log.Printf("failed to send email for notifID=%s", notifID, err)
		dbErr := s.notificationRepo.UpdateStatusByID(ctx, *notifID, entities.StatusFailed, "")
		return dbErr
	}

	providerID := "smtp"
	dbErr := s.notificationRepo.UpdateStatusByID(ctx, *notifID, entities.StatusDelivered, providerID)
	return dbErr
}

func (s *DeliveryServiceImpl) SendPush(ctx context.Context, notifID *uuid.UUID, deviceToken string, message string) error {
	if deviceToken == "" {
		log.Printf("device token is empty for notifID=%s", notifID)
		dbErr := s.notificationRepo.UpdateStatusByID(ctx, *notifID, entities.StatusSkipped, "")
		return dbErr
	}

	log.Println("SENDING PUSH TO:", deviceToken)

	messagePayload := &fcm.Message{
		Token: deviceToken,
		Notification: &fcm.Notification{
			Title: "New Notification",
			Body:  message,
		},
		Data: map[string]string{
			"notification_id": notifID.String(),
		},
	}

	log.Printf("message %s", messagePayload)
	if s.fcmClient == nil {
		log.Println("FCM CLIENT IS NIL")
		return fmt.Errorf("fcm client is nil")
	}

	response, err := s.fcmClient.Send(ctx, messagePayload)

	if err != nil {
		log.Printf("FCM SEND FAILED notifID=%s err=%v", notifID, err)
		return s.notificationRepo.UpdateStatusByID(ctx, *notifID, entities.StatusFailed, "")
	}
	log.Printf("FCM SUCCESS notifID=%s response=%s", notifID, response)

	providerID := response
	return s.notificationRepo.UpdateStatusByID(ctx, *notifID, entities.StatusDelivered, providerID)
}
