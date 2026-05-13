package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

func InitFCM() *messaging.Client {
	ctx := context.Background()

	opt := option.WithCredentialsFile("notifpref_firebase.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("firebase init failed: %v", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("fcm client init failed: %v", err)
	}

	log.Println("FCM initialized successfully")

	return client
}
