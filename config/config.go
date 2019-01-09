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
	ServerListeningAddress     string
	AMQPAddress                string
	// in seconds
	MaximumTaskProcessingTime           int
	CommunitiesProcessingPeriod         int
	ServerMaximumNumberOfResultsPerPage uint
	// time in seconds to consider user as new
	NewUserTime                        int
	NewUsersUpdatingPeriod             int
	UsersUpdatingPeriodIncreasingValue int
	UsersMinUpdatingPeriod             int
	UsersMaxUpdatingPeriod             int
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
	Settings.Debug = false
	Settings.Database = map[string]string{
		"Name":     "test",
		"Username": "test",
		"Password": "test",
	}
	Settings.ProxyProviderClientTimeout = 60
	Settings.ProxyProviderBaseURL = "https://eivailohciihi4uquapach7abei9iesh.d3d.info/api/v1/"
	Settings.Pikabu18BotToken = ""
	Settings.ServerListeningAddress = "0.0.0.0:8080"
	Settings.AMQPAddress = "amqp://guest:guest@localhost:5672/"
	Settings.MaximumTaskProcessingTime = 60
	Settings.CommunitiesProcessingPeriod = 3600
	Settings.ServerMaximumNumberOfResultsPerPage = 1024
	Settings.NewUserTime = 3600 * 24 * 7
	Settings.NewUsersUpdatingPeriod = 3600 * 24
	Settings.UsersUpdatingPeriodIncreasingValue = 3600
	Settings.UsersMinUpdatingPeriod = 3600 * 12
	Settings.UsersMaxUpdatingPeriod = 3600 * 24 * 7
}
