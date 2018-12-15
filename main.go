package main

import (
	"bitbucket.org/d3dev/parse_pikabu/config"
	"bitbucket.org/d3dev/parse_pikabu/models"
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
	wg.Add(2)

	// start server
	go func() {
		err := server.Run()
		wg.Done()
		if err != nil {
			panic(err)
		}
	}()

	// start task manager
	go func() {
		err := task_manager.Run()
		wg.Done()
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()

	// withoutParsers := flag.Bool("without-parsers", false, "run without parsers")
	// withoutBots := flag.Bool("without-bots", false, "run without bots")
	/*flag.Parse()

	err := models.InitDb()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	if !*withoutBots {
		wg.Add(1)
		go pikabu_18_bot.RunPikabu18Bot()
	}
	if !*withoutParsers {
		wg.Add(1)
		go parsers.Run()
	}
	wg.Wait()*/
}
