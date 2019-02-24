package parser

import (
	"fmt"
	"os"
	"sync"

	"bitbucket.org/d3dev/parse_pikabu/helpers"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/parser/logger"
	logging "github.com/op/go-logging"
)

func Main() {
	file, err := os.OpenFile("logs/parser.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	helpers.PanicOnError(err)
	loggingBackend := logging.NewLogBackend(file, "", 0)
	loggingBackendFormatter := logging.NewBackendFormatter(loggingBackend, logger.LogFormat)

	logging.SetBackend(loggingBackend, loggingBackendFormatter)

	if config.Settings.Debug {
		logging.SetLevel(logging.DEBUG, "parse_pikabu")
	} else {
		logging.SetLevel(logging.WARNING, "parse_pikabu")
	}

	logger.Log.Debug("parsers started")

	parsersConfig, err := NewParsersConfigFromFile("parsers.config.json")
	helpers.PanicOnError(err)

	var wg sync.WaitGroup

	for _, parserConfig := range parsersConfig.Configs {
		// var configs
		for i := uint(0); i < parserConfig.NumberOfInstances; i++ {
			var conf ParserConfig
			conf = parserConfig
			if i != 0 {
				conf.ParserId += "_copy_" + fmt.Sprint(i)
			}

			parser, err := NewParser(&conf)
			if err != nil {
				helpers.PanicOnError(err)
			}
			wg.Add(1)
			go func() {
				parser.Loop()
				wg.Done()
				err := parser.Cleanup()
				helpers.PanicOnError(err)
			}()
		}
	}

	wg.Wait()
}
