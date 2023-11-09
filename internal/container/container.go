package container

import (
	"github.com/wlchs/blog/internal/logger"
	"github.com/wlchs/blog/internal/repository"
)

// Container interface defining core application utilities such as logging and DB connectivity
type Container interface {
	GetLogger() logger.Logger
	GetRepository() repository.Repository
}

type container struct {
	logger     logger.Logger
	repository repository.Repository
}

func CreateContainer(log logger.Logger, rep repository.Repository) Container {
	return &container{
		logger:     log,
		repository: rep,
	}
}

// GetLogger returns the logger implementation stored in the container
func (cont container) GetLogger() logger.Logger {
	return cont.logger
}

// GetRepository returns the repository implementation stored in the container
func (cont container) GetRepository() repository.Repository {
	return cont.repository
}
