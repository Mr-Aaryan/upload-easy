package connection

import (
	"fmt"

	"github.com/t3rm1n4l/go-mega"
)

var (
	megaClient *mega.Mega
)

func InitializeMega(email, password string) error {
	once.Do(func() {
		client := mega.New()
		err := client.Login(email, password)
		if err != nil {
			initErr = fmt.Errorf("failed to login to MEGA: %w", err)
			return
		}
		megaClient = client
	})
	return initErr
}

func GetMegaClient() (*mega.Mega, error) {
	if megaClient == nil {
		return nil, fmt.Errorf("MEGA client is not initialized. Call InitializeMEGA first")
	}
	return megaClient, nil
}
