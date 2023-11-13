package container

import (
	"github.com/wlchs/blog/internal/jwt"
	"github.com/wlchs/blog/internal/repository"
	"go.uber.org/zap"
)

// Container interface defining core application utilities such as logging and DB connectivity
type Container interface {
	GetLogger() *zap.SugaredLogger

	GetPostRepository() repository.PostRepository
	GetUserRepository() repository.UserRepository

	GetJWTUtils() jwt.TokenUtils
}

// container is the concrete implementation of the Container interface.
type container struct {
	logger *zap.SugaredLogger

	postRepository repository.PostRepository
	userRepository repository.UserRepository

	jwtUtils jwt.TokenUtils
}

// CreateContainer instantiates the application container with all its necessary dependencies.
func CreateContainer(
	log *zap.SugaredLogger,
	postRepository repository.PostRepository,
	userRepository repository.UserRepository,
	jwtUtils jwt.TokenUtils,
) Container {
	return &container{log, postRepository, userRepository, jwtUtils}
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

// GetJWTUtils returns the JWT utility implementation stored in the container.
func (cont container) GetJWTUtils() jwt.TokenUtils {
	return cont.jwtUtils
}
