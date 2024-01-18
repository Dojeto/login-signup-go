package main

import (
	"fmt"
	"go-jwt/controllers"
	"go-jwt/initializers"
	"go-jwt/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	fmt.Println("test 2")
	router := gin.Default()

	router.POST("/sign-up", controllers.Signup)
	router.POST("/login", controllers.Login)
	router.GET("/validate", middleware.RequierAuth, controllers.Validate)
	router.Run()
}
