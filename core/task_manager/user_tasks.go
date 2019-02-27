package task_manager

import (
	"time"

	"github.com/go-errors/errors"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-pg/pg"
)

func addMissingTasksWorker() error {
	for true {
		var users []models.PikabuUser
		// very slow query
		err := models.Db.Model(&users).
			ColumnExpr("pikabu_user.*").
			Join("LEFT JOIN parse_user_tasks AS parse_user_task").
			JoinOn("pikabu_user.pikabu_id = parse_user_task.pikabu_id").
			Where("parse_user_task.pikabu_id IS NULL").
			Limit(1024).
			Select()
		if err != nil {
			return err
		}
		for _, user := range users {
			err := AddParseUserTask(user.PikabuId, user.Username)
			if err != nil {
				return err
			}
		}

		time.Sleep(1 * time.Hour)
	}

	return nil
}

func addMissingUsersWorker() error {
	for true {
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
WHERE pikabu_id + 1 <> next_nr LIMIT 10;
`, offset)
			if err == pg.ErrNoRows {
				break
			} else if err != nil {
				return err
			}

			for _, gap := range gaps {
				if gap.GapEnd > offset {
					offset = gap.GapEnd
				}

				for i := gap.GapStart; i <= gap.GapEnd; i++ {
					deletedUser := models.PikabuDeletedOrNeverExistedUser{
						PikabuId:            i,
						LastUpdateTimestamp: 0,
						NextUpdateTimestamp: 0,
					}
					_, err := models.Db.Model(&deletedUser).
						OnConflict("(pikabu_id) DO NOTHING").
						Insert()
					if err != nil {
						return err
					}
					err = AddParseUserTaskIfNotExists(deletedUser.PikabuId, "")
					if err != nil {
						return err
					}
				}
			}
		}

		// parse new users
		// set offset to max value
		var lastUser models.PikabuUser
		err := models.Db.Model(&lastUser).
			Order("pikabu_id DESC").
			Limit(1).
			Select()
		if err != nil {
			return err
		}

		for i := 0; i < config.Settings.NumberOfNewUsersGap; i++ {
			deletedUser := models.PikabuDeletedOrNeverExistedUser{
				PikabuId:            lastUser.PikabuId + 1 + uint64(i),
				LastUpdateTimestamp: 0,
				NextUpdateTimestamp: 0,
			}
			_, err := models.Db.Model(&deletedUser).
				OnConflict("(pikabu_id) DO NOTHING").
				Insert()
			if err != nil {
				return err
			}
			err = AddParseUserTaskIfNotExists(deletedUser.PikabuId, "")
			if err != nil {
				return err
			}
		}

		// try to parse again
		var deletedUsers []models.PikabuDeletedOrNeverExistedUser
		err = models.Db.Model(&deletedUsers).
			Where("next_update_timestamp <= ?", time.Now().Unix()).
			Limit(1024).
			Select()
		for _, deletedUser := range deletedUsers {
			err := AddParseUserTask(deletedUser.PikabuId, "")
			if err != nil {
				return err
			}
		}

		time.Sleep(30 * time.Minute)
	}

	return nil
}

func processUserTasks() error {
	// update users
	// TODO: iterate over all users here
	// TODO: improve performance
	users := []models.PikabuUser{}
	err := models.Db.Model(&users).
		ColumnExpr("pikabu_user.*").
		Join("LEFT JOIN parse_user_tasks AS parse_user_task").
		JoinOn("pikabu_user.pikabu_id = parse_user_task.pikabu_id").
		Where("next_update_timestamp <= ? AND parse_user_task.is_done = true", time.Now().Unix()).
		Order("next_update_timestamp").
		Limit(1024).
		Select()
	if err != pg.ErrNoRows && err != nil {
		return errors.New(err)
	}

	for _, user := range users {
		err = AddParseUserTask(user.PikabuId, user.Username)
		if err != nil {
			return err
		}
	}

	// update tasks
	parseUserTasks := []models.ParseUserTask{}
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

	return nil
}

func AddParseUserTask(pikabuId uint64, username string) error {
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

func AddParseUserTaskIfNotExists(pikabuId uint64, username string) error {
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
