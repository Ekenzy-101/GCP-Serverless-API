package function

import (
	"context"
	"net/http"
	"time"

	"github.com/Ekenzy-101/GCP-Serverless/config"
	"github.com/Ekenzy-101/GCP-Serverless/helper"
	"github.com/Ekenzy-101/GCP-Serverless/model"
	"github.com/Ekenzy-101/GCP-Serverless/service"
	"github.com/Ekenzy-101/GCP-Serverless/types"
	"github.com/golang-jwt/jwt/v4"
)

func Login(w http.ResponseWriter, r *http.Request) {

}

func Logout(w http.ResponseWriter, r *http.Request) {

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
	err := user.HashAndSetPassword()
	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	ctx := context.Background()
	client, err := service.CreateFirestoreClient(ctx)
	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	docRef, _, err := client.Collection(config.UsersCollection).Add(ctx, user)
	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	user.ID = docRef.ID
	user.Password = ""
	cliams := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.AccessTokenTTLInSeconds * time.Second)),
		Subject:   user.ID,
	}
	accessToken, err := service.SignJWTToken(service.JWTOptions{
		SigningMethod: jwt.SigningMethodHS256,
		Secret:        config.AccessTokenSecret,
		Claims:        cliams,
	})

	if err != nil {
		helper.SendJSONResponse(w, http.StatusInternalServerError, types.M{"message": err.Error()})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     config.AccessTokenCookieName,
		MaxAge:   config.AccessTokenTTLInSeconds,
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true,
	})
	helper.SendJSONResponse(w, http.StatusOK, types.M{"user": user})
}
