package service

import (
	"cloud.google.com/go/firestore"
)

func GetFirestoreClient() *firestore.Client {
	return firestoreClient
}
