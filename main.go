package main

import (
	"firebase_go_auth/api"

	"github.com/gin-gonic/gin"
)

func main() {
	// Server route initialization
	router := gin.Default()

	// routes defination
	router.POST("/user/signup/", api.UserSignUp)
	// router.GET("/user/signin/", api.UserSignIn)
	// router.POST("/user/get/", api.UserGet)

	router.Run(":8080")
}
