package statistics

import (
	"fmt"
	"sync"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-errors/errors"
	"github.com/iancoleman/strcase"
)

// Run -
func Run() error {
	var wg sync.WaitGroup

	type ProcessNumberOfItemsInQueueConfig struct {
		TableName             string
		StatTableName         string
		UpdatingPeriodSeconds int
	}

	for _, config := range []ProcessNumberOfItemsInQueueConfig{
		{
			TableName:             "pikabu_users",
			StatTableName:         "number_of_users_to_process_entries",
			UpdatingPeriodSeconds: 60,
		},
		{
			TableName:             "pikabu_stories",
			StatTableName:         "number_of_stories_to_process_entries",
			UpdatingPeriodSeconds: 10 * 60,
		},
		{
			TableName:             "pikabu_comments",
			StatTableName:         "number_of_comments_to_process_entries",
			UpdatingPeriodSeconds: 60 * 60,
		},
	} {
		wg.Add(1)
		go func(config ProcessNumberOfItemsInQueueConfig) {
			err := ProcessNumberOfItemsInQueue(config.TableName, config.StatTableName, config.UpdatingPeriodSeconds)
			helpers.PanicOnError(err)
			wg.Done()
		}(config)
	}

	wg.Add(1)
	go func() {
		err := ProcessDistributions()
		helpers.PanicOnError(err)
		wg.Done()
	}()

	wg.Wait()

	return nil
}

// ProcessNumberOfItemsInQueue -
func ProcessNumberOfItemsInQueue(tableName string, statTableName string, updatingPeriod int) error {
	for {
		query := `
			WITH constants (curr_timestamp) AS (
				VALUES (extract(epoch from now())::int)
			) INSERT INTO ` + statTableName + ` (timestamp, value)
			SELECT 
				constants.curr_timestamp, 
				(
					SELECT COUNT(*) FROM ` + tableName + `
					WHERE 
						next_update_timestamp <= constants.curr_timestamp
				)
			FROM constants
			ON CONFLICT (timestamp) DO NOTHING;
		`
		_, err := models.Db.Exec(query)
		if err != nil {
			logger.Log.Errorf("error during processing query %v\n", query)
			logger.LogError(errors.New(err))
		}

		time.Sleep(time.Duration(updatingPeriod) * time.Second)
	}
}

// ProcessDistributions -
func ProcessDistributions() error {
	for true {
		// TODO: save last timestamp on redis
		time.Sleep(10 * time.Minute)
		for _, distributionFieldModel := range models.GeneratedDistributionFields {
			baseTableNameSnakeCase := strcase.ToSnake(distributionFieldModel.BaseTableName)
			baseColumnNameSnakeCase := strcase.ToSnake(distributionFieldModel.BaseColumnName)
			err := ProcessDistribution(baseTableNameSnakeCase, baseColumnNameSnakeCase, distributionFieldModel.BucketSize)
			if err != nil {
				return err
			}
		}

		err := ProcessDistribution("pikabu_user", "updating_period", 3600)
		if err != nil {
			return err
		}

		time.Sleep(1 * time.Hour)
	}

	return nil
}

// ProcessDistribution -
func ProcessDistribution(tableName string, columnName string, bucketSize int) error {
	distributionTableName := tableName + "_" + columnName + "_distribution_" + fmt.Sprint(bucketSize)
	if columnName == "updating_period" {
		columnName = "(next_update_timestamp - last_update_timestamp)"
	}
	tableName += "s"
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
