package routes

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"

	// User
	userHandler "Notification_Preferences/internal/user/handler/rest"
	userRepository "Notification_Preferences/internal/user/repository"
	userUseCase "Notification_Preferences/internal/user/usecase"
)

func RegisterPublicRoutes(app fiber.Router, db *mongo.Database) {

	api := app.Group("/api/v1")

	// === Dependency Wiring ===

	// User
	userRepo := userRepository.NewGormUserRepository(db)
	userService := userUseCase.NewUserService(userRepo)
	userHandler := userHandler.NewHttpUserHandler(userService)

	// === Public Routes ===

	// Auth routes (separated from /users)
	authGroup := api.Group("/auth")
	authGroup.Post("/signup", userHandler.Register)
	authGroup.Post("/signin", userHandler.Login)

	// User routes
	userGroup := api.Group("/users")
	userGroup.Get("/", userHandler.FindAllUsers)
	userGroup.Get("/:id", userHandler.FindUserByID)
	userGroup.Patch("/:id", userHandler.PatchUser)
	userGroup.Delete("/:id", userHandler.DeleteUser)

}
