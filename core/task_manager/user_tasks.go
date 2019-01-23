package task_manager

import (
	"strings"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-pg/pg"
)

func processUserTasks() error {
	// update tasks
	parseUserByUsernameTasks := []models.ParseUserByUsernameTask{}
	err := models.Db.Model(&parseUserByUsernameTasks).
		Where(
			"is_done = false AND is_taken = true AND added_timestamp < ?",
			models.TimestampType(time.Now().Unix())-models.TimestampType(config.Settings.MaximumTaskProcessingTime)).
		Limit(1024).
		Select()
	if err != pg.ErrNoRows && err != nil {
		return err
	}
	for _, task := range parseUserByUsernameTasks {
		err := AddParseUserByUsernameTask(task.Username)
		if err != nil {
			return err
		}
	}

	// update users
	users := []models.PikabuUser{}
	err = models.Db.Model(&users).
		Where("next_update_timestamp <= ?", time.Now().Unix()).
		Order("next_update_timestamp").
		Select()
	if err != pg.ErrNoRows && err != nil {
		return err
	}

	for _, user := range users {
		err = AddParseUserByUsernameTask(user.Username)
		if err != nil {
			return err
		}
	}

	// init db
	/*
		for _, username := range []string{
			"admin",
			"l4rever",
			"moderator",
			"lactarius",
			"apres",
			"dev",
			"code501",
		} {
			user := &models.PikabuUser{}
			exists, err := models.Db.Model(user).
				Where("LOWER(username) = ?", strings.ToLower(username)).
				Exists()

			if err != nil {
				return err
			}

			if !exists {
				err = AddParseUserByUsernameTask(username)
				if err != nil {
					return err
				}
			}
		}
	*/

	return nil
}

func AddParseUserByUsernameTask(username string) error {
	username = strings.ToLower(username)

	task := &models.ParseUserByUsernameTask{}

	err := models.Db.Model(task).
		Where("LOWER(username) = ?", username).
		Select()

	if err != pg.ErrNoRows && err != nil {
		return err
	}

	exists := err != pg.ErrNoRows

	task.AddedTimestamp = models.TimestampType(time.Now().Unix())
	task.IsDone = false
	task.IsTaken = true
	task.Username = username

	if !exists {
		err := models.Db.Insert(task)
		if err != nil {
			return err
		}
	} else {
		err := models.Db.Update(task)
		if err != nil {
			return err
		}
	}

	return PushTaskToQueue(task)
}

// TODO: unused
func AddParseUserByIdTask(id uint64) error {
	exists, err := models.Db.Model(&models.ParseUserByIdTask{}).
		Where("pikabu_id = ?", id).
		Exists()

	if err != nil || exists {
		return err
	}

	task := &models.ParseUserByIdTask{
		Task: models.Task{
			AddedTimestamp: models.TimestampType(time.Now().Unix()),
		},
		PikabuId: id,
	}

	return models.Db.Insert(task)
}
