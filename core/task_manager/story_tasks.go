package task_manager

import (
	"time"

	"github.com/go-errors/errors"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-pg/pg"
)

func addMissingStoriesWorker() error {
	for true {
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

		// parse new stories

		time.Sleep(time.Duration(config.Settings.WaitBeforeAddingTasksSeconds) * time.Second)
	}

	return nil
}

func storyTasksWorker() error {
	// TODO: do it
	for true {
		// update tasks
		/*parseStoryTasks := []models.ParseUserTask{}
		err = models.Db.Model(&parseUserTasks).
			Where(
				"is_done = false AND is_taken = true AND added_timestamp < ?",
				models.TimestampType(time.Now().Unix())-models.TimestampType(config.Settings.MaximumTaskProcessingTime)).
			Limit(1024).
			Select()
		if err != pg.ErrNoRows && err != nil {
			return errors.New(err)
		}
		for _, task := range parseUserTasks {
			err := AddParseUserTask(task.PikabuId, task.Username)
			if err != nil {
				return err
			}
		}
		*/

		// update stories

		time.Sleep(time.Duration(config.Settings.WaitBeforeAddingTasksSeconds) * time.Second)
	}

	return nil
}

func AddParseStoryTask(pikabuId uint64) error {
	// TODO: fix
	// username = strings.ToLower(username)

	task := &models.ParseUserTask{}

	err := models.Db.Model(task).
		Where("pikabu_id = ?", pikabuId).
		Select()

	if err != pg.ErrNoRows && err != nil {
		return errors.New(err)
	}

	exists := err != pg.ErrNoRows

	task.PikabuId = pikabuId
	task.AddedTimestamp = models.TimestampType(time.Now().Unix())
	task.IsDone = false
	task.IsTaken = true
	task.Username = username

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

	return PushTaskToQueue(task)
}

func AddParseStoryTaskIfNotExists(pikabuId uint64) error {
	// TODO: fix
	task := &models.ParseUserTask{}

	err := models.Db.Model(task).
		Where("pikabu_id = ?", pikabuId).
		Select()

	if err != pg.ErrNoRows && err != nil {
		return errors.New(err)
	}

	exists := err != pg.ErrNoRows
	if exists {
		return nil
	}

	task.PikabuId = pikabuId
	task.AddedTimestamp = models.TimestampType(time.Now().Unix())
	task.IsDone = false
	task.IsTaken = true
	task.Username = username

	err = models.Db.Insert(task)
	if err != nil {
		return errors.New(err)
	}

	return PushTaskToQueue(task)
}
