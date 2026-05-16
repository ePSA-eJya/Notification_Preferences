package app

import (
	"log"
	"os"

	"Notification_Preferences/internal/broker"
	"Notification_Preferences/pkg/config"
	"Notification_Preferences/pkg/database"
	"Notification_Preferences/pkg/middleware"
	"Notification_Preferences/pkg/routes"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// rest
// func SetupRestServer(db *mongo.Database, cfg *config.Config) (*fiber.App, error) {
// 	app := fiber.New()
// 	middleware.FiberMiddleware(app)
// 	// comment out Swagger when testing
// 	// routes.SwaggerRoute(app)
// 	routes.RegisterPublicRoutes(app, db)
// 	routes.RegisterPrivateRoutes(app, db)
// 	routes.RegisterNotFoundRoute(app)
// 	return app, nil
// }

func SetupRestServer(db *mongo.Database, cfg *config.Config) (*fiber.App, error) {
	app := fiber.New()

	var publisher *broker.KafkaProducer
	if cfg != nil && cfg.KafkaBrokerURL != "" {
		publisher = broker.NewKafkaProducer(cfg.KafkaBrokerURL, cfg.KafkaEventTopic)
	}

	middleware.FiberMiddleware(app)
	
	// Create uploads directory relative to current working directory
	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		log.Printf("Warning: Failed to create uploads directory: %v", err)
	}
	
	// Serve uploads directory as static files
	// The /uploads prefix will map to the uploads directory
	app.Static("/uploads", uploadsDir)
	log.Printf("Static files configured. Uploads directory: %s", uploadsDir)
	
	routes.RegisterPublicRoutes(app, db, publisher)
	routes.RegisterPrivateRoutes(app, db, publisher, cfg)
	routes.RegisterNotFoundRoute(app)

	return app, nil
}

// grpc
// func SetupGrpcServer(db *mongo.Database, cfg *config.Config) (*grpc.Server, error) {
// 	s := grpc.NewServer()
// 	orderRepo := orderRepository.NewGormOrderRepository(db)
// 	orderService := orderUseCase.NewOrderService(orderRepo)

// 	orderHandler := GrpcOrderHandler.NewGrpcOrderHandler(orderService)
// 	orderpb.RegisterOrderServiceServer(s, orderHandler)
// 	return s, nil
// }

// dependencies

func SetupDependencies(env string) (*mongo.Database, *config.Config, error) {
	cfg := config.LoadConfig(env)

	log.Println("MONGODB_URI =", os.Getenv("MongoURI"))
	log.Println("db_name =", os.Getenv("DB_NAME"))

	db, err := database.ConnectMongo(cfg.MongoURI, cfg.DB_NAME)

	if err != nil {
		return nil, nil, err
	}

	return db, cfg, nil
}
