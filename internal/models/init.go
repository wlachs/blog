package models

import (
	"github.com/wlchs/blog/internal/database"
	"github.com/wlchs/blog/internal/utils"
)

func InitModels() error {
	if err := database.Agent.AutoMigrate(&Post{}); err != nil {
		return err
	}
	if err := database.Agent.AutoMigrate(&User{}); err != nil {
		return err
	}
	utils.LOG.Debugln("db models initialized")
	return nil
}
