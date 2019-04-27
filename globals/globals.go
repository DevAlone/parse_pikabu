package globals

import (
	"sort"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/models"
)

// ParserResults is a channel for parser results
var ParserResults chan *models.ParserResult

// DoNotParseUsers -
var DoNotParseUsers = false

// ParserParseUserTasks is a channel for ParseUserTasks
var ParserParseUserTasks chan *models.ParseUserTask

// ParserParseStoryTasks is a channel for ParseStoryTasks
var ParserParseStoryTasks chan *models.ParseStoryTask

// ParserSimpleTasks is a channel for SimpleTasks
var ParserSimpleTasks chan *models.SimpleTask

type taskProbability struct {
	Probability float64
	ChannelName string
}

// indicates the probabilities for taking specific task
// keep them in range from 0.0 to 1.0, example
// [ 0.1, 0.5, 0.9 ]
var taskProbabilities = []taskProbability{
	taskProbability{Probability: 0.6, ChannelName: "ParserParseUserTasks"},
	taskProbability{Probability: 0.5, ChannelName: "ParserSimpleTasks"},
	taskProbability{Probability: 0.4, ChannelName: "ParserParseStoryTasks"},
}

func initTaskProbabilities() {
	// normalize values
	sum := 0.0
	for _, p := range taskProbabilities {
		sum += p.Probability
	}
	for i := range taskProbabilities {
		taskProbabilities[i].Probability /= sum
	}

	sort.Slice(taskProbabilities, func(i, j int) bool {
		return taskProbabilities[i].Probability < taskProbabilities[j].Probability
	})

	prevProbability := 0.0
	for i := range taskProbabilities {
		taskProbabilities[i].Probability += prevProbability
		if i == len(taskProbabilities)-1 {
			taskProbabilities[i].Probability = 1.0
		}
		prevProbability = taskProbabilities[i].Probability
	}

}

// Init initializes channels
func Init() error {
	ParserResults = make(chan *models.ParserResult, config.Settings.MaxNumberOfTasksInQueue)

	ParserParseUserTasks = make(chan *models.ParseUserTask, config.Settings.MaxNumberOfTasksInQueue)
	ParserParseStoryTasks = make(chan *models.ParseStoryTask, config.Settings.MaxNumberOfTasksInQueue)
	ParserSimpleTasks = make(chan *models.SimpleTask, config.Settings.MaxNumberOfTasksInQueue)

	initTaskProbabilities()

	return nil
}
