package taskmanager

import (
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-errors/errors"
)

// GetParserTask - fetches a task
func GetParserTask() (interface{}, error) {
	taskID, taskData := CoreTaskManager.WaitAndGetTask()
	switch taskID {
	case UpdateUserTask:
		// will panic on wrong type
		return taskData.(*models.ParseUserTask), nil
	case ParseNewUserTask:
		return taskData.(*models.ParseUserTask), nil
	case ParseDeletedOrNeverExistedUserTask:
		return taskData.(*models.ParseUserTask), nil
	case UpdateStoryTask:
		return taskData.(*models.ParseStoryTask), nil
	case ParseNewStoryTask:
		return taskData.(*models.ParseStoryTask), nil
	case ParseDeletedOrNeverExistedStoryTask:
		return taskData.(*models.ParseStoryTask), nil
	case ParseAllCommunitiesTask:
		return taskData.(*models.SimpleTask), nil

	default:
		return nil, errors.Errorf("forgot to handle task with id %v", taskID)
	}
}
