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
	for {
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
}

// ProcessDistribution -
func ProcessDistribution(tableName string, columnName string, bucketSize int) error {
	distributionTableName := tableName + "_" + columnName + "_distribution_" + fmt.Sprint(bucketSize)
	if columnName == "updating_period" {
		columnName = "(next_update_timestamp - last_update_timestamp)"
	}
	tableName += "s"
	request := `
DELETE FROM ` + distributionTableName + `; 
INSERT INTO ` + distributionTableName + `
SELECT timestamp, value FROM (
	WITH stats_min_max AS (
		SELECT MIN(` + columnName + `) as min_value, MAX(` + columnName + `) as max_value
		FROM ` + tableName + `
	), stats_values_range AS (
        SELECT ((max_value - min_value) / (86400))::int AS values_range
        FROM stats_min_max
    )
	SELECT 
		width_bucket(
            ` + columnName + `::int, 
			min_value::int, 
			CASE 
                WHEN max_value::int = min_value::int THEN max_value::int + 1
                ELSE max_value::int
            END,
			CASE 
                WHEN values_range = 0 THEN 1
                ELSE values_range
            END
		) as bucket,
		MIN(` + columnName + `) as timestamp,
		COUNT(*) AS value
	FROM 
		` + tableName + `, stats_min_max, stats_values_range
	GROUP BY
		bucket
	ORDER BY
		bucket
) distribution;
`
	_, err := models.Db.Exec(request)
	if err != nil {
		logger.Log.Errorf("error during execution request %v", request)
		return errors.New(err)
	}
	return nil
}
