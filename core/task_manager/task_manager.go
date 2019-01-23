package task_manager

import (
	"strings"
	"time"

	"github.com/go-pg/pg"
)

func Run() error {
	for true {
		if err := processUserTasks(); err != nil {
			return err
		}
		if err := processCommunityTasks(); err != nil {
			return err
		}

		time.Sleep(1 * time.Second)
	}

	return Cleanup()
}

func CompleteTask(tx *pg.Tx, tableName, fieldName string, fieldValue interface{}) error {
	// TODO: refactor to be able to pass the actual model, not a table's name
	var err error
	switch value := fieldValue.(type) {
	case string:
		_, err = tx.Model().Exec(`
			UPDATE `+tableName+` 
			SET is_done = true
			WHERE is_done = false AND LOWER(`+fieldName+`) = ?
		`, strings.ToLower(value))
	default:
		_, err = tx.Model().Exec(`
			UPDATE `+tableName+` 
			SET is_done = true
			WHERE is_done = false AND `+fieldName+` = ?
		`, fieldValue)
	}

	return err
}
