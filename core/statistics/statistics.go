package statistics

import (
	"sync"
	"time"

	"github.com/DevAlone/parse_pikabu/core/logger"
	"github.com/DevAlone/parse_pikabu/core/taskmanager"
	"github.com/DevAlone/parse_pikabu/globals"
	"github.com/DevAlone/parse_pikabu/helpers"
	"github.com/DevAlone/parse_pikabu/models"
	"github.com/go-errors/errors"
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
		// TODO: comment out?
		{
			TableName:             "pikabu_comments",
			StatTableName:         "number_of_comments_to_process_entries",
			UpdatingPeriodSeconds: 12 * 60 * 60,
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

	wg.Add(1)
	go func() {
		err := processTaskManagerChannelsSize()
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

func processTaskManagerChannelsSize() error {
	for {
		if taskmanager.CoreTaskManager == nil {
			time.Sleep(time.Second)
			continue
		}

		for channelName, channelSize := range map[string]int64{
			"update_user_tasks":                          taskmanager.CoreTaskManager.MustTaskChannelSize(taskmanager.UpdateUserTask),
			"parse_new_user_tasks":                       taskmanager.CoreTaskManager.MustTaskChannelSize(taskmanager.ParseNewUserTask),
			"parse_deleted_or_never_existed_user_tasks":  taskmanager.CoreTaskManager.MustTaskChannelSize(taskmanager.ParseDeletedOrNeverExistedUserTask),
			"update_story_tasks":                         taskmanager.CoreTaskManager.MustTaskChannelSize(taskmanager.UpdateStoryTask),
			"parse_new_story_tasks":                      taskmanager.CoreTaskManager.MustTaskChannelSize(taskmanager.ParseNewStoryTask),
			"parse_deleted_or_never_existed_story_tasks": taskmanager.CoreTaskManager.MustTaskChannelSize(taskmanager.ParseDeletedOrNeverExistedStoryTask),
			"parse_all_communities_tasks":                taskmanager.CoreTaskManager.MustTaskChannelSize(taskmanager.ParseAllCommunitiesTask),
			"parser_results":                             int64(len(globals.ParserResults)),
		} {
			err := processTaskManagerChannelSize(channelName, channelSize)
			if err != nil {
				return err
			}
		}

		time.Sleep(time.Minute)
	}
}

func processTaskManagerChannelSize(channelName string, channelSize int64) error {
	tableName := "number_of_items_in_channel_" + channelName
	query := `
		WITH constants (curr_timestamp, value) AS (
			VALUES (
	            extract(epoch from now())::int,
				?0
			)
		) INSERT INTO ` + tableName + ` (timestamp, value)
		SELECT 
			constants.curr_timestamp, 
			constants.value
		FROM constants
		ON CONFLICT (timestamp) DO NOTHING;
	`
	_, err := models.Db.Exec(query, channelSize)
	if err != nil {
		return errors.New(err)
	}
	return nil
}
