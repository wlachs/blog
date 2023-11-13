package app

import (
	"github.com/wlchs/blog/internal/container"
	"github.com/wlchs/blog/internal/controller"
	"github.com/wlchs/blog/internal/jwt"
	"github.com/wlchs/blog/internal/logger"
	"github.com/wlchs/blog/internal/repository"
)

// Run initializes the application:
// - Create logger
// - Establish DB connection
// - Define configuration container
// - Bind application routes
func Run() {
	log := logger.CreateLogger()

	rep := repository.CreateRepository()
	postRepository := repository.CreatePostRepository(log, rep)
	userRepository := repository.CreateUserRepository(log, rep)

	jwtUtils := jwt.CreateJWTUtils()

	cont := container.CreateContainer(
		log,
		postRepository,
		userRepository,
		jwtUtils,
	)

	controller.CreateRoutes(cont)
}
