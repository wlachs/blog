package logger

import (
	"fmt"
	"go.uber.org/zap"
	"os"
)

// CreateLogger initializes the logging framework
func CreateLogger() *zap.SugaredLogger {
	return getLogger()
}

// getLogger tries to fetch the logger.
// If, for some reason, the method fails, the application exits
func getLogger() *zap.SugaredLogger {
	l, err := getLoggerByMode()
	if err != nil {
		fmt.Printf("logging initialization failed: %s", err)
		os.Exit(1)
	}
	return l.Sugar()
}

// getLoggerByMode gets the logger configured according to the application mode.
// In development the development config is used, in release the production version.
func getLoggerByMode() (*zap.Logger, error) {
	mode := os.Getenv("GIN_MODE")
	if mode == "release" {
		return zap.NewProduction()
	} else {
		return zap.NewDevelopment()
	}
}
