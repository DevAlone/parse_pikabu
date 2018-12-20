package task_manager

import (
	"bitbucket.org/d3dev/parse_pikabu/config"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg"
	"github.com/go-redis/redis"
	"time"
)

func processCommunityTasks() error {
	redisClient := helpers.GetRedisClient()
	redisKey := "parse_pikabu_core_tasks_communities_last_processing_timestamp"
	lastTimestamp, err := redisClient.Get(redisKey).Int64()

	if err != redis.Nil && err != nil {
		return errors.New(err)
	}
	if err == redis.Nil {
		lastTimestamp = 0
	}

	if lastTimestamp < time.Now().Unix()-int64(config.Settings.CommunitiesProcessingPeriod) {
		err := processCommunities()
		if err != nil {
			return errors.New(err)
		}
		lastTimestamp := time.Now().Unix()
		err = redisClient.Set(redisKey, lastTimestamp, 0).Err()
		if err != nil {
			return errors.New(err)
		}
	}

	return nil
}

func processCommunities() error {
	task := &models.SimpleTask{
		Name: "parse_communities",
	}
	err := models.Db.Model(task).
		Where("name = ?name").
		Select()
	if err != pg.ErrNoRows && err != nil {
		return err
	}
	isExpired := task.AddedTimestamp <
		models.TimestampType(time.Now().Unix())-
			models.TimestampType(config.Settings.MaximumTaskProcessingTime)

	if err == pg.ErrNoRows || isExpired {
		task.AddedTimestamp = models.TimestampType(time.Now().Unix())
		task.IsDone = false
		task.IsTaken = false

		if err == pg.ErrNoRows {
			return models.Db.Insert(task)
		} else {
			return models.Db.Update(task)
		}
	}

	return nil
}
