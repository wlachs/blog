package logger

import (
	"fmt"
	"go.uber.org/zap"
	"os"
)

// Logger is a custom interface encapsulating the underlying logging framework.
// Goal is to achieve library-agnostic logging.
type Logger interface {
	Debug(msg string, data ...interface{})
	Info(msg string, data ...interface{})
	Warn(msg string, data ...interface{})
	Error(msg string, data ...interface{})
}

// logger struct to invoke methods on
type logger struct {
	Zap *zap.SugaredLogger
}

// CreateLogger initializes the logging framework
func CreateLogger() Logger {
	return &logger{Zap: getLogger()}
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

// Debug logs a message with log level DEBUG
func (log *logger) Debug(msg string, data ...interface{}) {
	log.Zap.Debugf(msg, data)
}

// Info logs a message with log level INFO
func (log *logger) Info(msg string, data ...interface{}) {
	log.Zap.Infof(msg, data)
}

// Warn logs a message with log level WARN
func (log *logger) Warn(msg string, data ...interface{}) {
	log.Zap.Warnf(msg, data)
}

// Error logs a message with log level ERROR
func (log *logger) Error(msg string, data ...interface{}) {
	log.Zap.Errorf(msg, data)
}
