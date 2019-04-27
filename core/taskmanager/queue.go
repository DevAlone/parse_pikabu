package taskmanager

import (
	"reflect"

	"bitbucket.org/d3dev/parse_pikabu/globals"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-errors/errors"
)

// PushTaskToQueue -
// TODO: check if it's used
func PushTaskToQueue(taskPtr interface{}) error {
	switch t := taskPtr.(type) {
	case *models.ParseUserTask:
		globals.ParserParseUserTasks <- t
		return nil
	case *models.SimpleTask:
		globals.ParserSimpleTasks <- t
		return nil
	case *models.ParseStoryTask:
		globals.ParserParseStoryTasks <- t
		return nil
	default:
		return errors.Errorf("trying to push undeclared type of task %v %v", reflect.TypeOf(t), t)
	}
}
