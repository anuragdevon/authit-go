package firebase_conn

import (
	"context"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"

	"firebase_go_auth/email"
)

func FirebaseInit() (context.Context, *auth.Client, error) {
	ctx := context.Background()

	serviceAccountKeyFilePath, err := filepath.Abs("./firebase/serviceAccountKey.json")
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
		return err
	}

	return email.SendMail(emailID, link)

}
