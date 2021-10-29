package main

import (
	"log"
	"net/http"

	function "github.com/Ekenzy-101/GCP-Serverless"
	"github.com/Ekenzy-101/GCP-Serverless/config"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := startApplication(); err != nil {
		log.Fatal(err)
	}
}

func startApplication() error {
	err := config.LoadEnvironmentalVariables()
	if err != nil {
		return err
	}

	router := gin.Default()

	authRouter := router.Group("/auth")
	authRouter.POST("/login", withHandler(function.Login))
	authRouter.POST("/logout", withHandler(function.Logout))
	authRouter.POST("/register", withHandler(function.Register))

	postRouter := router.Group("/posts")
	postRouter.DELETE("", withHandler(function.DeletePost))
	postRouter.GET("", func(c *gin.Context) {
		c.JSON(log.Flags(), 1)
		if c.Query("id") == "" {
			function.GetPosts(c.Writer, c.Request)
			return
		}

		function.GetPost(c.Writer, c.Request)
	})
	postRouter.PUT("", withHandler(function.UpdatePost))

	return router.Run(":" + config.Port)
}

func withHandler(fn func(w http.ResponseWriter, r *http.Request)) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(c.Writer, c.Request)
	}
}
