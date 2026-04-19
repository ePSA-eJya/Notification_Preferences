package app

import (
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"notification-pref/internal/entities"
	GrpcOrderHandler "notification-pref/internal/order/handler/grpc"
	orderRepository "notification-pref/internal/order/repository"
	orderUseCase "notification-pref/internal/order/usecase"
	"notification-pref/pkg/config"
	"notification-pref/pkg/database"
	"notification-pref/pkg/middleware"
	"notification-pref/pkg/routes"
	orderpb "notification-pref/proto/order"
)

// rest
func SetupRestServer(db *gorm.DB, cfg *config.Config) (*fiber.App, error) {
	app := fiber.New()
	middleware.FiberMiddleware(app)
	// comment out Swagger when testing
	// routes.SwaggerRoute(app)
	routes.RegisterPublicRoutes(app, db)
	routes.RegisterPrivateRoutes(app, db)
	routes.RegisterNotFoundRoute(app)
	return app, nil
}

// grpc
func SetupGrpcServer(db *gorm.DB, cfg *config.Config) (*grpc.Server, error) {
	s := grpc.NewServer()
	orderRepo := orderRepository.NewGormOrderRepository(db)
	orderService := orderUseCase.NewOrderService(orderRepo)

	orderHandler := GrpcOrderHandler.NewGrpcOrderHandler(orderService)
	orderpb.RegisterOrderServiceServer(s, orderHandler)
	return s, nil
}

// dependencies
func SetupDependencies(env string) (*gorm.DB, *config.Config, error) {
	cfg := config.LoadConfig(env)

	db, err := database.Connect(cfg.DatabaseDSN)
	if err != nil {
		return nil, nil, err
	}

	if env == "test" {
		_ = db.Migrator().DropTable(&entities.Order{}, &entities.User{})
	}
	if err := db.AutoMigrate(&entities.Order{}, &entities.User{}); err != nil {
		return nil, nil, err
	}

	return db, cfg, nil
}
