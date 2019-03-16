package taskmanager

import (
	"sync"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-pg/pg"
)

func storyTasksWorker() error {
	var wg sync.WaitGroup

	for _, f := range []func() error{
		addMissingStoriesWorker,
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
			return err
		}
		if count == 0 {
			// init database
			storyID := uint64(6570769)
			err := AddParseStoryTask(storyID)
			if err != nil {
				return err
			}
		}

		// TODO: parse gaps

		// TODO: parse new stories
		{
			var lastStory models.PikabuStory
			err := models.Db.Model(lastStory).
				Order("pikabu_id DESC").
				Limit(1).
				Select()
			if err != nil {
				return err
			}
			for i := 0; i < config.Settings.NumberOfNewStoriesGap; i++ {
				storyID := lastStory.PikabuID + 1 + uint64(i)
				err = AddParseStoryTask(storyID)
				if err != nil {
					return err
				}
			}
		}

		time.Sleep(time.Duration(config.Settings.WaitBeforeAddingNewStoryTasksSeconds) * time.Second)
	}
}

func updateStoriesWorker() error {
	for {
		// TODO: update stories

		time.Sleep(time.Duration(config.Settings.WaitBeforeAddingTasksSeconds) * time.Second)
	}
}

// AddParseStoryTask queues task for parsing story
func AddParseStoryTask(pikabuID uint64) error {
	timestamp := models.TimestampType(time.Now().Unix())

	deletedOrNeverExistedStory := models.PikabuStory{
		PikabuID: pikabuID,
	}
	err := models.Db.Select(&deletedOrNeverExistedStory)
	if err != pg.ErrNoRows && err != nil {
		return err
	}
	if err == nil {
		// ignore recently added tasks
		if deletedOrNeverExistedStory.TaskTakenAtTimestamp+models.TimestampType(config.Settings.MaximumParseStoryTaskProcessingTime) >= timestamp {
			return nil
		}
		deletedOrNeverExistedStory.TaskTakenAtTimestamp = timestamp
		err := models.Db.Update(&deletedOrNeverExistedStory)
		if err != nil {
			return err
		}
	} else {
		deletedOrNeverExistedStory.LastUpdateTimestamp = 0
		deletedOrNeverExistedStory.NextUpdateTimestamp = timestamp
		deletedOrNeverExistedStory.TaskTakenAtTimestamp = timestamp
		err := models.Db.Insert(&deletedOrNeverExistedStory)
		if err != nil {
			return err
		}
	}

	story := models.PikabuStory{
		PikabuID: pikabuID,
	}
	err = models.Db.Select(&story)
	if err != pg.ErrNoRows && err != nil {
		return err
	}
	if err == nil {
		// ignore recently added tasks
		if story.TaskTakenAtTimestamp+models.TimestampType(config.Settings.MaximumParseStoryTaskProcessingTime) >= timestamp {
			return nil
		}

		story.TaskTakenAtTimestamp = timestamp
		err := models.Db.Update(&story)
		if err != nil {
			return err
		}
	}

	return PushTaskToQueue(&models.ParseStoryTask{
		PikabuID:       pikabuID,
		AddedTimestamp: timestamp,
	})
}
