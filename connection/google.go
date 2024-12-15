package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var (
	srv *drive.Service
)

func InitializeGoogleDrive() error {
	once.Do(func() {
		const scope = "https://www.googleapis.com/auth/drive.file"
		ctx := context.Background()

		b, err := os.ReadFile("./googleutils/credentials.json")
		if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
		}
		config, err := google.ConfigFromJSON(b, scope)
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
		}
		client := getClient(config)
		srv, initErr = drive.NewService(ctx, option.WithHTTPClient(client))
	})
	return initErr

}

func GetGoogleDrive() (*drive.Service, error) {
	if srv == nil {
		return nil, fmt.Errorf("google drive is not initialized. Call Initialize first")
	}
	return srv, nil
}

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "./googleutils/token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authUrl := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser, authorize the application, and paste the authorization code here:\n%v\n", authUrl)
	//maybe needed to ask for authcode
	var authCode string
	fmt.Print("Enter authorization code: ")
	fmt.Scan(&authCode)

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
