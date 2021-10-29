package service

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type JWTOptions struct {
	SigningMethod jwt.SigningMethod
	Claims        jwt.Claims
	Secret        string
	Token         string
}

func SignJWTToken(options JWTOptions) (string, error) {
	token := jwt.NewWithClaims(options.SigningMethod, options.Claims)
	signedToken, err := token.SignedString([]byte(options.Secret))

	return signedToken, err
}

func VerifyJWTToken(options JWTOptions) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(
		options.Token,
		options.Claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(options.Secret), nil
		},
	)

	if err != nil {
		return nil, errors.New("Invalid or expired token")
	}

	return token.Claims, nil
}
