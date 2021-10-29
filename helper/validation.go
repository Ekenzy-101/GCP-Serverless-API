package helper

import (
	"reflect"
	"regexp"
	"strings"

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
