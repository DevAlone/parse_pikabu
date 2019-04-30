package taskmanager

import (
	"time"

	"github.com/go-errors/errors"
)

// TaskManager - task manager which guarantees at least N execution of specific task during cycle
type TaskManager struct {
	tasks              []*Task
	currentTaskIndex   uint
	tasksMap           map[string]*Task
	totalNumberOfTasks uint
}

// NewTaskManager - creates a new task manager.
func NewTaskManager(taskDeclarations map[string]TaskDeclaration) (*TaskManager, error) {
	tm := &TaskManager{}
	tm.currentTaskIndex = 0
	tm.totalNumberOfTasks = 0
	tm.tasks = []*Task{}
	tm.tasksMap = map[string]*Task{}

	for taskName, taskDeclaration := range taskDeclarations {
		task := &Task{
			Name:       taskName,
			Counter:    0,
			Importance: taskDeclaration.Importance,
			Channel:    make(chan interface{}, taskDeclaration.ChannelSize),
		}
		tm.tasks = append(tm.tasks, task)
		tm.tasksMap[taskName] = task
	}

	if len(tm.tasks) == 0 {
		return nil, errors.Errorf("You gotta declare at least one task")
	}

	return tm, nil
}

// PushTask - pushes a new task, will throw an error if the task wasn't previously registred
// Time complexity is O(1)
func (tm *TaskManager) PushTask(name string, data interface{}) error {
	task, found := tm.tasksMap[name]
	if !found {
		return errors.Errorf("task with name '%v' not found", name)
	}
	task.Channel <- data
	// TODO: if it's empty, set the index to task
	tm.totalNumberOfTasks++
	return nil
}

// WaitAndGetTask - waits until there is a task and returns one
// Average and best time complexity is O(1), worst is O(N)
// where N is the number of types of tasks(not tasks themselves)
func (tm *TaskManager) WaitAndGetTask() (string, interface{}) {
	for {
		for tm.totalNumberOfTasks == 0 {
			// TODO: find a better way
			time.Sleep(1 * time.Second)
		}

		task := tm.tasks[tm.currentTaskIndex]
		// if task is there and we're not out of allowed tasks per cycle
		if len(task.Channel) != 0 && task.Counter < task.Importance {
			// TODO: consider using mutex
			tm.totalNumberOfTasks--
			task.Counter++
			return task.Name, <-task.Channel
		}
		task.Counter = 0

		tm.currentTaskIndex++
		if tm.currentTaskIndex >= uint(len(tm.tasks)) {
			tm.currentTaskIndex = 0
		}
	}
}

// TaskDeclaration - declares task
// Importance - number of times that task will be executed in cycle
type TaskDeclaration struct {
	Importance  uint
	ChannelSize uint
}

// Task - internal type to represent a task
type Task struct {
	Name       string
	Counter    uint
	Importance uint
	Channel    chan interface{}
}
