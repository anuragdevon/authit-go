package main

import (
	"firebase_go_auth/api"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error in loading env!")
	}
	// Server route initialization
	router := gin.Default()

	// routes defination
	router.POST("/user/signup/", api.UserSignUp)
	// router.GET("/user/signin/", api.UserSignIn)
	// router.POST("/user/get/", api.UserGet)

	router.Run(":8080")
}
