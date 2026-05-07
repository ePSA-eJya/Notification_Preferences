package routes

import (
	feedHandler "Notification_Preferences/internal/feed/handler/rest"
	feedRepository "Notification_Preferences/internal/feed/repository"
	feedUseCase "Notification_Preferences/internal/feed/usecase"
	userHandler "Notification_Preferences/internal/user/handler/rest"
	userRepository "Notification_Preferences/internal/user/repository"
	userUseCase "Notification_Preferences/internal/user/usecase"
	middleware "Notification_Preferences/pkg/middleware"

	"Notification_Preferences/internal/broker"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterPrivateRoutes(app fiber.Router, db *mongo.Database, publisher *broker.KafkaProducer, eventTopic string) {

	route := app.Group("/api/v1", middleware.JWTMiddleware())

	feedRepo := feedRepository.NewMongoFeedRepository(db)
	userRepo := userRepository.NewMongoUserRepository(db)
	feedService := feedUseCase.NewFeedUsecase(feedRepo, userRepo, publisher, eventTopic)
	feedHTTPHandler := feedHandler.NewHttpFeedHandler(feedService)

	userService := userUseCase.NewUserServiceWithPublisher(userRepo, publisher)
	userHandler := userHandler.NewHttpUserHandler(userService)

	route.Post("/posts", feedHTTPHandler.CreatePost)
	route.Post("/posts/:id/like", feedHTTPHandler.LikePost)
	route.Post("/posts/:id/comment", feedHTTPHandler.CommentOnPost)
	route.Get("/feed", feedHTTPHandler.GetFeed)

	route.Get("/me", userHandler.GetUser)
	route.Post("/users/:id/follow", userHandler.FollowUser)
	route.Delete("/users/:id/follow", userHandler.UnfollowUser)

}
