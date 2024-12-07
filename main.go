package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	// "github.com/gorilla/mux"

	"upload-easy/cloudinaryutils"
	"upload-easy/googleutils"
	"upload-easy/megautils"

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

	filePath := flag.String("file", "", "Path to the file to be uploaded")
	flag.Parse()

	if(*filePath == ""){
		fmt.Println("Error: Please provide a file path using --file flag")
		os.Exit(1)
	}

	// filePath := "./uploads/file.png"
	file_, err := os.Open(*filePath)
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
		// CloudinaryUrl := goDotEnv("CLOUDINARY_URL")
		CloudinaryUrl := "cloudinary://532797552144317:JaOdM6I1Ds5vXUveqHSROkM8s3I@dhshp6y6p"
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
