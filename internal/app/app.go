package app

import (
	"github.com/wlachs/blog/internal/container"
	"github.com/wlachs/blog/internal/controller"
	"github.com/wlachs/blog/internal/db"
	"github.com/wlachs/blog/internal/jwt"
	"github.com/wlachs/blog/internal/logger"
	"github.com/wlachs/blog/internal/repository"
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
