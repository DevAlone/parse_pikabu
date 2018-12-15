package task_manager

import "time"

func Run() error {
	for true {
		if err := processUserTasks(); err != nil {
			return err
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}
