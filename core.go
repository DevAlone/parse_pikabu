package main

import (
	"bitbucket.org/d3dev/parse_pikabu/config"
	"bitbucket.org/d3dev/parse_pikabu/logger"
	"github.com/op/go-logging"
	"os"
	"sync"

	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/results_processor"
	"bitbucket.org/d3dev/parse_pikabu/server"
	"bitbucket.org/d3dev/parse_pikabu/statistics"
	"bitbucket.org/d3dev/parse_pikabu/task_manager"
)

func Main() {
	file, err := os.OpenFile("logs/parse_pikabu.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		handleError(err)
	}
	loggingBackend := logging.NewLogBackend(file, "", 0)
	loggingBackendFormatter := logging.NewBackendFormatter(loggingBackend, logger.LogFormat)

	logging.SetBackend(loggingBackend, loggingBackendFormatter)

	if config.Settings.Debug {
		logging.SetLevel(logging.DEBUG, "parse_pikabu/core")
	} else {
		logging.SetLevel(logging.WARNING, "parse_pikabu/core")
	}

	logger.Log.Debug("core started")

	err = models.InitDb()
	if err != nil {
		handleError(err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	// start server
	go func() {
		err := server.Run()
		handleError(err)
		wg.Done()
	}()

	wg.Add(1)
	// start task manager
	go func() {
		err := task_manager.Run()
		handleError(err)
		wg.Done()
	}()

	wg.Add(1)
	// start results processor
	go func() {
		err := results_processor.Run()
		handleError(err)
		wg.Done()
	}()

	wg.Add(1)
	// statistics
	go func() {
		err := statistics.Run()
		handleError(err)
		wg.Done()
	}()

	wg.Wait()
}
