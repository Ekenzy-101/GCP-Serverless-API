package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	AccessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
	AppEnv            = os.Getenv("APP_ENV")
	ProjectID         = os.Getenv("PROJECT_ID")
	PostsCollection   = fmt.Sprintf("serverless-%v-posts", AppEnv)
	UsersCollection   = fmt.Sprintf("serverless-%v-users", AppEnv)
)

var (
	Port string
)

func LoadEnvironmentalVariables(filenames ...string) error {
	if err := godotenv.Load(filenames...); err != nil {
		return err
	}

	AccessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
	AppEnv = os.Getenv("APP_ENV")
	ProjectID = os.Getenv("PROJECT_ID")
	PostsCollection = fmt.Sprintf("serverless-%v-posts", AppEnv)
	UsersCollection = fmt.Sprintf("serverless-%v-users", AppEnv)

	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "5000"
	}

	return nil
}
