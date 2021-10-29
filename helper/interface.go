package helper

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/Ekenzy-101/GCP-Serverless/types"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

var (
	lowercaseRegex        = regexp.MustCompile(`[a-z]+`)
	nameRegex             = regexp.MustCompile(`^[a-zA-z ]+$`)
	numberRegex           = regexp.MustCompile(`\d+`)
	specialCharacterRegex = regexp.MustCompile(`\W+`)
	uppercaseRegex        = regexp.MustCompile(`[A-Z]+`)
)

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(jsonTagName)
	err := validate.RegisterValidation("name", validateName)
	if err != nil {
		panic(err)
	}

	err = validate.RegisterValidation("password", validatePassword)
	if err != nil {
		panic(err)
	}
}

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

func SendJSONResponse(w http.ResponseWriter, statusCode int, obj interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		panic(err)
	}
}

func jsonTagName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}

func validateName(fl validator.FieldLevel) bool {
	return nameRegex.MatchString(fl.Field().String())
}

func validatePassword(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return specialCharacterRegex.MatchString(value) &&
		lowercaseRegex.MatchString(value) &&
		uppercaseRegex.MatchString(value) &&
		numberRegex.MatchString(value)
}
