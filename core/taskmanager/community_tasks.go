package taskmanager

import (
	"time"

	"github.com/DevAlone/parse_pikabu/core/config"
	"github.com/DevAlone/parse_pikabu/helpers"
	"github.com/DevAlone/parse_pikabu/models"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg"
	"github.com/go-redis/redis"
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
		err := AddParseCommunitiesTask()
		if err != nil {
			return err
		}

		lastTimestamp := time.Now().Unix()
		err = redisClient.Set(redisKey, lastTimestamp, 0).Err()
		if err != nil {
			return errors.New(err)
		}
	}

	return nil
}

// AddParseCommunitiesTask -
func AddParseCommunitiesTask() error {
	task := &models.SimpleTask{
		Name: "parse_communities_pages",
	}
	err := models.Db.Model(task).
		Where("name = ?name").
		Select()
	if err != pg.ErrNoRows && err != nil {
		return errors.New(err)
	}

	exists := err != pg.ErrNoRows

	task.AddedTimestamp = models.TimestampType(time.Now().Unix())
	task.IsDone = false
	task.IsTaken = true

	if !exists {
		err := models.Db.Insert(task)
		if err != nil {
			return errors.New(err)
		}
	} else {
		err := models.Db.Update(task)
		if err != nil {
			return errors.New(err)
		}
	}

	return CoreTaskManager.PushTask(ParseAllCommunitiesTask, task)
}
