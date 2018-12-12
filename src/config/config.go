package config

import (
	"encoding/json"
	"os"
)

var Settings struct {
	Debug                              bool
	Database                           map[string]string
	ParserParseForwardCommentsAtTime   int
	ParserParseForwardClientTimeout    int
	ProxyProviderClientTimeout         int
	ProxyProviderBaseURL               string
	ParserParseForwardClientUsername   string
	ParserParseForwardClientPassword   string
	ParserParseForwardSleepTime        int
	ParserParseForwardSleepOnErrorTime int
	Pikabu18BotToken                   string
}

func updateSettingsFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Settings)

	return err
}

func init() {
	err := updateSettingsFromFile("config/default_settings.json")
	if err != nil {
		panic(err)
	}
	if _, err := os.Stat("config/settings.json"); err == nil {
		err = updateSettingsFromFile("config/settings.json")
		if err != nil {
			panic(err)
		}
	}
}
