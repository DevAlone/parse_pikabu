package globals

import (
	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/models"
)

var ParserResults chan *models.ParserResult
var ParserParseUserTasks chan *models.ParseUserTask
var ParserSimpleTasks chan *models.SimpleTask

// test feature, uses go channels instead of amqp
var SingleProcessMode = false

func Init() error {
	ParserResults = make(chan *models.ParserResult, config.Settings.MaxNumberOfTasksInQueue)
	ParserParseUserTasks = make(chan *models.ParseUserTask, config.Settings.MaxNumberOfTasksInQueue)
	ParserSimpleTasks = make(chan *models.SimpleTask, config.Settings.MaxNumberOfTasksInQueue)

	return nil
}
