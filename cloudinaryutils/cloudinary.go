package cloudinaryutils

import (
	"context"
	"fmt"
	"os"
	"upload-easy/connection"

	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadToCloudianary(file *os.File, Directory string) error {
	//initialize cloudinary
	cld, err := connection.GetCloudinary()
	if err != nil {
		return fmt.Errorf("failed to Initailize Cloudinary: %v", err)
	}

	//upload file to cloudinary
	uploadResult, err := cld.Upload.Upload(context.TODO(), file, uploader.UploadParams{
		Folder: func() string {
			if Directory != "" {
				return Directory
			}
			return "golang"
		}(),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to Cloudinary: %v", err)
	}

	fmt.Printf("File uploaded successfully! URL: %s\n", uploadResult.SecureURL)
	return nil
}

func CreateFolderCloudinary(folderName string) error {
	cld, err := connection.GetCloudinary()
	if err != nil {
		return fmt.Errorf("failed to Initailize Cloudinary: %v", err)

	}
	_, err = cld.Admin.CreateFolder(context.TODO(), admin.CreateFolderParams{Folder: folderName})
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	fmt.Println("Folder created successfully", folderName)
	return nil
}
