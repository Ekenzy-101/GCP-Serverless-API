package function

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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

	doc, err := client.Collection(config.UsersCollection).Doc(claims.Subject).Get(ctx)
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

	docRef := client.Collection(config.PostsCollection).NewDoc()
	url, err := service.GeneratePresignedURL(fmt.Sprintf("posts/%v", docRef.ID))
	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	post.Image = strings.Split(url, "?")[0]
	post.User = types.M{"id": claims.Subject, "name": doc.Data()["name"]}
	_, err = docRef.Create(ctx, post)
	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	post.ID = docRef.ID
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

	doc, err := client.Collection(config.PostsCollection).Doc(postId).Get(ctx)
	if status.Code(err) == codes.NotFound {
		helper.SendJSONResponse(w, http.StatusNotFound, types.M{"message": "Post with the given id does not exist"})
		return
	}

	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	post := doc.Data()
	post["id"] = postId
	helper.SendJSONResponse(w, http.StatusOK, types.M{"post": post})
}

func GetPosts(w http.ResponseWriter, r *http.Request) {

}

func UpdatePost(w http.ResponseWriter, r *http.Request) {

}
