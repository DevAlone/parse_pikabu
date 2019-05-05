package taskmanager

import (
	"sync"
	"time"

	"github.com/go-errors/errors"
)

// TaskManager - task manager which guarantees at least N execution of specific task during cycle
type TaskManager struct {
	tasks              []*Task
	currentTaskIndex   uint
	tasksMap           map[int]*Task
	totalNumberOfTasks uint
	mutex              *sync.Mutex
}

// NewTaskManager - creates a new task manager.
func NewTaskManager(taskDeclarations map[int]TaskDeclaration) (*TaskManager, error) {
	tm := &TaskManager{}
	tm.currentTaskIndex = 0
	tm.totalNumberOfTasks = 0
	tm.tasks = []*Task{}
	tm.tasksMap = map[int]*Task{}
	tm.mutex = &sync.Mutex{}

	for taskID, taskDeclaration := range taskDeclarations {
		task := &Task{
			ID:         taskID,
			Counter:    0,
			Importance: taskDeclaration.Importance,
			Channel:    make(chan interface{}, taskDeclaration.ChannelSize),
		}
		tm.tasks = append(tm.tasks, task)
		tm.tasksMap[taskID] = task
	}

	if len(tm.tasks) == 0 {
		return nil, errors.Errorf("You gotta declare at least one task")
	}

	return tm, nil
}

// PushTask - pushes a new task, will throw an error if the task wasn't previously registred
// Time complexity is O(1)
func (tm *TaskManager) PushTask(id int, data interface{}) error {
	tm.mutex.Lock()
	task, found := tm.tasksMap[id]
	tm.mutex.Unlock()

	if !found {
		return errors.Errorf("task with id '%v' not found", id)
	}
	task.Channel <- data
	// TODO: if it's empty, set the index to task

	tm.mutex.Lock()
	tm.totalNumberOfTasks++
	tm.mutex.Unlock()

	return nil
}

// WaitAndGetTask - waits until there is a task and returns one
// Average and best time complexity is O(1), worst is O(N)
// where N is the number of types of tasks(not tasks themselves)
func (tm *TaskManager) WaitAndGetTask() (int, interface{}) {
	for {
		for tm.totalNumberOfTasks == 0 {
			// TODO: find a better way
			time.Sleep(1 * time.Second)
		}

		tm.mutex.Lock()
		task := tm.tasks[tm.currentTaskIndex]
		tm.mutex.Unlock()
		// if task is there and we're not out of allowed tasks per cycle
		if len(task.Channel) != 0 && task.Counter < task.Importance {
			// TODO: consider using mutex
			tm.mutex.Lock()
			tm.totalNumberOfTasks--
			tm.mutex.Unlock()
			task.Counter++
			return task.ID, <-task.Channel
		}
		task.Counter = 0

		tm.mutex.Lock()
		tm.currentTaskIndex++
		if tm.currentTaskIndex >= uint(len(tm.tasks)) {
			tm.currentTaskIndex = 0
		}
		tm.mutex.Unlock()
	}
}

// MustTaskChannelSize -
func (tm *TaskManager) MustTaskChannelSize(taskID int) int64 {
	tm.mutex.Lock()
	task, found := tm.tasksMap[taskID]
	tm.mutex.Unlock()

	if !found {
		panic(errors.Errorf("task with id %v not found", taskID))
	}

	return int64(len(task.Channel))
}

// TaskDeclaration - declares task
// Importance - number of times that task will be executed in cycle
type TaskDeclaration struct {
	Importance  uint
	ChannelSize uint
}

// Task - internal type to represent a task
type Task struct {
	ID         int
	Counter    uint
	Importance uint
	Channel    chan interface{}
}
