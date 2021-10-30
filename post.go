package function

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/Ekenzy-101/GCP-Serverless/config"
	"github.com/Ekenzy-101/GCP-Serverless/helper"
	"github.com/Ekenzy-101/GCP-Serverless/model"
	"github.com/Ekenzy-101/GCP-Serverless/service"
	"github.com/Ekenzy-101/GCP-Serverless/types"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	value, err := helper.AuthorizeRequest(r)
	if err != nil {
		helper.SendJSONResponse(w, http.StatusUnauthorized, types.M{"message": err.Error()})
		return
	}

	claims, ok := value.(*jwt.RegisteredClaims)
	if !ok {
		helper.SendJSONResponse(w, http.StatusUnauthorized, types.M{"message": "Payload's format is invalid"})
		return
	}

	ctx := context.Background()
	client, err := service.CreateFirestoreClient(ctx)
	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	userDocument, err := client.Collection(config.UsersCollection).Doc(claims.Subject).Get(ctx)
	if status.Code(err) == codes.NotFound {
		helper.SendJSONResponse(w, http.StatusNotFound, types.M{"message": "User not found"})
		return
	}

	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	post := &model.Post{}
	if messages := helper.ValidateRequestBody(r, post); messages != nil {
		helper.SendJSONResponse(w, http.StatusBadRequest, messages)
		return
	}

	documentRef := client.Collection(config.PostsCollection).NewDoc()
	url, err := service.GeneratePresignedURL(fmt.Sprintf("posts/%v", documentRef.ID))
	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	post.Image = strings.Split(url, "?")[0]
	post.User = types.M{"id": claims.Subject, "name": userDocument.Data()["name"]}
	result, err := documentRef.Create(ctx, post)
	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	post = post.SetID(documentRef.ID).SetTimestamps(result.UpdateTime, result.UpdateTime)
	helper.SendJSONResponse(w, http.StatusOK, types.M{"post": post, "url": url})
}

func DeletePost(w http.ResponseWriter, r *http.Request) {

}

func GetPost(w http.ResponseWriter, r *http.Request) {
	postId := r.URL.Query().Get("id")
	ctx := context.Background()
	client, err := service.CreateFirestoreClient(ctx)
	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	document, err := client.Collection(config.PostsCollection).Doc(postId).Get(ctx)
	if status.Code(err) == codes.NotFound {
		helper.SendJSONResponse(w, http.StatusNotFound, types.M{"message": "Post with the given id does not exist"})
		return
	}

	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	post, err := model.NewPostFromDocument(document)
	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	post = post.SetID(postId).SetTimestamps(document.CreateTime, document.UpdateTime)
	helper.SendJSONResponse(w, http.StatusOK, types.M{"post": post})
}

func GetPosts(w http.ResponseWriter, r *http.Request) {

}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	value, err := helper.AuthorizeRequest(r)
	if err != nil {
		helper.SendJSONResponse(w, http.StatusUnauthorized, types.M{"message": err.Error()})
		return
	}

	claims, ok := value.(*jwt.RegisteredClaims)
	if !ok {
		helper.SendJSONResponse(w, http.StatusUnauthorized, types.M{"message": "Payload's format is invalid"})
		return
	}

	requestBody := &model.Post{}
	if messages := helper.ValidateRequestBody(r, requestBody); messages != nil {
		helper.SendJSONResponse(w, http.StatusBadRequest, messages)
		return
	}

	postId := r.URL.Query().Get("id")
	ctx := context.Background()
	client, err := service.CreateFirestoreClient(ctx)
	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	documentRef := client.Collection(config.PostsCollection).Doc(postId)
	document, err := documentRef.Get(ctx)
	if status.Code(err) == codes.NotFound {
		helper.SendJSONResponse(w, http.StatusNotFound, types.M{"message": "Post with the given id does not exist"})
		return
	}

	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	post, err := model.NewPostFromDocument(document)
	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	if post.User["id"] != claims.Subject {
		helper.SendJSONResponse(w, http.StatusForbidden, types.M{"message": "You are not allowed to update this post"})
		return
	}

	updates := []firestore.Update{
		{Path: "title", Value: requestBody.Title},
		{Path: "content", Value: requestBody.Content},
	}
	result, err := documentRef.Update(ctx, updates)
	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	post = post.SetContent(requestBody.Content).SetID(postId).
		SetTimestamps(document.CreateTime, result.UpdateTime).SetTitle(requestBody.Title)
	helper.SendJSONResponse(w, http.StatusOK, types.M{"post": post})
}
