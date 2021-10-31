package service

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"github.com/Ekenzy-101/GCP-Serverless/config"
)

var (
	firestoreClient *firestore.Client
	storageClient   *storage.Client
)

// Initialize all services' clients. This will only run during an instance's cold start.
func init() {
	ctx := context.Background()
	var err error
	firestoreClient, err = firestore.NewClient(ctx, config.ProjectID)
	if err != nil {
		log.Fatalf("firestore.NewClient %v", err)
	}

	storageClient, err = storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient %v", err)
	}
}
