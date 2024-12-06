package googleutils

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
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// drive
func UploadToDrive(file *os.File) error {
	// initialize google drive
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
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	//upload file
	if err := uploadFile(srv, file); err != nil {
		return fmt.Errorf("unable to upload file to Drive: %v", err)
	}

	fmt.Printf("File uploaded successfully: %s\n", file.Name())
	return nil
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

func uploadFile(srv *drive.Service, file *os.File) error {
	buf := make([]byte, 512)
	_, err := file.Read(buf)
	if err != nil {
		log.Fatal("Unable to read file for content type detection", err)
		return err
	}

	contentType := http.DetectContentType(buf)

	_, err = file.Seek(0, 0)
	if err != nil {
		log.Fatal("Unable to reset file pointer", err)
		return err
	}

	fileMetaData := &drive.File{
		Name: file.Name(),
	}
	_, err = srv.Files.Create(fileMetaData).Media(file, googleapi.ContentType(contentType)).Do()
	if err != nil {
		log.Fatal("Unable to upload file", err)
	}
	return nil
}
