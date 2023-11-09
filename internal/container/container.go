package container

import (
	"github.com/wlchs/blog/internal/logger"
	"github.com/wlchs/blog/internal/repository"
)

// Container interface defining core application utilities such as logging and DB connectivity
type Container interface {
	GetRepository() repository.Repository
	GetLogger() logger.Logger
}
