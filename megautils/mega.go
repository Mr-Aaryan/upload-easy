package megautils

import (
	"log"
	"os"
	"path/filepath"

	"github.com/t3rm1n4l/go-mega"
)

func UploadToMega(file *os.File, email string, password string) error {
	// Create a new Mega client
	m := mega.New()

	// Login to your MEGA account
	err := m.Login(email, password)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
		return err
	}

	log.Println("login successful")

	srcPath, err := filepath.Abs(file.Name())
	if err != nil {
		log.Fatalf("failed to fetch file: %v", err)
		return err
	}
	// parentNode := m.FS.GetRoot()
	parentNode := m.FS.GetRoot()

	// Upload the file
	progress := make(chan int)
	go func() {
		for p := range progress {
			log.Printf("Upload progress: %d%%\n", p)
		}
	}()
	fileData, err := file.Stat()
	if err != nil {
		log.Fatalf("failed to read file stats %s", err)
	}
	log.Println(fileData.Name())
	node, err := m.UploadFile(srcPath, parentNode, fileData.Name(), &progress)
	if err != nil {
		log.Fatalf("Upload failed: %v", err)
		return err
	}

	log.Printf("File uploaded successfully: %v", node)
	return nil
}
