package function

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Ekenzy-101/GCP-Serverless/config"
	"github.com/Ekenzy-101/GCP-Serverless/helper"
	"github.com/Ekenzy-101/GCP-Serverless/model"
	"github.com/Ekenzy-101/GCP-Serverless/service"
	"github.com/Ekenzy-101/GCP-Serverless/types"
	"google.golang.org/api/iterator"
)

func Login(w http.ResponseWriter, r *http.Request) {
	requestBody := &types.LoginRequestBody{}
	if messages := helper.ValidateRequestBody(r, requestBody); messages != nil {
		helper.SendJSONResponse(w, http.StatusBadRequest, messages)
		return
	}

	ctx := context.Background()
	client := service.GetFirestoreClient()
	iter := client.Collection(config.UsersCollection).Where("email", "==", requestBody.Email).Documents(ctx)
	document, err := iter.Next()
	if err != nil && !errors.Is(err, iterator.Done) {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	if document == nil {
		helper.SendJSONResponse(w, http.StatusBadRequest, types.M{"message": "Invalid email or password"})
		return
	}

	user := &model.User{}
	if err := document.DataTo(user); err != nil {
		helper.SendJSONResponse(w, http.StatusBadRequest, types.M{"message": err.Error()})
		return
	}

	matches, err := user.ComparePassword(requestBody.Password)
	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	if !matches {
		helper.SendJSONResponse(w, http.StatusBadRequest, types.M{"message": "Invalid email or password"})
		return
	}

	user.SetIDAndPassword(document.Ref.ID)
	accessToken, err := user.GenerateAccessToken()
	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		MaxAge:   config.AccessTokenTTLInSeconds,
		Name:     config.AccessTokenCookieName,
		Path:     "/",
		Secure:   true,
		Value:    accessToken,
	})
	helper.SendJSONResponse(w, http.StatusOK, types.M{"user": user})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		MaxAge:   -1,
		Name:     config.AccessTokenCookieName,
		Path:     "/",
		Secure:   true,
		Value:    "",
	})
	helper.SendJSONResponse(w, http.StatusNoContent, nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	requestBody := &types.RegisterRequestBody{}
	if messages := helper.ValidateRequestBody(r, requestBody); messages != nil {
		helper.SendJSONResponse(w, http.StatusBadRequest, messages)
		return
	}

	user := &model.User{
		Name:      requestBody.Name,
		Email:     requestBody.Email,
		Password:  requestBody.Password,
		CreatedAt: time.Now(),
	}
	if err := user.HashAndSetPassword(); err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	ctx := context.Background()
	client := service.GetFirestoreClient()
	iter := client.Collection(config.UsersCollection).Where("email", "==", requestBody.Email).Documents(ctx)
	document, err := iter.Next()
	if err != nil && !errors.Is(err, iterator.Done) {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	if document != nil {
		helper.SendJSONResponse(w, http.StatusBadRequest, types.M{"message": "A user with this email already exists"})
		return
	}

	documentRef, _, err := client.Collection(config.UsersCollection).Add(ctx, user)
	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	user.SetIDAndPassword(documentRef.ID)
	accessToken, err := user.GenerateAccessToken()
	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		MaxAge:   config.AccessTokenTTLInSeconds,
		Name:     config.AccessTokenCookieName,
		Path:     "/",
		Secure:   true,
		Value:    accessToken,
	})
	helper.SendJSONResponse(w, http.StatusOK, types.M{"user": user})
}
