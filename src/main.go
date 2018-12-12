package main

import (
	"flag"
	. "models"
	"parsers"
	"sync"
	"telegram/pikabu_18_bot"
)

func main() {
	withoutParsers := flag.Bool("without-parsers", false, "run without parsers")
	withoutBots := flag.Bool("without-bots", false, "run without bots")
	flag.Parse()

	err := InitDb()
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
	wg.Wait()
}
