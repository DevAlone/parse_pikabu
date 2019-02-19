package core

import (
	"os"
	"sync"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	logging "github.com/op/go-logging"

	"bitbucket.org/d3dev/parse_pikabu/core/results_processor"
	"bitbucket.org/d3dev/parse_pikabu/core/server"
	"bitbucket.org/d3dev/parse_pikabu/core/statistics"
	"bitbucket.org/d3dev/parse_pikabu/core/task_manager"
	"bitbucket.org/d3dev/parse_pikabu/models"
)

func Main() {
	file, err := os.OpenFile("logs/parse_pikabu.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
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

	err = models.InitDb()
	if err != nil {
		helpers.PanicOnError(err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	// start server
	go func() {
		err := server.Run()
		helpers.PanicOnError(err)
		wg.Done()
	}()

	wg.Add(1)
	// start task manager
	go func() {
		err := task_manager.Run()
		helpers.PanicOnError(err)
		wg.Done()
	}()

	wg.Add(1)
	// start results processor
	go func() {
		err := results_processor.Run()
		helpers.PanicOnError(err)
		wg.Done()
	}()

	wg.Add(1)
	// statistics
	go func() {
		err := statistics.Run()
		helpers.PanicOnError(err)
		wg.Done()
	}()

	wg.Wait()
}
