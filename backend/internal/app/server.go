package app

import (
	"Notification_Preferences/pkg/database"
	"Notification_Preferences/utils"
	"log"

	"github.com/joho/godotenv"
)

func loadEnv() {
	p := "../../.env"
	if err := godotenv.Load(p); err == nil {
		log.Println("Loaded env:", p)
		return
	}

	log.Println("⚠️ No env file found, using system env")
}
func Start() {
	loadEnv()
	// Setup dependencies: database and configuration
	db, cfg, err := SetupDependencies("dev")
	if err != nil {
		log.Fatalf("❌ Failed to setup dependencies: %v", err)
	}

	// Setup REST server
	restApp, err := SetupRestServer(db, cfg)
	if err != nil {
		log.Fatalf("❌ Failed to setup REST server: %v", err)
	}

	// Setup gRPC server
	// grpcServer, err := SetupGrpcServer(db, cfg)
	// if err != nil {
	// 	log.Fatalf("❌ Failed to setup gRPC server: %v", err)
	// }

	// Start REST and gRPC servers
	go utils.StartRestServer(restApp, cfg)
	// go utils.StartGrpcServer(grpcServer, cfg)

	// Graceful shutdown listener
	utils.WaitForShutdown([]func(){
		func() {
			log.Println("Shutting down REST server...")
			if err := restApp.Shutdown(); err != nil {
				log.Printf("Error shutting down REST server: %v", err)
			}
		},
		// func() {
		// 	log.Println("Shutting down gRPC server...")
		// 	grpcServer.GracefulStop()
		// },
		func() {
			if err := database.Close(); err != nil {
				log.Printf("Error closing DB: %v", err)
			}
		},
	})

}
