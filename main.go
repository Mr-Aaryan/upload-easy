package main

import (
	"fmt"
	"log"
	"os"

	// "github.com/gorilla/mux"

	"learn3/cloudinaryutils"
	"learn3/googleutils"
	"learn3/megautils"

	"github.com/joho/godotenv"
)

func goDotEnv(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func main() {
	uploadFunc()
}

// common function
func uploadFunc() error {

	filePath := "./uploads/file.png"
	file_, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file_.Close()

	fmt.Println("File opened successfully:", file_.Name())

	fmt.Println("Where do you want to upload?")
	fmt.Println("1. Cloudinary\n2. Google\n3. Mega")
	var choise int
	fmt.Scan(&choise)

	switch choise {
	case 1:
		CloudinaryUrl := goDotEnv("CLOUDINARY_URL")
		if err := cloudinaryutils.UploadToCloudianary(file_, CloudinaryUrl); err != nil {
			return fmt.Errorf("failed to upload file to cloudinary: %v", err)
		}
	case 2:
		if err := googleutils.UploadToDrive(file_); err != nil {
			return fmt.Errorf("failed to upload file to drive: %v", err)

		}
	case 3:
		MegaEmail := goDotEnv("MEGA_EMAIL")
		MegaPassword := goDotEnv("MEGA_PASSWORD")
		if err := megautils.UploadToMega(file_, MegaEmail, MegaPassword); err != nil {
			return fmt.Errorf("failed to upload file to mega: %v", err)
		}
	default:
		fmt.Println("Invalid choise")
		return nil
	}
	return nil
}
