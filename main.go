package parse_pikabu

import (
	"bitbucket.org/d3dev/parse_pikabu/config"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/results_processor"
	"bitbucket.org/d3dev/parse_pikabu/server"
	"bitbucket.org/d3dev/parse_pikabu/task_manager"
	"flag"
	"sync"
)

func main() {
	configFilePath := flag.String("config", "config.json", "config file")

	if configFilePath == nil {
		panic("configFilePath is nil")
	}

	err := config.UpdateSettingsFromFile(*configFilePath)
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
