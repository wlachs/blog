package services

import (
	"github.com/wlchs/blog/internal/container"
	"github.com/wlchs/blog/internal/models"
	"os"
)

// InitActions makes sure that the application is initialized for first time use.
func InitActions(cont container.Container) {
	log := cont.GetLogger()

	// Initialize DB models
	if err := models.InitModels(); err != nil {
		log.Error("DB model initialization failed: %s", err)
		os.Exit(1)
	}

	// Create user if it doesn't exist yet
	if err := RegisterFirstUser(); err != nil {
		log.Error("first user registration failed: %s", err)
	}

	log.Info("init actions done")
}
