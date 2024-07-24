package main

import (
	"go-jwt/controllers"
	"go-jwt/initializers"
	"go-jwt/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	router := gin.Default()

	// router.SetTrustedProxies([]string{"172.17.0.2"})

	router.POST("/sign-up", controllers.Signup)
	router.POST("/login", controllers.Login)
	router.GET("/validate", middleware.RequierAuth, controllers.Validate)
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusAccepted, gin.H{
			"message": "I am Working",
		})
	})
	router.Run()
}
