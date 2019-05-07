package config

import (
	"encoding/json"
	"os"
)

// Settings is struct for global settings of app
var Settings struct {
	Debug                      bool
	LogSQLQueries              bool
	Database                   map[string]string
	ProxyProviderClientTimeout int
	ProxyProviderBaseURL       string
	Pikabu18BotToken           string
	ServerListeningAddress     string
	AMQPAddress                string
	// in seconds
	MaximumParseUserTaskProcessingTime  int
	MaximumParseStoryTaskProcessingTime int
	CommunitiesProcessingPeriod         int
	ServerMaximumNumberOfResultsPerPage uint
	// time in seconds to consider user as new
	NewUserTime                        int
	NewUsersUpdatingPeriod             int
	UsersUpdatingPeriodIncreasingValue int
	UsersMinUpdatingPeriod             int
	UsersMaxUpdatingPeriod             int
	MaxNumberOfTasksInQueue            int
	// actual number if number of thread multiplied by this value
	NumberOfTasksProcessorsMultiplier    int
	NumberOfNewUsersGap                  int
	NumberOfNewStoriesGap                int
	WaitBeforeAddingNewUserTasksSeconds  int
	WaitBeforeAddingNewStoryTasksSeconds int
	AddNewUsersEachNMinutes              int

	UpdateUserTaskImportance                      uint
	ParseNewUserTaskImportance                    uint
	ParseDeletedOrNeverExistedUserTaskImportance  uint
	UpdateStoryTaskImportance                     uint
	ParseNewStoryTaskImportance                   uint
	ParseDeletedOrNeverExistedStoryTaskImportance uint
	ParseAllCommunitiesTaskImportance             uint
}

// UpdateSettingsFromFile fills settings from the file
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
	Settings.LogSQLQueries = false
	Settings.Database = map[string]string{
		"Name":     "test",
		"Username": "test",
		"Password": "test",
	}
	Settings.ProxyProviderClientTimeout = 60
	Settings.ProxyProviderBaseURL = ""
	Settings.Pikabu18BotToken = ""
	Settings.ServerListeningAddress = "0.0.0.0:8080"
	Settings.AMQPAddress = "amqp://guest:guest@localhost:5672/"
	Settings.MaximumParseUserTaskProcessingTime = 2 * 60
	Settings.MaximumParseStoryTaskProcessingTime = 2 * 60
	Settings.CommunitiesProcessingPeriod = 3600
	Settings.ServerMaximumNumberOfResultsPerPage = 1024
	Settings.NewUserTime = 3600 * 24 * 7
	Settings.NewUsersUpdatingPeriod = 3600 * 24
	Settings.UsersUpdatingPeriodIncreasingValue = 6 * 3600
	Settings.UsersMinUpdatingPeriod = 3600 * 12
	Settings.UsersMaxUpdatingPeriod = 3600 * 24 * 7
	Settings.MaxNumberOfTasksInQueue = 128
	Settings.NumberOfTasksProcessorsMultiplier = 1
	Settings.NumberOfNewUsersGap = 128
	Settings.NumberOfNewStoriesGap = 1024
	Settings.WaitBeforeAddingNewUserTasksSeconds = 10
	Settings.WaitBeforeAddingNewStoryTasksSeconds = 60
	Settings.AddNewUsersEachNMinutes = 10

	Settings.UpdateUserTaskImportance = 100
	Settings.ParseNewUserTaskImportance = 10
	Settings.ParseDeletedOrNeverExistedUserTaskImportance = 1
	Settings.ParseNewStoryTaskImportance = 10
	Settings.UpdateStoryTaskImportance = 1
	Settings.ParseDeletedOrNeverExistedStoryTaskImportance = 5
	Settings.ParseAllCommunitiesTaskImportance = 1
}
