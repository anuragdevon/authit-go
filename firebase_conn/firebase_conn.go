package firebase_conn

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"

	"firebase_go_auth/email"
	"firebase_go_auth/utils"
)

func FirebaseInit() (context.Context, *auth.Client, error) {
	ctx := context.Background()

	serviceAccountKeyFilePath, err := filepath.Abs("./serviceAccountKey.json")
	if err != nil {
		log.Println("Unable to load Firebase Cred Files!")
		return ctx, nil, err
	}

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println("Unable to initialize Firebase App!")
		return ctx, nil, err
	}

	client, err := app.Auth(context.Background())

	if err != nil {
		log.Println("Unable to initialize Firebase Auth Client!")
		return ctx, nil, err
	}

	return ctx, client, err
}

func EmailVerification(emailID string, client *auth.Client, ctx context.Context) error {
	link, err := client.EmailVerificationLinkWithSettings(ctx, emailID, nil)
	if err != nil {
		log.Println("Error while generating email verification link: ", err)
		return err
	}

	return email.SendMail(emailID, link)
}

func SignInWithEmailPassword(email string, password string) (*http.Response, error) {

	API_KEY := os.Getenv("API_KEY")
	endpoint := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key={API_KEY}"

	replacer := strings.NewReplacer("{API_KEY}", API_KEY)
	endpoint = replacer.Replace(endpoint)

	payload := map[string]interface{}{"email": email, "password": password, "returnSecureToken": true}

	resp, err := utils.InternalRequest(payload, "POST", endpoint)

	return resp, err
}

func RenewTokens(RefreshToken string) (*http.Response, error) {
	API_KEY := os.Getenv("API_KEY")
	endpoint := "https://securetoken.googleapis.com/v1/token?key={API_KEY}"

	replacer := strings.NewReplacer("{API_KEY}", API_KEY)
	endpoint = replacer.Replace(endpoint)

	payload := map[string]interface{}{"grant_type": "refresh_token", "refresh_token": RefreshToken}

	resp, err := utils.InternalRequest(payload, "POST", endpoint)

	return resp, err
}
