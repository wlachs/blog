package app

import (
	"github.com/wlchs/blog/internal/container"
	"github.com/wlchs/blog/internal/logger"
	"github.com/wlchs/blog/internal/repository"
	"github.com/wlchs/blog/internal/services"
	"github.com/wlchs/blog/internal/transport/rest"
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
	cont := container.CreateContainer(log, rep)

	services.InitActions(cont)
	rest.CreateRoutes(cont)
}
