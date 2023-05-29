package app

import (
	"fmt"
	"os"

	"github.com/wlchs/blog/internal/database"
	"github.com/wlchs/blog/internal/models"
	"github.com/wlchs/blog/internal/transport/rest"
)

// Initialize the application:
// - Establish DB connection
// - Create potentially non-existing DB tables
// - Bind application routes
func Run() {
	if db_err := database.InitDB(); db_err != nil {
		fmt.Println("db initialization error:", db_err)
		os.Exit(1)
	}

	if migration_err := models.InitModels(); migration_err != nil {
		fmt.Println("model migration error:", migration_err)
		os.Exit(1)
	}

	if route_err := rest.InitRoutes(); route_err != nil {
		fmt.Println("routing error:", route_err)
		os.Exit(1)
	}
}
