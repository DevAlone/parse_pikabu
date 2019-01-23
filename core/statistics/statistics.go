package statistics

import (
	"sync"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
)

func Run() error {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		ProcessNumberOfUsersInQueue()
		wg.Done()
	}()

	wg.Done()

	return nil
}

func ProcessNumberOfUsersInQueue() error {
	// number_of_users_to_process
	for true {
		_, err := models.Db.Exec(`
			WITH constants (curr_timestamp) AS (
				VALUES (extract(epoch from now())::int)
			) INSERT INTO number_of_users_to_process_entries (timestamp, value)
			SELECT 
				constants.curr_timestamp, 
				(
					SELECT COUNT(*) FROM pikabu_users
					WHERE 
						next_update_timestamp <= constants.curr_timestamp
				)
			FROM constants
			ON CONFLICT (timestamp) DO NOTHING;
		`)
		if err != nil {
			logger.Log.Error(err)
		}

		time.Sleep(60 * time.Second)
	}

	return nil
}
