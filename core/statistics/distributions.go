package statistics

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DevAlone/parse_pikabu/core/logger"
	"github.com/DevAlone/parse_pikabu/helpers"
	"github.com/DevAlone/parse_pikabu/models"
	"github.com/go-errors/errors"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

const updatingPeriodSeconds = 12 * 3600

// ProcessDistributions -
func ProcessDistributions() error {
	for {
		redisClient := helpers.GetRedisClient()

		err := redisClient.SetNX("parse_pikabu/core/statistics/distributions/last_update_timestamp", 0, 0).Err()
		if err != nil {
			return err
		}

		lastUpdateTimestampStr, err := redisClient.Get("parse_pikabu/core/statistics/distributions/last_update_timestamp").Result()
		if err != nil {
			return err
		}
		lastUpdateTimestamp, err := strconv.ParseInt(lastUpdateTimestampStr, 10, 64)
		if err != nil {
			return err
		}

		if time.Now().Unix() < lastUpdateTimestamp+updatingPeriodSeconds {
			time.Sleep(1 * time.Hour)
			continue
		}

		// TODO: save last timestamp on redis
		// time.Sleep(10 * time.Minute)
		for _, distributionFieldModel := range models.GeneratedDistributionFields {
			baseTableNameSnakeCase := strcase.ToSnake(distributionFieldModel.BaseTableName)
			baseColumnNameSnakeCase := strcase.ToSnake(distributionFieldModel.BaseColumnName)
			err := ProcessDistribution(baseTableNameSnakeCase, baseColumnNameSnakeCase, distributionFieldModel.BucketSize)
			if err != nil {
				return err
			}
		}

		err = ProcessDistribution("pikabu_user", "updating_period", 3600)
		if err != nil {
			return err
		}

		err = redisClient.Set("parse_pikabu/core/statistics/distributions/last_update_timestamp", time.Now().Unix(), 0).Err()
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
	tableName = inflection.Plural(tableName)
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
