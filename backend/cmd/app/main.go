package main

import "Notification_Preferences/internal/app"

func main() {

	app.Start() // Call server.go
}

// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"log"
// 	"net/http"

// 	firebase "firebase.google.com/go/v4"
// 	"firebase.google.com/go/v4/messaging"
// 	"google.golang.org/api/option"
// )

// type TokenRequest struct {
// 	Token    string `json:"token"`
// 	Platform string `json:"platform"`
// }

// var client *messaging.Client

// var savedTokens []string

// func enableCors(w http.ResponseWriter) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// 	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// }

// func main() {

// 	ctx := context.Background()

// 	opt := option.WithCredentialsFile("serviceAccount.json")

// 	app, err := firebase.NewApp(ctx, nil, opt)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	client, err = app.Messaging(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	http.HandleFunc("/device-token", saveToken)
// 	http.HandleFunc("/send", sendPush)

// 	log.Println("Backend running on :8080")

// 	err = http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func saveToken(w http.ResponseWriter, r *http.Request) {

// 	enableCors(w)

// 	if r.Method == "OPTIONS" {
// 		return
// 	}

// 	var req TokenRequest

// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		http.Error(w, err.Error(), 400)
// 		return
// 	}

// 	log.Println("TOKEN RECEIVED:")
// 	log.Println(req.Token)

// 	savedTokens = append(savedTokens, req.Token)

// 	w.Write([]byte("token saved"))
// }

// func sendPush(w http.ResponseWriter, r *http.Request) {

// 	enableCors(w)

// 	if r.Method == "OPTIONS" {
// 		return
// 	}

// 	ctx := context.Background()

// 	log.Println("TOKENS COUNT:", len(savedTokens))

// 	for _, token := range savedTokens {

// 		log.Println("SENDING TO TOKEN:")
// 		log.Println(token)

// 		msg := &messaging.Message{
// 			Token: token,
// 			Notification: &messaging.Notification{
// 				Title: "Hello from Go 🚀",
// 				Body:  "Push notification working successfully",
// 			},
// 		}

// 		resp, err := client.Send(ctx, msg)
// 		if err != nil {
// 			log.Println("FCM ERROR:", err)
// 			continue
// 		}

// 		log.Println("FCM SUCCESS:", resp)
// 	}

// 	w.Write([]byte("notifications attempted"))
// }
