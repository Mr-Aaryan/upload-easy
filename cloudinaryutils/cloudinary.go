package cloudinaryutils

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadToCloudianary(file *os.File, CloudinaryUrl string) error {
	//initialize cloudinary
	cld, err := cloudinary.NewFromURL(CloudinaryUrl)
	if err != nil {
		return fmt.Errorf("failed to Initailize Cloudinary: %v", err)
		
	}

	//upload file to cloudinary
	uploadResult, err := cld.Upload.Upload(context.TODO(), file, uploader.UploadParams{
		Folder: "golang",
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to Cloudinary: %v", err)
	}

	fmt.Printf("File uploaded successfully! URL: %s\n", uploadResult.SecureURL)
	return nil
}
