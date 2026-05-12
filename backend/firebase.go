package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	opt := option.WithCredentialsFile("notifpref.json")
	fmt.Print(opt)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatal(err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatal(err)
	}

	token := "DEVICE_FCM_TOKEN"

	msg := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: "Hello from Linux",
			Body:  "Push notification working!",
		},
	}

	resp, err := client.Send(ctx, msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Sent message:", resp)
}
