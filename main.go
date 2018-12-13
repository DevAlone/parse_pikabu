package main

import (
	"bitbucket.org/d3dev/parse_pikabu/config"
	"flag"
	"fmt"
	"github.com/nats-io/go-nats"
	"time"
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

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	err = nc.Publish("foo", []byte("Hello, fucking World!"))
	if err != nil {
		panic(err)
	}

	_, err = nc.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Recieved a message: %s\n", string(m.Data))
	})

	if err != nil {
		panic(err)
	}


	// msg, err
	for true {
		time.Sleep(1 * time.Second)
		err = nc.Publish("foo", []byte("Hello, fucking World!1"))
		if err != nil {
			panic(err)
		}
	}
}

