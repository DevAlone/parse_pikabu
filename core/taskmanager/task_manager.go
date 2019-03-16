package taskmanager

import (
	"strings"
	"sync"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-pg/pg/orm"

	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"github.com/go-pg/pg"
)

// Run runs goroutines to process tasks
func Run() error {
	var wg sync.WaitGroup

	for _, f := range []func() error{
		storyTasksWorker,
		userTasksWorker,
	} {
		wg.Add(1)
		go func(handler func() error) {
			helpers.PanicOnError(handler())
			wg.Done()
		}(f)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			helpers.PanicOnError(processUserTasks())
			helpers.PanicOnError(processCommunityTasks())

			time.Sleep(time.Duration(config.Settings.WaitBeforeAddingTasksSeconds) * time.Second)
		}
	}()

	wg.Wait()

	return Cleanup()
}

// CompleteTask completes a task
func CompleteTask(tx *pg.Tx, tableName, fieldName string, fieldValue interface{}) error {
	// TODO: refactor to be able to pass the actual model, not a table's name
	var q *orm.Query
	if tx == nil {
		q = models.Db.Model()
	} else {
		q = tx.Model()
	}
	var err error
	switch value := fieldValue.(type) {
	case string:
		_, err = q.Exec(`
			UPDATE `+tableName+` 
			SET is_done = true
			WHERE is_done = false AND LOWER(`+fieldName+`) = ?
		`, strings.ToLower(value))
	default:
		_, err = q.Exec(`
			UPDATE `+tableName+` 
			SET is_done = true
			WHERE is_done = false AND `+fieldName+` = ?
		`, fieldValue)
	}

	return err
}
