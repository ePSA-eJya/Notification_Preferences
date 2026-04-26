package service

import (
	"Notification_Preferences/internal/delivery/service"
	"Notification_Preferences/internal/entities"
	"Notification_Preferences/internal/notification/repository"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

type NotificationServiceImpl struct {
	repo            repository.NotificationRepository
	followRepo      repository.FollowRepository
	userRepo        repository.UserRepository
	preferenceRepo  repository.PreferenceRepository
	deliveryService service.DeliveryService
	broker          MessageBroker
}

func NewNotificationService(repo repository.NotificationRepository, preferenceRepo repository.PreferenceRepository, broker MessageBroker) NotificationService {
	return &NotificationServiceImpl{repo: repo, preferenceRepo: preferenceRepo, broker: broker}
}

func (s *NotificationServiceImpl) FormatMessage(event *entities.Event) (string, err) {
	actorName, err := s.userRepo.GetNameByUserID(ctx, event.ActorID)
	if err != nil {
		log.Printf("failed to fetch actor name for userID=%s", event.ActorID)
		return "", err
	}

	var message string
	switch event.ActionType {
	case "POST":
		message = "%s published a new post.", actorName
	case "LIKE":
		message = "%s liked your post.", actorName
	case "COMMENT":
		message = "%s commented on your post.", actorName
	case "FOLLOW":
		message = "%s started following you.", actorName
	default:
		message = "New activity from %s.", actorName
	}

	return message, nil
}

func (s *NotificationServiceImpl) GetInitialStatus(enabled bool, targetStatus entities.DeliveryStatus) entities.DeliveryStatus {
	if !enabled {
		return entities.StatusSkipped
	}
	return targetStatus
}

func (s *NotificationServiceImpl) CreateNotification(event *entities.Event, recipientID uuid.UUID, enabledChannels map[entities.ChannelType]bool) (*entities.Notification, error) {
	message, err := s.FormatMessage(event)
	if err != nil {
		log.Printf("failed to format message: %v", err)
		return nil, err
	}

	channel := entities.NotificationChannels{
		InApp: entities.ChannelDelivery{
			// InApp is SENT immediately because it's now in the DB for the user to see
			Status: s.GetInitialStatus(enabledChannels[entities.InAppChannel], entities.StatusSent),
		},
		Push: entities.ChannelDelivery{
			Status: s.GetInitialStatus(enabledChannels[entities.PushChannel], entities.StatusPending),
		},
		Email: entities.ChannelDelivery{
			Status: s.GetInitialStatus(enabledChannels[entities.EmailChannel], entities.StatusPending),
		},
	}

	notification := entities.Notification{
		ID:          uuid.New(),
		RecipientID: recipientID,
		EventID:     event.ID,
		Message:     message,
		Channels:    channel,
		CreatedAt:   time.Now(),
	}

	return &notification, nil
}

func (s *NotificationServiceImpl) GetRecipientsByActionType(ctx context.Context, event *entities.Event) ([]uuid.UUID, error) {

	switch event.ActionType {

	case entities.Liked, entities.Commented:
		// send to owner of post
		return []uuid.UUID{event.EntityID}, nil

	case entities.Followed: //send to jisko followed
		return []uuid.UUID{event.EntityID}, nil

	case entities.Posted:
		// send to followers
		return s.userRepo.GetFollowers(ctx, event.ActorID)

	default:
		return nil, nil
	}
}

// func (s *NotificationServiceImpl) SendInAppNotif(event *entities.Event, recipientID uuid.UUID) {
// notification := s.CreateNotification(event, recipientID, entities.InAppChannel)
// err := s.repo.Create(ctx, notification)
// if err != nil {
// 	log.Printf("failed to save in-app notification: %v", err)
// }
// }

func (s *NotificationServiceImpl) SendPushNotif(ctx context.Context, event *entities.Event, notificationID *uuid.UUID, recipientID uuid.UUID) {
	deviceToken, err := s.userRepo.GetDeviceTokenByUserID(ctx, recipientID)
	if err != nil || deviceToken == "" {
		log.Printf("Skipping Push: No token found for user %s", recipientID)
		return
	}

	// 3. Publish to the Push Queue
	// The Delivery Service's SendPush will pick this up
	// pushPayload := entities.PushPayload{
	// 	NotificationID: notificationID,
	// 	DeviceToken:    deviceToken,
	// 	// Message:        notif.Message,
	// }
	// s.broker.Publish(entities.PushQueue, pushPayload)
	var message string

	s.deliveryService.SendPush(ctx, notificationID, deviceToken, message)

}

func (s *NotificationServiceImpl) SendEmailNotif(ctx context.Context, event *entities.Event, notificationID *uuid.UUID, recipientID uuid.UUID) {
	// 1. Fetch recipient email
	email, err := s.userRepo.GetEmailByUserID(ctx, recipientID)
	if err != nil || email == "" {
		log.Printf("Skipping Email: No email address for user %s", recipientID)
		return
	}

	// 3. Publish to the Email Queue
	// The Delivery Service's SendGmail will pick this up
	// emailPayload := entities.EmailPayload{
	// 	NotificationID: notif.ID,
	// 	RecipientEmail: []string{email},
	// 	Subject:        "New Activity on Your Account",
	// 	Body:           notif.Message,
	// }
	// s.broker.Publish(entities.EmailQueue, emailPayload)
	var message string
	s.deliveryService.SendGmail(ctx, notificationID, email, "New Activity on Your Account", message)
}

// GetActionPrefs extracts the channel settings for a specific action type from the user's full preferences.
func (s *NotificationServiceImpl) GetActionPrefs(actionType entities.ActionType, prefs entities.NotificationPreferences) entities.ChannelConfig {
	switch actionType {
	case entities.Posted:
		return prefs.Posts
	case entities.Liked:
		return prefs.Likes
	case entities.Commented:
		return prefs.Comments
	case entities.Followed:
		return prefs.Follows
	default:
		// Default to no notifications if action type is unrecognized
		return entities.ChannelConfig{
			InApp: entities.PrefNone,
			Push:  entities.PrefNone,
			Email: entities.PrefNone,
		}
	}
}

func (s *NotificationServiceImpl) CheckIfFollow(recipientID uuid.UUID, actorID uuid.UUID) bool {
	isFollowing, err := s.followRepo.IsFollowing(context.Background(), recipientID, actorID)
	if err != nil {
		log.Printf("Error checking follow status: %v", err)
		return false
	}
	return isFollowing
}

func (s *NotificationServiceImpl) ShouldNotify(recipientID uuid.UUID, actorID uuid.UUID, prefLevel entities.PreferenceLevel) bool {
	switch prefLevel {
	case entities.PrefAll:
		// User wants all notifications in this category
		return true

	case entities.PrefFollowers:
		// User only wants notifications from accounts they follow
		// s.Follow should check if recipientID follows actorID
		return s.CheckIfFollow(recipientID, actorID)

	case entities.PrefNone:
		return false

	default:
		return false
	}
}

// acha toh ham whi notifs bhejenge joh recipient ne allow kr rakhe h
// i send when i post (to all), like/comment (to author of the post)
func (s *NotificationServiceImpl) ProcessEvent(ctx context.Context, event *entities.Event) error {
	// action type is fixed --> P,L,C,F
	// iss user ka iss action k lie recipients kon hai?
	recipients, recepErr := s.GetRecipientsByActionType(ctx, event)
	if recepErr != nil {
		log.Printf("failed to fetch recipients for userID=%s", event.ActorID)
		return recepErr
	}

	for _, recipientID := range recipients {
		var prefs entities.NotificationPreferences
		// this is the recipients ka pref to receive notifs all actions(post/like/comment/follow)
		prefs, prefsErr := s.preferenceRepo.GetPreferenceByUserID(ctx, recipientID)
		if prefsErr != nil {
			log.Printf("failed to fetch preferences for userID=%s", recipientID)
			continue
		}
		//  for this particular action
		actionPrefs := s.GetActionPrefs(event.ActionType, prefs)
		enabledChannels := map[entities.ChannelType]bool{
			entities.InAppChannel: s.ShouldNotify(recipientID, event.ActorID, actionPrefs.InApp),
			entities.EmailChannel: s.ShouldNotify(recipientID, event.ActorID, actionPrefs.Email),
			entities.PushChannel:  s.ShouldNotify(recipientID, event.ActorID, actionPrefs.Push),
		}

		notification, err := s.CreateNotification(event, recipientID, enabledChannels)
		dbErr := s.repo.Create(ctx, notification)
		if dbErr != nil {
			log.Printf("failed to save notification: %v", err)
			continue
		}

		// if enabledChannels[entities.InAppChannel] {
		// 	s.SendInAppNotif(event, recipientID) // Persist to Mongo
		// }

		if enabledChannels[entities.PushChannel] {
			s.SendPushNotif(ctx, event, &notification.ID, recipientID) // Queue for Delivery Service
		}

		if enabledChannels[entities.EmailChannel] {
			s.SendEmailNotif(ctx, event, &notification.ID, recipientID) // Queue for Delivery Service
		}

	}

	return nil
}
