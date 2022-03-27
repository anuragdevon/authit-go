package firebase_conn

import (
	"context"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
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
