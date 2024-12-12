package connection

import (
	"fmt"
	"sync"

	"github.com/cloudinary/cloudinary-go/v2"
)

var (
	cld     *cloudinary.Cloudinary
	once    sync.Once
	initErr error
)

// initialize cloudinary
func InitializeCloudinary(CloudinaryUrl string) error {
	once.Do(func() {
		cld, initErr = cloudinary.NewFromURL(CloudinaryUrl)
	})
	return initErr
}

// returns the initialized cloudinary instance
func GetCloudinary() (*cloudinary.Cloudinary, error) {
	if cld == nil {
		return nil, fmt.Errorf("cloudinary is not initialized. Call Initialize first")
	}
	return cld, nil
}
