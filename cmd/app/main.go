package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// import "Notification_Preferences/internal/app"

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
		log.Panicln("No .env file found")
	}

	uri := os.Getenv("MONGO_URI")

	fmt.Println(uri)
	// app.Start() // Call server.go
}
