package globals

import (
	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/models"
)

// ParserResults is a channel for parser results
var ParserResults chan *models.ParserResult

// ParserParseUserTasks is a channel for ParseUserTasks
var ParserParseUserTasks chan *models.ParseUserTask

// ParserParseStoryTasks is a channel for ParseStoryTasks
var ParserParseStoryTasks chan *models.ParseStoryTask

// ParserSimpleTasks is a channel for SimpleTasks
var ParserSimpleTasks chan *models.SimpleTask

// DoNotParseUsers -
var DoNotParseUsers = false

// Init initializes channels
func Init() error {
	ParserResults = make(chan *models.ParserResult, config.Settings.MaxNumberOfTasksInQueue)
	ParserParseUserTasks = make(chan *models.ParseUserTask, config.Settings.MaxNumberOfTasksInQueue)
	ParserParseStoryTasks = make(chan *models.ParseStoryTask, config.Settings.MaxNumberOfTasksInQueue)
	ParserSimpleTasks = make(chan *models.SimpleTask, config.Settings.MaxNumberOfTasksInQueue)

	return nil
}
