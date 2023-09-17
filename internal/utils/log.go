package utils

import (
	"go.uber.org/zap"
	"os"
)

var (
	LOG *zap.SugaredLogger
)

// InitLogger initializes the logging framework
func InitLogger() error {
	if logger, err := getLoggerByMode(); err == nil {
		LOG = logger.Sugar()
		LOG.Debugln("logging initialized")
		return nil
	} else {
		return err
	}
}

// getLoggerByMode gets the logger configured according to the application mode. In development the development config is used,
// in release the production version.
func getLoggerByMode() (*zap.Logger, error) {
	mode := os.Getenv("GIN_MODE")
	if mode == "release" {
		return zap.NewProduction()
	} else {
		return zap.NewDevelopment()
	}
}
