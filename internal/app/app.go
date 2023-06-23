package app

import (
	"fmt"
	"os"

	"github.com/wlchs/blog/internal/database"
	"github.com/wlchs/blog/internal/models"
	"github.com/wlchs/blog/internal/transport/rest"
)

// Run initializes the application:
// - Establish DB connection
// - Create potentially non-existing DB tables
// - Bind application routes
func Run() {
	if dbErr := database.InitDB(); dbErr != nil {
		fmt.Println("db initialization error:", dbErr)
		os.Exit(1)
	}

	if migrationErr := models.InitModels(); migrationErr != nil {
		fmt.Println("model migration error:", migrationErr)
		os.Exit(1)
	}

	if routeErr := rest.InitRoutes(); routeErr != nil {
		fmt.Println("routing error:", routeErr)
		os.Exit(1)
	}
}
