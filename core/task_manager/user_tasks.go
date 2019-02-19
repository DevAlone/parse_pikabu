package task_manager

import (
	"time"

	"github.com/go-errors/errors"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-pg/pg"
)

func processUserTasks() error {
	// update tasks
	parseUserTasks := []models.ParseUserTask{}
	err := models.Db.Model(&parseUserTasks).
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

	// update users
	/*
		users := []models.PikabuUser{}
		err = models.Db.Model(&users).
			ColumnExpr("pikabu_user.*").
			Join("LEFT JOIN parse_user_tasks AS parse_user_task").
			JoinOn("pikabu_user.pikabu_id = parse_user_task.pikabu_id").
			Where("next_update_timestamp <= ? AND (parse_user_task.pikabu_id IS NULL OR parse_user_task.is_done = true)", time.Now().Unix()).
			Order("next_update_timestamp").
			Limit(1024).
			Select()
		if err != pg.ErrNoRows && err != nil {
			return errors.New(err)
		}
	*/
	users := []models.PikabuUser{}
	err = models.Db.Model(&users).
		Where("next_update_timestamp <= ?", time.Now().Unix()).
		Order("next_update_timestamp").
		Limit(1024).
		Select()

	for _, user := range users {
		err = AddParseUserTask(user.PikabuId, user.Username)
		if err != nil {
			return err
		}
	}

	// TODO: parse new users by their id

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
