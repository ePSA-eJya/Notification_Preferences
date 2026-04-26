package app

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"Notification_Preferences/pkg/config"
	"Notification_Preferences/pkg/database"
	"Notification_Preferences/pkg/middleware"
	"Notification_Preferences/pkg/routes"

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

	middleware.FiberMiddleware(app)
	routes.RegisterPublicRoutes(app, db)
	routes.RegisterPrivateRoutes(app, db)
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

	// ctx := context.Background()

	// fcmClient, err := config.InitFCM(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Println("MONGODB_URI =", os.Getenv("MongoURI"))

	db, err := database.ConnectMongo(cfg.MongoURI, cfg.DBName)
	if err != nil {
		return nil, nil, err
	}

	return db, cfg, nil
}
