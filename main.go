package main

import (
	"firebase_go_auth/api"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
		}
		c.Next()
	}
}

func main() {
	err := godotenv.Load()
	port := os.Getenv("PORT")
	if err != nil {
		log.Fatal("Error in loading env!")
	}
	// Server route initialization
	router := gin.Default()

	router.Use(corsMiddleware())

	// routes defination
	router.POST("/user/signup", api.UserSignUp)
	router.POST("/user/signin", api.UserSignIn)
	router.POST("/user/get", api.UserGet)

	router.Run(":" + port)
}
