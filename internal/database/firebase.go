package database

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	"cloud.google.com/go/firestore"
)

var (
	FirebaseApp     *firebase.App
	FirestoreClient *firestore.Client
	Ctx             context.Context
)

func InitFirebase(serviceAccountPath string) {
	Ctx = context.Background()
	opt := option.WithCredentialsFile(serviceAccountPath)

	app, err := firebase.NewApp(Ctx, nil, opt)
	if err != nil {
		log.Fatalf("❌ Firebase init error: %v", err)
	}

	client, err := app.Firestore(Ctx)
	if err != nil {
		log.Fatalf("❌ Firestore client init error: %v", err)
	}

	FirebaseApp = app
	FirestoreClient = client
	log.Println("✅ Firebase Firestore initialized")
}

func CloseFirebase() {
	if FirestoreClient != nil {
		if err := FirestoreClient.Close(); err != nil {
			log.Printf("⚠️ Error closing Firestore: %v", err)
		} else {
			log.Println("🛑 Firestore client closed")
		}
	}
}
