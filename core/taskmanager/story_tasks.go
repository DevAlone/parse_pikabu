package taskmanager

import (
	"sync"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/globals"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg"
)

func storyTasksWorker() error {
	var wg sync.WaitGroup

	for _, f := range []func() error{
		addMissingStoriesWorker,
		addNewStoriesWorker,
		updateStoriesWorker,
	} {
		wg.Add(1)
		go func(handler func() error) {
			helpers.PanicOnError(handler())
			wg.Done()
		}(f)
	}

	wg.Wait()

	return nil
}

func addMissingStoriesWorker() error {
	for {
		count, err := models.Db.Model((*models.PikabuStory)(nil)).Count()
		if err != nil {
			return errors.New(err)
		}
		if count == 0 {
			// init database
			storyID := uint64(6579293)
			err := AddParseStoryTask(storyID)
			if err != nil {
				return err
			}
			time.Sleep(time.Duration(config.Settings.WaitBeforeAddingNewStoryTasksSeconds) * time.Second)
			continue
		}

		// TODO: parse gaps

		time.Sleep(time.Duration(config.Settings.WaitBeforeAddingNewStoryTasksSeconds) * time.Second)
	}
}

func addNewStoriesWorker() error {
	for {
		time.Sleep(10 * time.Second)

		var lastStory models.PikabuStory
		err := models.Db.Model(&lastStory).
			Order("pikabu_id DESC").
			Limit(1).
			Select()
		if err == pg.ErrNoRows {
			continue
		}
		if err != nil {
			return errors.New(err)
		}
		for i := 0; i < config.Settings.NumberOfNewStoriesGap; i++ {
			storyID := lastStory.PikabuID + 1 + uint64(i)
			err = AddParseStoryTask(storyID)
			if err != nil {
				return err
			}
		}

		var firstStory models.PikabuStory
		err = models.Db.Model(&firstStory).
			Order("pikabu_id ASC").
			Limit(1).
			Select()
		if err == pg.ErrNoRows {
			continue
		}
		if err != nil {
			return errors.New(err)
		}
		for i := 0; i < config.Settings.NumberOfNewStoriesGap; i++ {
			storyID := firstStory.PikabuID - 1 - uint64(i)
			if storyID == 0 {
				break
			}
			err = AddParseStoryTask(storyID)
			if err != nil {
				return err
			}
		}
		time.Sleep(time.Duration(config.Settings.WaitBeforeAddingNewStoryTasksSeconds) * time.Second)
	}
}

func updateStoriesWorker() error {
	for {
		time.Sleep(time.Duration(config.Settings.WaitBeforeAddingNewStoryTasksSeconds) * time.Second)

		if len(globals.ParserParseStoryTasks) >= config.Settings.MaxNumberOfTasksInQueue/2 {
			// wait for queue to become empty
			continue
		}

		var storiesToUpdate []models.PikabuStory
		err := models.Db.Model(&storiesToUpdate).
			Where("is_permanently_deleted = false AND next_update_timestamp < ? AND task_taken_at_timestamp < ?", time.Now().Unix(), time.Now().Unix()-int64(config.Settings.MaximumParseStoryTaskProcessingTime)).
			Limit(1024).
			Select()
		if err != pg.ErrNoRows && err != nil {
			return err
		}
		for _, story := range storiesToUpdate {
			err := AddParseStoryTask(story.PikabuID)
			if err != nil {
				return err
			}
		}

		var deletedOrNeverExistedStories []models.PikabuDeletedOrNeverExistedStory
		err = models.Db.Model(&deletedOrNeverExistedStories).
			Where("next_update_timestamp < ? AND task_taken_at_timestamp > ?", time.Now().Unix(), time.Now().Unix()-int64(config.Settings.MaximumParseStoryTaskProcessingTime)).
			Limit(1024).
			Select()
		if err != pg.ErrNoRows && err != nil {
			return err
		}
		for _, story := range storiesToUpdate {
			err := AddParseStoryTask(story.PikabuID)
			if err != nil {
				return err
			}
		}
	}
}

// AddParseStoryTask queues task for parsing story
func AddParseStoryTask(pikabuID uint64) error {
	timestamp := models.TimestampType(time.Now().Unix())

	deletedOrNeverExistedStory := models.PikabuDeletedOrNeverExistedStory{
		PikabuID: pikabuID,
	}
	err := models.Db.Select(&deletedOrNeverExistedStory)
	if err != pg.ErrNoRows && err != nil {
		return errors.New(err)
	}
	if err == nil {
		// ignore recently added tasks
		if deletedOrNeverExistedStory.TaskTakenAtTimestamp+models.TimestampType(config.Settings.MaximumParseStoryTaskProcessingTime) >= timestamp {
			return nil
		}
		deletedOrNeverExistedStory.TaskTakenAtTimestamp = timestamp
		err := models.Db.Update(&deletedOrNeverExistedStory)
		if err != nil {
			return errors.New(err)
		}
	} else {
		deletedOrNeverExistedStory.LastUpdateTimestamp = 0
		deletedOrNeverExistedStory.NextUpdateTimestamp = timestamp
		deletedOrNeverExistedStory.TaskTakenAtTimestamp = timestamp
		err := models.Db.Insert(&deletedOrNeverExistedStory)
		if err != nil {
			return errors.New(err)
		}
	}

	story := models.PikabuStory{
		PikabuID: pikabuID,
	}
	err = models.Db.Select(&story)
	if err != pg.ErrNoRows && err != nil {
		return errors.New(err)
	}
	if err == nil {
		// ignore recently added tasks
		if story.TaskTakenAtTimestamp+models.TimestampType(config.Settings.MaximumParseStoryTaskProcessingTime) >= timestamp {
			return nil
		}

		story.TaskTakenAtTimestamp = timestamp
		err := models.Db.Update(&story)
		if err != nil {
			return errors.New(err)
		}
	}

	return PushTaskToQueue(&models.ParseStoryTask{
		PikabuID:       pikabuID,
		AddedTimestamp: timestamp,
	})
}
