package main

import (
	"firebase_go_auth/api"
	"log"
	"time"

	"github.com/gin-contrib/cors"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:8080"},
		AllowMethods:     []string{"PUT", "POST", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
	// routes defination
	router.POST("/user/signup/", api.UserSignUp)
	router.POST("/user/signin/", api.UserSignIn)
	// router.POST("/user/get/", api.UserGet)

	router.Run(":8080")
}
