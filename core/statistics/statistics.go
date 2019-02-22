package statistics

import (
	"fmt"
	"sync"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/models"
)

func Run() error {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		err := ProcessNumberOfUsersInQueue()
		helpers.PanicOnError(err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := ProcessDistributions()
		helpers.PanicOnError(err)
		wg.Done()
	}()

	wg.Wait()

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

func ProcessDistributions() error {
	for true {
		err := ProcessDistribution("pikabu_users", "signup_timestamp", 86400)
		if err != nil {
			return err
		}

		time.Sleep(10 * time.Minute)
	}

	return nil
}

func ProcessDistribution(tableName string, columnName string, bucketSize int) error {
	distributionTableName := tableName + "_" + columnName + "_distribution_" + fmt.Sprint(bucketSize)
	_, err := models.Db.Exec(`
DELETE FROM ` + distributionTableName + `; 
INSERT INTO ` + distributionTableName + `
SELECT timestamp, value FROM (
	WITH stats AS (
		SELECT MIN(` + columnName + `) as min_value, MAX(` + columnName + `) as max_value
		FROM ` + tableName + `
	)
	SELECT 
		width_bucket(` + columnName + `::int, min_value::int, max_value::int, ((max_value - min_value) / (` + fmt.Sprint(bucketSize) + `))::int) as bucket,
		MIN(` + columnName + `) as timestamp,
		COUNT(*) AS value
	FROM 
		` + tableName + `, stats
	GROUP BY
		bucket
	ORDER BY
		bucket
) distribution;
`)
	return err
}
