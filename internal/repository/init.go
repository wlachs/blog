package repository

import (
	"github.com/wlchs/blog/internal/container"
)

func InitModels(cont container.Container) error {
	log := cont.GetLogger()
	rep := cont.GetRepository()
	if err := rep.AutoMigrate(&Post{}); err != nil {
		return err
	}
	if err := rep.AutoMigrate(&User{}); err != nil {
		return err
	}
	log.Debug("db models initialized")
	return nil
}
