package firebase_conn

import (
	"context"
	"log"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"

	"firebase_go_auth/email"
)

func FirebaseInit() (context.Context, *auth.Client, error) {
	ctx := context.Background()

	serviceAccountKeyFilePath, err := filepath.Abs("./firebase_conn/serviceAccountKey.json")
	// serviceAccountKeyFilePath, err := filepath.Rel(".", "serviceAccountKey.json")
	if err != nil {
		panic("Unable to load Firebase Cred Files!")
	}
	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context.Background(), nil, opt)
	client, err := app.Auth(context.Background())
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

// func SignInWithEmailPasword() error {
// 	API_KEY := os.Getenv("API_KEY")
// 	r := strings.NewReplacer(API_KEY)
// 	endpoint := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key={API_KEY}"
// 	r.Replace(endpoint)
// 	payload := {"email":"[user@example.com]","password":"[PASSWORD]","returnSecureToken":true}
// 	err := utils.InternalRequest(payload, "POST", endpoint)

// 	return err

// }
