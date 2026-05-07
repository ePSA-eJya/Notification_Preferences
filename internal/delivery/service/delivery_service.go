package service

import (
	"Notification_Preferences/internal/entities"
	"Notification_Preferences/internal/notification/repository"
	"Notification_Preferences/pkg/config"
	"context"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"firebase.google.com/go/messaging"
	"github.com/google/uuid"
)

type DeliveryServiceImpl struct {
	// repo      repository.DeliveryRepository
	notificationRepo repository.NotificationRepository
	smtpHost         string
	smtpPort         string
	email            string
	password         string
	fcmClient        *messaging.Client
}

func NewDeliveryService(notificationRepo repository.NotificationRepository, config config.SMTPConfig, fcmClient *messaging.Client) DeliveryService {
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
	// recipient_email, dberr := s.repo.GetEmailByUserID(ctx, recipientID)
	// if dberr != nil {
	// 	log.Printf("failed to fetch recipient email for notifID=%s", notifID)
	// 	return fmt.Errorf("failed to fetch recipient email: %w", err)
	// }

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
		log.Printf("failed to send email for notifID=%s", notifID)
		dbErr := s.notificationRepo.UpdateStatusByID(ctx, notifID, entities.StatusFailed, "")
		return dbErr
	}

	providerID := "smtp"
	dbErr := s.notificationRepo.UpdateStatusByID(ctx, notifID, entities.StatusDelivered, providerID)
	return dbErr
}

func (s *DeliveryServiceImpl) SendPush(ctx context.Context, notifID *uuid.UUID, deviceToken string, message string) error {
	if deviceToken == "" {
		log.Printf("device token is empty for notifID=%s", notifID)
		dbErr := s.notificationRepo.UpdateStatusByID(ctx, notifID, entities.StatusSkipped, "")
		return dbErr
	}

	messagePayload := &messaging.Message{
		Token: deviceToken,
		Notification: &messaging.Notification{
			Title: "New Notification",
			Body:  message,
		},
		Data: map[string]string{
			"notification_id": notifID.String(),
		},
	}

	response, err := s.fcmClient.Send(ctx, messagePayload)

	if err != nil {
		log.Printf("failed to send push notification for notifID=%s", notifID)
		dbErr := s.notificationRepo.UpdateStatusByID(ctx, notifID, entities.StatusFailed, "")
		return dbErr
	}
	log.Printf("Successfully sent message: %s", response)

	providerID := response
	dbErr := s.notificationRepo.UpdateStatusByID(ctx, notifID, entities.StatusDelivered, providerID)
	return dbErr
}
