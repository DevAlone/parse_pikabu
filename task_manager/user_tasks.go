package task_manager

import (
	"bitbucket.org/d3dev/parse_pikabu/config"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-pg/pg"
	"strings"
	"time"
)

func processUserTasks() error {
	// update tasks
	parseUserByUsernameTasks := []models.ParseUserByUsernameTask{}
	err := models.Db.Model(&parseUserByUsernameTasks).
		Where(
			"is_done = false AND is_taken = true AND added_timestamp < ?",
			models.TimestampType(time.Now().Unix())-models.TimestampType(config.Settings.MaximumTaskProcessingTime)).
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

	if exists {
		expired := task.AddedTimestamp < models.TimestampType(time.Now().Unix())-
			models.TimestampType(config.Settings.MaximumTaskProcessingTime)

		if !expired && !task.IsDone {
			return nil
		}

		task.IsTaken = false
		task.IsDone = false
		task.AddedTimestamp = models.TimestampType(time.Now().Unix())

		return models.Db.Update(task)
	}

	task = &models.ParseUserByUsernameTask{
		Task: models.Task{
			AddedTimestamp: models.TimestampType(time.Now().Unix()),
		},
		Username: username,
	}

	return models.Db.Insert(task)

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
