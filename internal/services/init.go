package services

import (
	"github.com/wlchs/blog/internal/models"
	"github.com/wlchs/blog/internal/utils"
)

// RunInitActions makes sure that the application is initialized for first time use.
func RunInitActions() error {
	// Initialize logging
	if err := utils.InitLogger(); err != nil {
		return err
	}
	// Initialize DB models
	if err := models.InitModels(); err != nil {
		return err
	}
	// Create user if it doesn't exist yet
	if err := RegisterFirstUser(); err != nil {
		return err
	}
	return nil
}
