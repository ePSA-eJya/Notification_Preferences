package routes

import (
	"context"

	"Notification_Preferences/internal/broker"
	deliveryUseCase "Notification_Preferences/internal/delivery/usecase"
	feedHandler "Notification_Preferences/internal/feed/handler/rest"
	feedRepository "Notification_Preferences/internal/feed/repository"
	feedUseCase "Notification_Preferences/internal/feed/usecase"
	notifRepository "Notification_Preferences/internal/notification/repository"
	notifUseCase "Notification_Preferences/internal/notification/usecase"
	preferenceRepository "Notification_Preferences/internal/preference/repository"
	userHandler "Notification_Preferences/internal/user/handler/rest"
	userRepository "Notification_Preferences/internal/user/repository"
	userUseCase "Notification_Preferences/internal/user/usecase"
	"Notification_Preferences/pkg/config"
	middleware "Notification_Preferences/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterPrivateRoutes(app fiber.Router, db *mongo.Database, publisher *broker.KafkaProducer, cfg *config.Config) {

	route := app.Group("/api/v1", middleware.JWTMiddleware())

	feedRepo := feedRepository.NewMongoFeedRepository(db)
	userRepo := userRepository.NewMongoUserRepository(db)
	prefRepo := preferenceRepository.NewMongoPreferenceRepository(db)
	notifRepo := notifRepository.NewNotificationRepository(db)

	deliveryService := deliveryUseCase.NewDeliveryService(notifRepo, cfg.SMTP, nil) // Assume nil fcmClient for now
	notifService := notifUseCase.NewNotificationService(notifRepo, userRepo, userRepo, prefRepo, deliveryService)

	eventTopic := ""
	if cfg != nil {
		eventTopic = cfg.KafkaEventTopic
	}

	// Start Kafka Consumer
	if cfg != nil && cfg.KafkaBrokerURL != "" {
		consumer := broker.NewKafkaConsumer([]string{cfg.KafkaBrokerURL}, cfg.KafkaEventTopic, "notification-group", notifService)
		go consumer.Start(context.Background())
	}

	feedService := feedUseCase.NewFeedUsecase(feedRepo, userRepo, publisher, eventTopic)
	feedHTTPHandler := feedHandler.NewHttpFeedHandler(feedService)

	userService := userUseCase.NewUserServiceWithPublisher(userRepo, publisher)
	userHTTPHandler := userHandler.NewHttpUserHandler(userService)

	route.Post("/posts", feedHTTPHandler.CreatePost)
	route.Post("/posts/:id/like", feedHTTPHandler.LikePost)
	route.Post("/posts/:id/comment", feedHTTPHandler.CommentOnPost)
	route.Get("/feed", feedHTTPHandler.GetFeed)

	route.Get("/me", userHTTPHandler.GetUser)
	route.Post("/users/:id/follow", userHTTPHandler.FollowUser)
	route.Delete("/users/:id/follow", userHTTPHandler.UnfollowUser)
}
