package container

import (
	"github.com/wlchs/blog/internal/repository"
	"go.uber.org/zap"
)

// Container interface defining core application utilities such as logging and DB connectivity
type Container interface {
	GetLogger() *zap.SugaredLogger
	GetRepository() repository.Repository
}

type container struct {
	logger     *zap.SugaredLogger
	repository repository.Repository
}

func CreateContainer(log *zap.SugaredLogger, rep repository.Repository) Container {
	return &container{
		logger:     log,
		repository: rep,
	}
}

// GetLogger returns the logger implementation stored in the container
func (cont container) GetLogger() *zap.SugaredLogger {
	return cont.logger
}

// GetRepository returns the repository implementation stored in the container
func (cont container) GetRepository() repository.Repository {
	return cont.repository
}
