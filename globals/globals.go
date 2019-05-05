package globals

import (
	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/models"
)

// ParserResults is a channel for parser results
var ParserResults chan *models.ParserResult

// DoNotParseUsers -
var DoNotParseUsers = false

// DoNotParseStories -
var DoNotParseStories = false

// Init initializes channels
func Init() error {
	ParserResults = make(chan *models.ParserResult, config.Settings.MaxNumberOfTasksInQueue)

	return nil
}
