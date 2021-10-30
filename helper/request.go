package helper

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Ekenzy-101/GCP-Serverless/config"
	"github.com/Ekenzy-101/GCP-Serverless/service"
	"github.com/Ekenzy-101/GCP-Serverless/types"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

// Obj should be a pointer to a value
func ValidateRequestBody(r *http.Request, obj interface{}) interface{} {
	err := json.NewDecoder(r.Body).Decode(obj)
	if err != nil {
		return types.M{"message": err.Error()}
	}

	err = validate.Struct(obj)
	validationErrors := validator.ValidationErrors{}
	if errors.As(err, &validationErrors) {
		return GenerateErrorMessages(validationErrors)
	}

	if err != nil {
		return types.M{"message": err.Error()}
	}

	return nil
}

func AuthorizeRequest(r *http.Request) (jwt.Claims, error) {
	cookie, err := r.Cookie(config.AccessTokenCookieName)
	if err != nil {
		return nil, errors.New("Cookie not found")
	}

	return service.VerifyJWTToken(service.JWTOptions{
		Claims: &jwt.RegisteredClaims{},
		Secret: config.AccessTokenSecret,
		Token:  cookie.Value,
	})
}

func SendJSONResponse(w http.ResponseWriter, statusCode int, obj interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if obj != nil {
		if err := json.NewEncoder(w).Encode(obj); err != nil {
			panic(err)
		}
	}
}
