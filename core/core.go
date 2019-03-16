package core

import (
	"sync"

	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/core/resultsprocessor"
	"bitbucket.org/d3dev/parse_pikabu/core/server"
	"bitbucket.org/d3dev/parse_pikabu/core/statistics"
	"bitbucket.org/d3dev/parse_pikabu/core/taskmanager"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/models"
)

func Main() {
	logger.Init()

	err := models.InitDb()
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
		err := taskmanager.Run()
		helpers.PanicOnError(err)
		wg.Done()
	}()

	wg.Add(1)
	// start results processor
	go func() {
		err := resultsprocessor.Run()
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

	logger.Cleanup()
}
