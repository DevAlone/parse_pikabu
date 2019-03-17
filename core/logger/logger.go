package logger

import (
	"os"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"github.com/go-errors/errors"
	"github.com/sirupsen/logrus"
)

// Log is a log fore core part of the project
var Log *logrus.Logger
var logFile *os.File

// Init initializes logs
func Init() {
	var err error
	logFile, err = os.OpenFile("logs/core.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}

	Log = logrus.New()
	Log.SetOutput(logFile)
	// Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	if config.Settings.Debug {
		Log.SetLevel(logrus.DebugLevel)
	} else {
		Log.SetLevel(logrus.WarnLevel)
	}
}

// Cleanup cleans up after logs
func Cleanup() {
	_ = logFile.Close()
}

// LogError logs go-errors error
func LogError(err error) {
	if err == nil {
		return
	}
	if e, ok := err.(*errors.Error); ok {
		Log.Error(e.ErrorStack())
	} else {
		Log.Error(err)
	}
}
