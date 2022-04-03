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
	// Image   string `json:"image"`Tuesday (5th April)
}

type loginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserGetData struct {
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
}

func UserSignUp(c *gin.Context) {

	// Extract Input
	var input newUser
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{"error_message": "Invalid Input!"})
		return
	}

	ctx, client, err := firebase_conn.FirebaseInit()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "Internal Server Error"})
		return
	}

	// Check if user already exists
	_, err = client.GetUserByEmail(ctx, input.Email)
	if err == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error_message": "User already exists!"})
		return
	}

	defaultPhotoUrl := "someImageName.png"

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
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "Unable to create user!"})
		return
	}

	// Send Email Verification
	err = firebase_conn.EmailVerification(input.Email, client, ctx)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "Verification Email Sending Failed!"})
		return

	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"message": "User Created Successfully!"})
}

func UserSignIn(c *gin.Context) {
	// Extract Input
	var input loginUser
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{"error_message": "Invalid Input!"})
		return
	}

	ctx, client, err := firebase_conn.FirebaseInit()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "Internal Server Error!"})
		return
	}

	// Check if user exists
	user, err := client.GetUserByEmail(ctx, input.Email)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{"error_message": "No Such User Exists!"})
		return
	}

	// Check if user is verified
	if !user.EmailVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error_message": "Email Not Verified!"})
		return
	}

	resp, err := firebase_conn.SignInWithEmailPassword(input.Email, input.Password)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "Internal Server Error!"})
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		decoder := json.NewDecoder(resp.Body)

		var data map[string]interface{}

		err = decoder.Decode(&data)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error_message": "Internal Server Error!"})
			return
		}
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, gin.H{"id_token": data["idToken"], "refresh_token": data["refreshToken"]})

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error_message": "Invalid Credentials!"})
	}
}

func UserGet(c *gin.Context) {
	// Extract Input
	var input UserGetData
	var data map[string]interface{}
	var uid string

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{"error_message": "Invalid Input!"})
		return
	}

	ctx, client, err := firebase_conn.FirebaseInit()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "Internal Server Error!"})
		return
	}

	decoded_token, err := client.VerifyIDTokenAndCheckRevoked(ctx, input.IdToken)

	// Renew tokens, if tokenID is expired or revoked
	if err != nil {
		resp, err := firebase_conn.RenewTokens(input.RefreshToken)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error_message": "Internal Server Error!"})
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			decoder := json.NewDecoder(resp.Body)
			// Check for control flow here
			err = decoder.Decode(&data)
			if err != nil {
				log.Println(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error_message": "Internal Server Error!"})
				return
			}
			// Check if refresh token is valid and exit
			byte_uid, _ := json.Marshal(data["user_id"])
			uid = string(byte_uid)
			uid = uid[1 : len(uid)-1]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error_message": "Invalid Credentials!"})
			return
		}
	} else {
		uid = decoded_token.UID
		data = map[string]interface{}{
			"id_token":      input.IdToken,
			"refresh_token": input.RefreshToken,
		}
	}

	user, err := client.GetUser(ctx, uid)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "Unable to get user info!"})
	}
	c.JSON(http.StatusOK, gin.H{"id_token": data["id_token"], "refresh_token": data["refresh_token"], "email": user.Email, "name": user.DisplayName, "user_id": uid})
}
