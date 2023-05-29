package models

import "github.com/wlchs/blog/internal/database"

func InitModels() error {
	if err := database.Agent.AutoMigrate(&Post{}); err != nil {
		return err
	}

	if err := database.Agent.AutoMigrate(&User{}); err != nil {
		return err
	}

	return nil
}
