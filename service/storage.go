package service

import (
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/Ekenzy-101/GCP-Serverless/config"
)

func GeneratePresignedURL(object string) (string, error) {
	return storage.SignedURL(config.BucketName, object, &storage.SignedURLOptions{
		Expires:        time.Now().Add(15 * time.Minute),
		Method:         http.MethodPut,
		Scheme:         storage.SigningSchemeV4,
		GoogleAccessID: config.ServiceAccountEmail,
		PrivateKey:     []byte(config.ServiceAccountKey),
	})
}
