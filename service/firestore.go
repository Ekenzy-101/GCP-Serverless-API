package service

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/Ekenzy-101/GCP-Serverless/config"
)

func CreateFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	return firestore.NewClient(ctx, config.ProjectID)
}
