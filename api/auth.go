package api

import (
	"encoding/json"
	"log"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"

	firebase_conn "firebase_go_auth/firebase_conn"
)

type newUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`

	// Should be given as base64 encoded image
	// Image   string `json:"image"`
}

type loginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserSignUp(c *gin.Context) {

	// Extract Input
	var input newUser
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{"response": "Invalid Input!"})
		return
	}

	ctx, client, err := firebase_conn.FirebaseInit()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"response": "Internal Server Error!"})
		return
	}

	// Check if user already exists
	_, err = client.GetUserByEmail(ctx, input.Email)
	if err == nil {
		log.Println(err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{"response": "User already exists!"})
		return
	}

	defaultPhotoUrl := ""

	params := (&auth.UserToCreate{}).
		Email(input.Email).
		EmailVerified(false).
		Password(input.Password).
		DisplayName(input.Name).
		PhotoURL(defaultPhotoUrl).
		Disabled(false)

	_, err = client.CreateUser(ctx, params)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"response": "Unable to create user!"})
		return
	}

	// Send Email Verification
	err = firebase_conn.EmailVerification(input.Email, client, ctx)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"response": "Verification Email Sending Failed!"})
		return

	}

	c.JSON(http.StatusOK, gin.H{"response": "User Created Successfully!"})
}

func UserSignIn(c *gin.Context) {

	// Extract Input
	var input loginUser
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{"response": "Invalid Input!"})
		return
	}

	ctx, client, err := firebase_conn.FirebaseInit()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"response": "Internal Server Error!"})
		return
	}

	// Check if user exists
	user, err := client.GetUserByEmail(ctx, input.Email)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{"response": "No Such User Exists!"})
		return
	}

	// Check if user is verified
	if !user.EmailVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"response": "Email Not Verified!"})
		return
	}

	resp, err := firebase_conn.SignInWithEmailPassword(input.Email, input.Password)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"response": "Internal Server Error!"})
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		decoder := json.NewDecoder(resp.Body)

		var data map[string]interface{}

		err = decoder.Decode(&data)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"response": "Internal Server Error!"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"response": gin.H{"id_token": data["idToken"], "refresh_token": data["refreshToken"], "uid": data["localId"].(string), "email": data["email"].(string), "name": data["displayName"].(string)}})

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"response": "Invalid Credentials!"})
	}
}
