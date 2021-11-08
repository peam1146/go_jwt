package main

import (
	// gin
	"fmt"
	"os"

	"main.go/controller"
	"main.go/databases"
	"main.go/middlewares"
	"main.go/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// init database
	databases.InitDatabase()
	authServices := services.NewAuthServices()
	authController := controller.NewAuthController(&authServices)

	r := gin.Default()

	r.POST("/register", authController.Register)
	r.GET("/login", authController.Login)
	// authentication middleware group
	auth := r.Group("/auth")
	auth.Use(middlewares.JWTMiddleware())
	{
		auth.GET("/user", authController.User)
	}

	// run server with auto tls
	// r.RunTLS(":443", "./cert/server.crt", "./cert/server.key")

	r.Run(":8080")
}
