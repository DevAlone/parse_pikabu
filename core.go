package main

import (
	"bitbucket.org/d3dev/parse_pikabu/config"
	"bitbucket.org/d3dev/parse_pikabu/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/results_processor"
	"bitbucket.org/d3dev/parse_pikabu/server"
	"bitbucket.org/d3dev/parse_pikabu/task_manager"
	"flag"
	"github.com/op/go-logging"
	"os"
	"sync"
)

func Main() {
	file, err := os.OpenFile("logs/parse_pikabu.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	// loggingBackend := logger.NewLogBackend(os.Stderr, "", 0)
	loggingBackend := logging.NewLogBackend(file, "", 0)
	loggingBackendFormatter := logging.NewBackendFormatter(loggingBackend, logger.LogFormat)

	logging.SetBackend(loggingBackend, loggingBackendFormatter)
	logger.Log.Debug("app started")

	configFilePath := flag.String("config", "config.json", "config file")

	if configFilePath == nil {
		panic("configFilePath is nil")
	}

	err = config.UpdateSettingsFromFile(*configFilePath)
	if err != nil {
		panic(err)
	}

	err = models.InitDb()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(3)

	// start server
	go func() {
		err := server.Run()
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	// start task manager
	go func() {
		err := task_manager.Run()
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	// start results processor
	go func() {
		err := results_processor.Run()
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	wg.Wait()
}
