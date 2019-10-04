package parser

import (
	"fmt"
	"sync"

	"github.com/DevAlone/parse_pikabu/helpers"
	"github.com/DevAlone/parse_pikabu/parser/logger"
)

// Main - parser's entry point
func Main() {
	logger.Init()

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
				conf.ParserID += "_copy_" + fmt.Sprint(i)
			}

			parser, err := NewParser(&conf)
			if err != nil {
				helpers.PanicOnError(err)
			}
			wg.Add(1)
			go func() {
				parser.Loop()
				wg.Done()
			}()
		}
	}

	wg.Wait()

	logger.Cleanup()
}
