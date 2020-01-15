package core

import (
	"sync"

	"github.com/DevAlone/parse_pikabu/core/logger"
	"github.com/DevAlone/parse_pikabu/core/resultsprocessor"
	"github.com/DevAlone/parse_pikabu/core/server"
	"github.com/DevAlone/parse_pikabu/core/statistics"
	"github.com/DevAlone/parse_pikabu/core/taskmanager"
	"github.com/DevAlone/parse_pikabu/helpers"
	"github.com/DevAlone/parse_pikabu/modelhooks"
	"github.com/DevAlone/parse_pikabu/models"
	"github.com/DevAlone/parse_pikabu/telegram"
)

// Main - entry point for core of the project
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

	wg.Add(1)
	// start telegram notifier
	go func() {
		err := telegram.RunTelegramNotifier()
		helpers.PanicOnError(err)
		wg.Done()
	}()

	wg.Add(1)
	// start modelhooks handler
	go func() {
		err := modelhooks.RunModelHooksHandler()
		helpers.PanicOnError(err)
		wg.Done()
	}()

	wg.Wait()

	logger.Cleanup()
}
