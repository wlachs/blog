package app

import (
	"github.com/wlchs/blog/internal/container"
	"github.com/wlchs/blog/internal/controller"
	"github.com/wlchs/blog/internal/db"
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
	database := db.ConnectToMySQL()
	rep := repository.CreateRepository(database)
	postRepository := repository.CreatePostRepository(log, rep)
	userRepository := repository.CreateUserRepository(log, rep)
	jwtUtils := jwt.CreateTokenUtils(log)

	cont := container.CreateContainer(
		log,
		postRepository,
		userRepository,
		jwtUtils,
	)

	controller.CreateRoutes(cont)
}
