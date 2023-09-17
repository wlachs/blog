package app

import (
	"github.com/wlchs/blog/internal/services"
	"github.com/wlchs/blog/internal/utils"
	"os"

	"github.com/wlchs/blog/internal/database"
	"github.com/wlchs/blog/internal/transport/rest"
)

// Run initializes the application:
// - Establish DB connection
// - Create potentially non-existing DB tables
// - Bind application routes
func Run() {
	if dbErr := database.InitDB(); dbErr != nil {
		utils.LOG.Errorf("db initialization error: %s", dbErr)
		os.Exit(1)
	}

	if err := services.RunInitActions(); err != nil {
		utils.LOG.Errorf("initialization error: %s", err)
		os.Exit(1)
	}

	if routeErr := rest.InitRoutes(); routeErr != nil {
		utils.LOG.Errorf("routing error: %s", routeErr)
		os.Exit(1)
	}
}
