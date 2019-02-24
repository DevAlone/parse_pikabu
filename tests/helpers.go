package tests

import (
	"os"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	logging "github.com/op/go-logging"
)

func initLogs() {
	file, err := os.OpenFile("../logs/tests.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		helpers.PanicOnError(err)
	}
	loggingBackend := logging.NewLogBackend(file, "", 0)
	loggingBackendFormatter := logging.NewBackendFormatter(loggingBackend, logger.LogFormat)

	logging.SetBackend(loggingBackend, loggingBackendFormatter)

	if config.Settings.Debug {
		logging.SetLevel(logging.DEBUG, "parse_pikabu")
	} else {
		logging.SetLevel(logging.WARNING, "parse_pikabu")
	}

}
