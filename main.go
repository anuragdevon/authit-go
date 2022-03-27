package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// Server route initialization
	router := gin.Default()

	// routes defination
	router.GET("/user/signup/", ??.Signup)
	router.GET("/user/signin/", ??.Signin)
	router.POST("/user/get/", ??.Userget)
	// start the server
	router.Run(":8080")

}
