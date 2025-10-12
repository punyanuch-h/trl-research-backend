package database

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var (
	FirebaseApp     *firebase.App
	FirestoreClient *firestore.Client
	Ctx             context.Context
)

func InitFirebase(serviceAccountPath string) {
	Ctx = context.Background()
	fmt.Println("Ctx", Ctx)
	fmt.Println("serviceAccountPath", serviceAccountPath)
	opt := option.WithCredentialsFile(serviceAccountPath)
	fmt.Println("opt", opt)

	app, err := firebase.NewApp(Ctx, nil, opt)
	fmt.Println("app", app)
	if err != nil {
		log.Fatalf("❌ Firebase init error: %v", err)
	}

	client, err := app.Firestore(Ctx)
	fmt.Println("client", client)
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
