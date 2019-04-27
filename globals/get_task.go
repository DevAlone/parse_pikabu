package globals

import (
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

// GetTaskPtr - fetches a task with a probability specified in taskProbabilities config or the one with the greatest probability if no such
func GetTaskPtr() (interface{}, error) {
	randNumber := rand.Float64()
	for _, taskProbability := range taskProbabilities {
		if taskProbability.Probability <= randNumber {
			// found our channel
			task, found, err := tryToGetTaskFromChannelByName(taskProbability.ChannelName)
			if err != nil {
				return nil, err
			}
			if found {
				return task, nil
			}
			break
		}
	}
	for {
		for i := range taskProbabilities {
			reversedIndex := len(taskProbabilities) - i - 1
			taskProbability := &taskProbabilities[reversedIndex]
			task, found, err := tryToGetTaskFromChannelByName(taskProbability.ChannelName)
			if err != nil {
				return nil, err
			}
			if found {
				return task, nil
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func tryToGetTaskFromChannelByName(name string) (interface{}, bool, error) {
	switch name {
	case "ParserParseUserTasks":
		if len(ParserParseUserTasks) <= 0 {
			return nil, false, nil
		}
		return <-ParserParseUserTasks, true, nil
	case "ParserSimpleTasks":
		if len(ParserSimpleTasks) <= 0 {
			return nil, false, nil
		}
		return <-ParserSimpleTasks, true, nil
	case "ParserParseStoryTasks":
		if len(ParserParseStoryTasks) <= 0 {
			return nil, false, nil
		}
		return <-ParserParseStoryTasks, true, nil
	}

	return nil, false, errors.Errorf("tryToGetTaskFromChannelByName(): wrong channel name %v", name)
}
