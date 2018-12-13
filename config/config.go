package config

import (
	"encoding/json"
	"os"
)

var Settings struct {
	Debug                      bool
	Database                   map[string]string
	ProxyProviderClientTimeout int
	ProxyProviderBaseURL       string
	Pikabu18BotToken           string
}

func UpdateSettingsFromFile(filename string) error {
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
	Settings.Debug = true
	Settings.Database = map[string]string{
		"Name":     "test",
		"Username": "test",
		"Password": "test",
	}
	Settings.ProxyProviderClientTimeout = 60
	Settings.ProxyProviderBaseURL = "https://eivailohciihi4uquapach7abei9iesh.d3d.info/api/v1/"
	Settings.Pikabu18BotToken = ""
}
