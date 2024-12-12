package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"upload-easy/cloudinaryutils"
	"upload-easy/connection"
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

type FileInfo struct {
	Path     string
	IsDir    bool
	Children []FileInfo
}

var recentDir = ""

func handleDir(filePath string) ([]FileInfo, error) {
	var files []FileInfo
	filePath = strings.Replace(filePath, "/*", "", 1)
	// fmt.Println(filePath)
	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		files = append(files, FileInfo{Path: path, IsDir: info.IsDir()})
		// fmt.Println(files)
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
	}
	return files, nil
}

func main() {
	filePath := flag.String("file", "", "Path to the file to be uploaded")
	googleUpload := flag.Bool("g", false, "Upload to Google Drive")
	cloudinaryUpload := flag.Bool("c", false, "Upload to Cloudinary")
	megaUpload := flag.Bool("m", false, "Upload to Mega")
	help := flag.Bool("help", false, "View help doc")

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

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *filePath == "" {
		fmt.Println("Error: Please provide a file path using --file flag")
		os.Exit(1)
	} else if f := strings.Contains(*filePath, "/*"); f {
		files, err := handleDir(*filePath)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(0)
		}
		for _, file := range files {
			fmt.Printf("Path: %s, IsDir: %v\n", file.Path, file.IsDir)
			if !file.IsDir {
				dir := path.Dir(file.Path)
				relativeDir := strings.TrimPrefix(dir, strings.Replace(*filePath, "/*", "", 1))
				if relativeDir == "." {
					relativeDir = "uploads" // Default to "uploads" for base files
				}
				err := uploadFunc(file.Path, *googleUpload, *cloudinaryUpload, *megaUpload, relativeDir)
				if err != nil {
					log.Fatalf("Upload failed: %v", err)
				}
			} else {
				recentDir = file.Path
				err := createFolderFunc(recentDir, *googleUpload, *cloudinaryUpload, *megaUpload)
				if err != nil {
					log.Fatalf("Folder creation failed: %v", err)
				}
			}
		}
		os.Exit(0)
	}

	// filePath := "./uploads/file.png"
	err := uploadFunc(*filePath, *googleUpload, *cloudinaryUpload, *megaUpload, "")
	if err != nil {
		log.Fatalf("Upload failed: %v", err)
	}
}

// common upload function
func uploadFunc(filePath string, googleUpload bool, cloudinaryUpload bool, megaUpload bool, Directory string) error {
	file_, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file_.Close()

	fmt.Println("File opened successfully:", file_.Name())

	switch {
	case cloudinaryUpload:
		CloudinaryUrl := goDotEnv("CLOUDINARY_URL")
		err := connection.InitializeCloudinary(CloudinaryUrl)
		if err != nil {
			return fmt.Errorf("failed to initialize cloudinary: %v", err)
		}
		if err = cloudinaryutils.UploadToCloudianary(file_, Directory); err != nil {
			return fmt.Errorf("failed to upload file to cloudinary: %v", err)
		}
	case googleUpload:
		if err := googleutils.UploadToDrive(file_); err != nil {
			return fmt.Errorf("failed to upload file to drive: %v", err)

		}
	case megaUpload:
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

// common folder create function
func createFolderFunc(folderName string, googleUpload bool, cloudinaryUpload bool, megaUpload bool) error {
	if folderName == "" {
		return fmt.Errorf("foldername is not defined: %v", folderName)
	}
	fmt.Println("folder:", folderName)
	switch {
	case cloudinaryUpload:
		CloudinaryUrl := goDotEnv("CLOUDINARY_URL")
		err := connection.InitializeCloudinary(CloudinaryUrl)
		if err != nil {
			return fmt.Errorf("failed to initialize cloudinary: %v", err)
		}
		err = cloudinaryutils.CreateFolderCloudinary(folderName)
		if err != nil {
			return fmt.Errorf("error creating folder in cloudinary: %v", err)
		}
	case googleUpload:
		return fmt.Errorf("option doesnt exist")

	case megaUpload:
		return fmt.Errorf("option doesnt exist")
	default:
		return fmt.Errorf("undefined flag")
	}
	return nil
}
