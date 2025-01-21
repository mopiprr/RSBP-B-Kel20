package service

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

type (
	FirebaseService interface {
		SendToToken(app *firebase.App) error
	}

	firebaseService struct {
	}
)

func NewFirebaseService() FirebaseService {
	return &firebaseService{}
}

func (s *firebaseService) SendToToken(app *firebase.App) error {
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	registrationToken := "d3adVu0tMDyBUTPkgc_l-0:APA91bHjb6-wWkT1ABGSasFqxrsOR3AdfcTjLc8b7f7yukWLt32GS4UA5XdIwZ8p98oOLp-CBcyuYaCYdEPRji_f2WSXO9JKb7XPjotm_3bdkk-7hJyxJS8JuUHt82xzGGJ6Aacy0QWb"

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Notification Test",
			Body:  "Hello React!!",
		},
		Token: registrationToken,
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Successfully sent message:", response)

	return nil
}
