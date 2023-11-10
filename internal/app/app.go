package app

import (
	"github.com/wlchs/blog/internal/container"
	"github.com/wlchs/blog/internal/controller"
	"github.com/wlchs/blog/internal/logger"
	"github.com/wlchs/blog/internal/repository"
	"github.com/wlchs/blog/internal/services"
)

// Run initializes the application:
// - Create logger
// - Establish DB connection
// - Define configuration container
// - Create potentially non-existing DB tables
// - Bind application routes
func Run() {
	log := logger.CreateLogger()
	rep := repository.CreateRepository()
	postRepository := repository.CreatePostRepository(log, rep)
	userRepository := repository.CreateUserRepository(log, rep)
	cont := container.CreateContainer(log, postRepository, userRepository)

	services.InitActions(cont)
	controller.CreateRoutes(cont)
}
