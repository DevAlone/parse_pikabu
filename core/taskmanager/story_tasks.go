package taskmanager

import (
	"sync"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
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
		updateDeletedOrNeverExistedStoriesWorker,
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
			for _, storyID := range []uint64{
				6579296,
				6579295,
				6579294,
				6579293,
				6579292,
				6579291,
				6579290,
			} {
				err := AddParseStoryTask(storyID, ParseNewStoryTask)
				if err != nil {
					return err
				}

			}
			time.Sleep(time.Duration(config.Settings.WaitBeforeAddingNewStoryTasksSeconds) * time.Second)
			continue
		}

		// parse gaps
		for offset := uint64(0); true; {
			var gaps []struct {
				GapStart uint64
				GapEnd   uint64
			}

			_, err := models.Db.Query(&gaps, `
SELECT
	pikabu_id + 1 as gap_start,
	next_nr - 1 as gap_end
FROM (
	SELECT 
		pikabu_id,
		lead(pikabu_id)
	OVER (ORDER BY pikabu_id) as next_nr
	FROM pikabu_stories
	WHERE pikabu_id > ?
) nr
WHERE pikabu_id + 1 <> next_nr LIMIT 10;
			`, offset)
			if err == pg.ErrNoRows {
				break
			} else if err != nil {
				return errors.New(err)
			}

			for _, gap := range gaps {
				if gap.GapEnd > offset {
					offset = gap.GapEnd
				} else {
					offset++
				}

				for i := gap.GapStart; i <= gap.GapEnd; i++ {
					err := AddParseStoryTaskIfNotExists(i, ParseDeletedOrNeverExistedStoryTask)
					if err != nil {
						return err
					}
				}
			}
		}

		time.Sleep(time.Duration(config.Settings.WaitBeforeAddingNewStoryTasksSeconds) * time.Second)
	}
}

func addNewStoriesWorker() error {
	for {
		time.Sleep(time.Minute)

		// go up
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
			err = AddParseStoryTask(storyID, ParseNewStoryTask)
			if err != nil {
				return err
			}
		}

		// go down
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
			err = AddParseStoryTask(storyID, ParseNewStoryTask)
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

		var storiesToUpdate []models.PikabuStory
		err := models.Db.Model(&storiesToUpdate).
			Where("next_update_timestamp < ? AND task_taken_at_timestamp < ?", time.Now().Unix(), time.Now().Unix()-int64(config.Settings.MaximumParseStoryTaskProcessingTime)).
			Order("next_update_timestamp").
			Limit(config.Settings.GetItemsToUpdateAtTime).
			Select()
		if err != pg.ErrNoRows && err != nil {
			return errors.New(err)
		}
		for _, story := range storiesToUpdate {
			err := AddParseStoryTask(story.PikabuID, UpdateStoryTask)
			if err != nil {
				return err
			}
		}
	}
}

func updateDeletedOrNeverExistedStoriesWorker() error {
	for {
		time.Sleep(time.Duration(config.Settings.WaitBeforeAddingNewStoryTasksSeconds) * time.Second)

		// deletedOrNeverExistedStoriesToUpdate
		var deletedOrNeverExistedStoriesToUpdate []models.PikabuDeletedOrNeverExistedStory
		err := models.Db.Model(&deletedOrNeverExistedStoriesToUpdate).
			Where("next_update_timestamp < ? AND task_taken_at_timestamp > ?", time.Now().Unix(), time.Now().Unix()-int64(config.Settings.MaximumParseStoryTaskProcessingTime)).
			Limit(config.Settings.GetItemsToUpdateAtTime).
			Select()
		if err != pg.ErrNoRows && err != nil {
			return errors.New(err)
		}
		for _, story := range deletedOrNeverExistedStoriesToUpdate {
			err := AddParseStoryTask(story.PikabuID, ParseDeletedOrNeverExistedStoryTask)
			if err != nil {
				return err
			}
		}
	}
}

// AddParseStoryTask queues task for parsing story
func AddParseStoryTask(pikabuID uint64, taskType int) error {
	return _AddParseStoryTask(pikabuID, taskType, false)
}

// AddParseStoryTaskIfNotExists queues task for parsing story
func AddParseStoryTaskIfNotExists(pikabuID uint64, taskType int) error {
	return _AddParseStoryTask(pikabuID, taskType, true)
}

func _AddParseStoryTask(pikabuID uint64, taskType int, ignoreIfExists bool) error {
	timestamp := models.TimestampType(time.Now().Unix())

	story := models.PikabuStory{
		PikabuID: pikabuID,
	}
	err := models.Db.Select(&story)
	if err != pg.ErrNoRows && err != nil {
		return errors.New(err)
	}

	if err == nil {
		// story exists
		if ignoreIfExists {
			return nil
		}

		// ignore recently added tasks
		if story.TaskTakenAtTimestamp+models.TimestampType(config.Settings.MaximumParseStoryTaskProcessingTime) >= timestamp {
			return nil
		}

		story.TaskTakenAtTimestamp = timestamp
		_, err := models.Db.Model(&story).Set("task_taken_at_timestamp = ?task_taken_at_timestamp").WherePK().Update()
		if err != nil {
			return errors.New(err)
		}
	} else {
		// story does not exist

		deletedOrNeverExistedStory := models.PikabuDeletedOrNeverExistedStory{
			PikabuID: pikabuID,
		}
		err := models.Db.Select(&deletedOrNeverExistedStory)
		if err != pg.ErrNoRows && err != nil {
			return errors.New(err)
		}
		if err == nil {
			// deletedOrNeverExistedStory exists
			if ignoreIfExists {
				return nil
			}

			// ignore recently added tasks
			if deletedOrNeverExistedStory.TaskTakenAtTimestamp+models.TimestampType(config.Settings.MaximumParseStoryTaskProcessingTime) >= timestamp {
				return nil
			}
			deletedOrNeverExistedStory.TaskTakenAtTimestamp = timestamp
			_, err := models.Db.Model(&deletedOrNeverExistedStory).Column("task_taken_at_timestamp").WherePK().Update()
			if err != pg.ErrNoRows && err != nil {
				return errors.New(err)
			}
		} else {
			// deletedOrNeverExistedStory does not exist
			deletedOrNeverExistedStory.LastUpdateTimestamp = 0
			deletedOrNeverExistedStory.NextUpdateTimestamp = timestamp
			deletedOrNeverExistedStory.TaskTakenAtTimestamp = timestamp
			err := models.Db.Insert(&deletedOrNeverExistedStory)
			if err != nil {
				return errors.New(err)
			}
		}
	}

	return CoreTaskManager.PushTask(taskType, &models.ParseStoryTask{
		PikabuID:       pikabuID,
		AddedTimestamp: timestamp,
	})
}

// ForceAddParseStoryTask queues task for parsing story without limiting
func ForceAddParseStoryTask(pikabuID uint64, taskType int) error {
	// TODO: refactor
	timestamp := models.TimestampType(time.Now().Unix())

	deletedOrNeverExistedStory := models.PikabuDeletedOrNeverExistedStory{
		PikabuID: pikabuID,
	}
	err := models.Db.Select(&deletedOrNeverExistedStory)
	if err != pg.ErrNoRows && err != nil {
		return errors.New(err)
	}
	if err == nil {
		// deletedOrNeverExistedStory exists

		deletedOrNeverExistedStory.TaskTakenAtTimestamp = timestamp
		err := models.Db.Update(&deletedOrNeverExistedStory)
		if err != pg.ErrNoRows && err != nil {
			return errors.New(err)
		}
	} else {
		// deletedOrNeverExistedStory does not exist
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
		story.TaskTakenAtTimestamp = timestamp
		err := models.Db.Update(&story)
		if err != nil {
			return errors.New(err)
		}
	}

	return CoreTaskManager.PushTask(taskType, &models.ParseStoryTask{
		PikabuID:       pikabuID,
		AddedTimestamp: timestamp,
	})
}
