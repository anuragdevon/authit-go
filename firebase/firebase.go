package firebase_conn

import (
	"context"
	"path/filepath"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func FirebaseInit() (*firebase.App, error){
	// Firebase initialization
	// ctx := context.Background()
	// app, err := firebase.NewApp(ctx, nil)
	// if err != nil {
	// 	log.Fatalf("error initializing app: %v\n", err)
	// }
	serviceAccountKeyFilePath, err := filepath.Abs("./serviceAccountKey.json")
	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}
	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)
	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context.Background(), nil, opt)
	// if err != nil {
	// 	panic("Firebase load error")
	// }

	return app, err
}
