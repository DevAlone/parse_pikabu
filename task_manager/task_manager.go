package task_manager

import (
	"github.com/go-pg/pg"
	"time"
)

func Run() error {
	for true {
		if err := processUserTasks(); err != nil {
			return err
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

func CompleteTask(tx *pg.Tx, tableName, fieldName string, fieldValue interface{}) error {
	// TODO: refactor to be able to pass the actual model, not a table's name
	_, err := tx.Model().Exec(`
		UPDATE `+tableName+` 
		SET is_done = true
		WHERE is_done = false AND `+fieldName+` = ?
	`, fieldValue)

	return err
}
