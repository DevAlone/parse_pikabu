package task_manager

import (
	"bitbucket.org/d3dev/parse_pikabu/models"
	"strings"
	"time"
)

func processUserTasks() error {
	for _, username := range []string{
		"admin",
		"l4rever",
	} {
		// print(username, "\n")
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
		} else {
			// check for update
		}
	}

	for _, id := range []uint64{1, 2, 3, 4} {
		err := AddParseUserByIdTask(id)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddParseUserByUsernameTask(username string) error {
	username = strings.ToLower(username)

	exists, err := models.Db.Model(&models.ParseUserByUsernameTask{}).
		Where("LOWER(username) = ?", username).
		Exists()

	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	task := &models.ParseUserByUsernameTask{
		Task: models.Task{
			AddedTimestamp: models.TimestampType(time.Now().Unix()),
		},
		Username: username,
	}

	return models.Db.Insert(task)
}

func AddParseUserByIdTask(id uint64) error {
	exists, err := models.Db.Model(&models.ParseUserByIdTask{}).
		Where("pikabu_id = ?", id).
		Exists()

	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	task := &models.ParseUserByIdTask{
		Task: models.Task{
			AddedTimestamp: models.TimestampType(time.Now().Unix()),
		},
		PikabuId: id,
	}

	return models.Db.Insert(task)
}
