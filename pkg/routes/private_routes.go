package routes

import (
	userHandler "notification-pref/internal/user/handler/rest"
	userRepository "notification-pref/internal/user/repository"
	userUseCase "notification-pref/internal/user/usecase"
	middleware "notification-pref/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterPrivateRoutes(app fiber.Router, db *gorm.DB) {

	route := app.Group("/api/v1", middleware.JWTMiddleware())

	userRepo := userRepository.NewGormUserRepository(db)
	userService := userUseCase.NewUserService(userRepo)
	userHandler := userHandler.NewHttpUserHandler(userService)

	route.Get("/me", userHandler.GetUser)

}
