package app

import (
	"fmt"
	"github.com/wlchs/blog/internal/services"
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
		fmt.Println("db initialization error:", dbErr)
		os.Exit(1)
	}

	if err := services.RunInitActions(); err != nil {
		fmt.Println("initialization error:", err)
		os.Exit(1)
	}

	if routeErr := rest.InitRoutes(); routeErr != nil {
		fmt.Println("routing error:", routeErr)
		os.Exit(1)
	}
}
