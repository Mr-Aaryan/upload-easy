package googleutils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"upload-easy/connection"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)


func UploadToDrive(file *os.File, parentId string) error {
	srv, err := connection.GetGoogleDrive()
	if(err != nil){
		return err
	}
	buf := make([]byte, 512)
	_, err = file.Read(buf)
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

	fileData, err := file.Stat()
	if err != nil {
		log.Fatal("Error reading file info")
	}
	fileMetaData := &drive.File{
		Name: fileData.Name(),
		Parents: []string{parentId},
	}
	_, err = srv.Files.Create(fileMetaData).Media(file, googleapi.ContentType(contentType)).Do()
	if err != nil {
		log.Fatal("Unable to upload file", err)
	}
	fmt.Printf("File uploaded successfully: %s\n", fileData.Name())
	return nil
}

func CreateDriveFolder(folderName string, parentId string) (string, error){
	srv, err := connection.GetGoogleDrive()
	if(err != nil){
		return "", err
	}
	folderMetaData := &drive.File{
		Name: folderName,
		MimeType: "application/vnd.google-apps.folder",
	}
	if(parentId != ""){
		folderMetaData.Parents = []string{parentId}
	}
	folder, err := srv.Files.Create(folderMetaData).Do()
	if err != nil {
		return "", fmt.Errorf("unable to create folder: %v", err)
	}
	fmt.Println("folder created successfully. FolderId: ", folder.Id)
	return folder.Id, nil
}