package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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
	googleUpload := flag.Bool("g", false, "Upload to Google Drive")
	cloudinaryUpload := flag.Bool("c", false, "Upload to Cloudinary")
	megaUpload := flag.Bool("m", false, "Upload to Mega")
	help := flag.Bool("help", false, "View help doc")

	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: upload-easy [options]

Options:
  --file <path> (required) Path to the file to be uploaded.
  -g                Upload to Google Drive.
  -m                Upload to Mega.
  -c                Upload to Cloudinary.
  --help            View this help documentation.

Example:
  upload-easy --file "./upload/file.png" -g
`)
	}

	if *help {
		flag.Usage()
		os.Exit(1)

	}

	if *filePath == "" {
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

	switch {
	case *cloudinaryUpload:
		CloudinaryUrl := goDotEnv("CLOUDINARY_URL")
		// CloudinaryUrl := "cloudinary://532797552144317:JaOdM6I1Ds5vXUveqHSROkM8s3I@dhshp6y6p"
		if err := cloudinaryutils.UploadToCloudianary(file_, CloudinaryUrl); err != nil {
			return fmt.Errorf("failed to upload file to cloudinary: %v", err)
		}
	case *googleUpload:
		if err := googleutils.UploadToDrive(file_); err != nil {
			return fmt.Errorf("failed to upload file to drive: %v", err)

		}
	case *megaUpload:
		MegaEmail := goDotEnv("MEGA_EMAIL")
		MegaPassword := goDotEnv("MEGA_PASSWORD")
		if err := megautils.UploadToMega(file_, MegaEmail, MegaPassword); err != nil {
			return fmt.Errorf("failed to upload file to mega: %v", err)
		}
	default:
		fmt.Println("Please specify a valid upload destination using -g, -c, or -m.")
	}
	return nil
}
