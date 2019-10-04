package taskmanager

import (
	"sync"
	"time"

	"github.com/go-errors/errors"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-pg/pg"
)

func userTasksWorker() error {
	var wg sync.WaitGroup

	for _, f := range []func() error{
		addMissingUsersWorker,
		addNewUsersWorker,
		updateUsersWorker,
		updateDeletedOrNeverExistedUsersWorker,
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

func addMissingUsersWorker() error {
	for {
		time.Sleep(10 * time.Minute)

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
    FROM pikabu_users
    WHERE pikabu_id > ?
) nr 
WHERE pikabu_id + 1 <> next_nr LIMIT 128;
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
					err = AddParseUserTaskIfNotExists(i, "", ParseNewUserTask)
					if err != nil {
						return err
					}
				}
			}
		}

		time.Sleep(6 * time.Hour)
	}
}

func updateDeletedOrNeverExistedUsersWorker() error {
	for {
		time.Sleep(1 * time.Minute)

		// try to parse again
		var deletedUsers []models.PikabuDeletedOrNeverExistedUser
		err := models.Db.Model(&deletedUsers).
			Where("next_update_timestamp <= ?", time.Now().Unix()).
			Limit(1024).
			Select()
		if err == pg.ErrNoRows {
			continue
		}
		if err != nil {
			return errors.New(err)
		}
		for _, deletedUser := range deletedUsers {
			err := AddParseUserTask(deletedUser.PikabuID, "", ParseDeletedOrNeverExistedUserTask)
			if err != nil {
				return err
			}
		}

		time.Sleep(30 * time.Minute)
	}
}

func addNewUsersWorker() error {
	for {
		time.Sleep(time.Minute)

		// parse new users
		var lastUser models.PikabuUser
		err := models.Db.Model(&lastUser).
			Order("pikabu_id DESC").
			Limit(1).
			Select()
		if err == pg.ErrNoRows {
			err := AddParseUserTask(1, "admin", ParseNewUserTask)
			if err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return errors.New(err)
		}

		for i := 0; i < config.Settings.NumberOfNewUsersGap; i++ {
			deletedUser := models.PikabuDeletedOrNeverExistedUser{
				PikabuID:            lastUser.PikabuID + 1 + uint64(i),
				LastUpdateTimestamp: 0,
				NextUpdateTimestamp: 0,
			}
			_, err := models.Db.Model(&deletedUser).
				OnConflict("(pikabu_id) DO NOTHING").
				Insert()
			if err != nil {
				return err
			}
			err = AddParseUserTaskIfNotExists(deletedUser.PikabuID, "", ParseDeletedOrNeverExistedUserTask)
			if err != nil {
				return err
			}
		}

		time.Sleep(time.Duration(config.Settings.AddNewUsersEachNMinutes) * time.Minute)
	}
}

func updateUsersWorker() error {
	for {
		time.Sleep(time.Duration(config.Settings.UpdateUsersEachNSeconds) * time.Second)

		var usersToUpdate []models.PikabuUser
		err := models.Db.Model(&usersToUpdate).
			Where(
				"next_update_timestamp < ? AND task_taken_at_timestamp < ?",
				time.Now().Unix(),
				time.Now().Unix()-int64(config.Settings.MaximumParseUserTaskProcessingTime),
			).
			Order("next_update_timestamp").
			Limit(config.Settings.GetItemsToUpdateAtTime).
			Select()
		if err != pg.ErrNoRows && err != nil {
			return err
		}

		for _, user := range usersToUpdate {
			err = AddParseUserTask(user.PikabuID, user.Username, UpdateUserTask)
			if err != nil {
				return err
			}
		}
	}
}

// AddParseUserTask -
func AddParseUserTask(pikabuID uint64, username string, taskType int) error {
	timestamp := models.TimestampType(time.Now().Unix())

	user := &models.PikabuUser{
		PikabuID: pikabuID,
	}

	err := models.Db.Select(user)
	if err != pg.ErrNoRows && err != nil {
		return errors.New(err)
	}
	if err == nil {
		// such user exists

		// ignore recently added tasks
		if user.TaskTakenAtTimestamp+models.TimestampType(config.Settings.MaximumParseUserTaskProcessingTime) >= timestamp {
			return nil
		}

		user.TaskTakenAtTimestamp = timestamp
		_, err := models.Db.Model(user).Column("task_taken_at_timestamp").WherePK().Update()
		if err != nil {
			return errors.New(err)
		}
	} else {
		// deletedOrNeverExistedUser

		deletedOrNeverExistedUser := &models.PikabuDeletedOrNeverExistedUser{
			PikabuID: pikabuID,
		}
		err := models.Db.Select(deletedOrNeverExistedUser)
		if err != pg.ErrNoRows && err != nil {
			return errors.New(err)
		}
		if err == nil {
			// exists

			// ignore recently added tasks
			if deletedOrNeverExistedUser.TaskTakenAtTimestamp+models.TimestampType(config.Settings.MaximumParseUserTaskProcessingTime) >= timestamp {
				return nil
			}
			deletedOrNeverExistedUser.TaskTakenAtTimestamp = timestamp
			_, err := models.Db.Model(deletedOrNeverExistedUser).Column("task_taken_at_timestamp").WherePK().Update()
			if err != pg.ErrNoRows && err != nil {
				return errors.New(err)
			}
		} else {
			deletedOrNeverExistedUser.LastUpdateTimestamp = 0
			deletedOrNeverExistedUser.NextUpdateTimestamp = timestamp
			deletedOrNeverExistedUser.TaskTakenAtTimestamp = timestamp
			err := models.Db.Insert(deletedOrNeverExistedUser)
			if err != nil {
				return errors.New(err)
			}
		}
	}

	return CoreTaskManager.PushTask(taskType, &models.ParseUserTask{
		PikabuID:       pikabuID,
		Username:       username,
		AddedTimestamp: timestamp,
	})
}

// AddParseUserTaskIfNotExists -
func AddParseUserTaskIfNotExists(pikabuID uint64, username string, taskType int) error {
	user := &models.PikabuUser{
		PikabuID: pikabuID,
	}

	err := models.Db.Select(user)
	if err != pg.ErrNoRows && err != nil {
		return errors.New(err)
	}
	if err == nil {
		return nil
		// such user exists
	}

	deletedOrNeverExistedUser := &models.PikabuDeletedOrNeverExistedUser{
		PikabuID: pikabuID,
	}
	err = models.Db.Select(deletedOrNeverExistedUser)
	if err != pg.ErrNoRows && err != nil {
		return errors.New(err)
	}
	// deletedOrNeverExistedUser entry found
	if err == nil {
		return nil
	}

	timestamp := models.TimestampType(time.Now().Unix())

	deletedOrNeverExistedUser.LastUpdateTimestamp = 0
	deletedOrNeverExistedUser.NextUpdateTimestamp = timestamp
	deletedOrNeverExistedUser.TaskTakenAtTimestamp = timestamp
	err = models.Db.Insert(deletedOrNeverExistedUser)
	if err != nil {
		return errors.New(err)
	}

	return CoreTaskManager.PushTask(taskType, &models.ParseUserTask{
		PikabuID:       pikabuID,
		Username:       username,
		AddedTimestamp: timestamp,
	})
}
