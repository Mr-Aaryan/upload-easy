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

	"upload-easy/utils"

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
	}

	// if its folder, this will trigger
	if info, err := os.Stat(*filePath); err == nil && info.IsDir() {
		err := processDirectory(*filePath, *googleUpload, *cloudinaryUpload, *megaUpload)
		if err != nil {
			log.Fatalf("Failed to process directory: %v", err)
		}
		os.Exit(0)
	}

	// it its file, this will trigger
	err := uploadFunc(*filePath, *googleUpload, *cloudinaryUpload, *megaUpload, "")
	if err != nil {
		log.Fatalf("Upload failed: %v", err)
	}
	os.Exit(0)
}

// common upload function
func uploadFunc(filePath string, googleUpload bool, cloudinaryUpload bool, megaUpload bool, Directory string) error {
	file_, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file_.Close()
	Directory = strings.TrimPrefix(Directory, "./")

	// fmt.Println("File opened successfully:", file_.Name())

	switch {
	case cloudinaryUpload:
		CloudinaryUrl := goDotEnv("CLOUDINARY_URL")
		err := connection.InitializeCloudinary(CloudinaryUrl)
		if err != nil {
			return fmt.Errorf("failed to initialize cloudinary: %v", err)
		}
		fmt.Println("Directory", Directory)
		fmt.Println("fileee", file_)
		if err = cloudinaryutils.UploadToCloudianary(file_, Directory); err != nil {
			return fmt.Errorf("cloudinary: %v", err)
		}
	case googleUpload:
		err := connection.InitializeGoogleDrive()
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		parentId := utils.PeekStack()

		if err = googleutils.UploadToDrive(file_, parentId); err != nil {
			return fmt.Errorf("drive: %v", err)
		}
	case megaUpload:
		MegaEmail := goDotEnv("MEGA_EMAIL")
		MegaPassword := goDotEnv("MEGA_PASSWORD")
		err := connection.InitializeMega(MegaEmail, MegaPassword)
		if err != nil {
			return fmt.Errorf("mega: %v", err)
		}
		parentId:= utils.PeekStack()
		if err := megautils.UploadToMega(file_, parentId); err != nil {
			return fmt.Errorf("mega: %v", err)
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
	folderName = strings.TrimPrefix(folderName, "./")
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
		err := connection.InitializeGoogleDrive()
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		parentId := utils.PeekStack()
		parentId, err = googleutils.CreateDriveFolder(path.Base(folderName), parentId)
		if err != nil {
			return fmt.Errorf("failed to create folder in drive: %v", err)
		}
		utils.PushToStack(parentId)
	case megaUpload:
		MegaEmail := goDotEnv("MEGA_EMAIL")
		MegaPassword := goDotEnv("MEGA_PASSWORD")
		err := connection.InitializeMega(MegaEmail, MegaPassword)
		if err != nil {
			return fmt.Errorf("failed to initialize mega: %v", err)
		}
		parentId := utils.PeekStack()
		err = megautils.CreateFolderMega(path.Base(folderName), parentId)
		if err != nil {
			return fmt.Errorf("failed to create folder in mega: %v", err)
		}
		utils.PushToStack(folderName)
	default:
		return fmt.Errorf("undefined flag")
	}
	return nil
}

// process dir recursively
func processDirectory(directory string, googleUpload bool, cloudinaryUpload bool, megaUpload bool) error {
	//list items of directory
	entries, err := os.ReadDir(directory)
	// fmt.Println(entries, len(entries))
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", directory, err)
	}
	err = createFolderFunc(directory, googleUpload, cloudinaryUpload, megaUpload)
	if err != nil {
		return fmt.Errorf("error creating folder: %v", err)
	}

	//iterate through dir entries
	for _, entry := range entries {
		fmt.Printf("\nfilePath: %v", entry)
		fullPath := filepath.Join(directory, entry.Name())
		if entry.IsDir() {
			err := processDirectory(fullPath, googleUpload, cloudinaryUpload, megaUpload)
			if err != nil {
				return fmt.Errorf("failed to process subdirectory %s: %v", fullPath, err)
			}
		} else {
			err := uploadFunc(fullPath, googleUpload, cloudinaryUpload, megaUpload, directory)
			if err != nil {
				return fmt.Errorf("failed to upload file %s: %v", fullPath, err)
			}
		}
	}

	utils.PopFromStack()
	return nil
}
