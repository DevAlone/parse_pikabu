package globals

import (
	"github.com/DevAlone/parse_pikabu/core/config"
	"github.com/DevAlone/parse_pikabu/models"
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
