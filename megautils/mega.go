package megautils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"upload-easy/connection"

	"github.com/t3rm1n4l/go-mega"
)

func UploadToMega(file *os.File, parentId string) error {
	// Get mega client
	m, err := connection.GetMegaClient()
	if err != nil {
		return err
	}

	srcPath, err := filepath.Abs(file.Name())
	if err != nil {
		log.Fatalf("failed to fetch file: %v", err)
		return err
	}
	var parentNode *mega.Node
	if parentId == "" {
		parentNode = m.FS.GetRoot()
	} else {
		// Resolve parentId to a *mega.Node
		fmt.Println()
		parentNode, err = FindNodeByName(m.FS, parentId)
		if(err != nil){
			return fmt.Errorf("error occured :'%s'", err)
		}
	}

	// Upload the file
	progress := make(chan int)

	//show file upload progress
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
	_, err = m.UploadFile(srcPath, parentNode, fileData.Name(), &progress)
	if err != nil {
		log.Fatalf("Upload failed: %v", err)
		return err
	}

	log.Printf("File uploaded successfully: %v", fileData.Name())
	return nil
}

func CreateFolderMega(directory string, parentId string) (error) {
	// Get mega client
	m, err := connection.GetMegaClient()
	if err != nil {
		return err
	}
	var parentNode *mega.Node
	if parentId == "" {
		parentNode = m.FS.GetRoot()
	} else {
		// Resolve parentId to a *mega.Node
		parentNode, err = FindNodeByName(m.FS, parentId)
		if parentNode == nil {
			return fmt.Errorf("parent folder with ID '%s' not found", parentId)
		}
		if(err != nil){
			return fmt.Errorf("error occured :'%s'", err)
		}
	}
	_, err = m.CreateDir(directory, parentNode)
	if err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", directory, err)
	}

	// Log or confirm the directory creation success
	fmt.Printf("Folder '%s' created successfully\n", directory)

	return nil
}

// helper func
func FindNodeByName(fs *mega.MegaFS, name string) (*mega.Node, error) {
	// If name is empty or "/", return the root node
	if name == "" || name == "/" {
		return fs.GetRoot(), nil
	}

	// Start from the root node
	currentNode := fs.GetRoot()

	// Split the path into components (e.g., "/folder/subfolder" -> ["folder", "subfolder"])
	components := splitPath(name)

	// Traverse the hierarchy to find the target node
	for _, component := range components {
		found := false

		// Use an SDK method to get children nodes
		children, err := fs.GetChildren(currentNode)
		if err != nil {
			return nil, fmt.Errorf("failed to get children of node '%s': %w", currentNode.GetName(), err)
		}

		// Search for the component in the children
		for _, child := range children {
			if child.GetName() == component {
				currentNode = child
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("node with path '%s' not found", name)
		}
	}

	return currentNode, nil
}

// Helper function to split a path into components
func splitPath(path string) []string {
	trimmed := strings.Trim(path, "/") // Remove leading and trailing slashes
	if trimmed == "" {
		return []string{}
	}
	return strings.Split(trimmed, "/")
}
