package routes

import (
	userHandler "Notification_Preferences/internal/user/handler/rest"
	userRepository "Notification_Preferences/internal/user/repository"
	userUseCase "Notification_Preferences/internal/user/usecase"
	middleware "Notification_Preferences/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterPrivateRoutes(app fiber.Router, db *mongo.Database, publisher userUseCase.EventPublisher) {

	route := app.Group("/api/v1", middleware.JWTMiddleware())

	userRepo := userRepository.NewMongoUserRepository(db)
	userService := userUseCase.NewUserServiceWithPublisher(userRepo, publisher)
	userHandler := userHandler.NewHttpUserHandler(userService)

	route.Get("/me", userHandler.GetUser)
	route.Post("/users/:id/follow", userHandler.FollowUser)
	route.Delete("/users/:id/follow", userHandler.UnfollowUser)

}
