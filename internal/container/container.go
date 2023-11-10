package container

import (
	"github.com/wlchs/blog/internal/repository"
	"go.uber.org/zap"
)

// Container interface defining core application utilities such as logging and DB connectivity
type Container interface {
	GetLogger() *zap.SugaredLogger
	GetPostRepository() repository.PostRepository
	GetUserRepository() repository.UserRepository
}

type container struct {
	logger         *zap.SugaredLogger
	postRepository repository.PostRepository
	userRepository repository.UserRepository
}

func CreateContainer(log *zap.SugaredLogger, postRepository repository.PostRepository, userRepository repository.UserRepository) Container {
	return &container{
		logger:         log,
		postRepository: postRepository,
		userRepository: userRepository,
	}
}

// GetLogger returns the logger implementation stored in the container
func (cont container) GetLogger() *zap.SugaredLogger {
	return cont.logger
}

// GetPostRepository returns the post repository implementation stored in the container
func (cont container) GetPostRepository() repository.PostRepository {
	return cont.postRepository
}

// GetUserRepository returns the user repository implementation stored in the container
func (cont container) GetUserRepository() repository.UserRepository {
	return cont.userRepository
}
