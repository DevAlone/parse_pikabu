package results_processor

import (
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/task_manager"
	"errors"
	"gogsweb.2-47.ru/d3dev/pikago"
	"sync"
)

var processUserProfileMutex = &sync.Mutex{}

func processUserProfile(userProfile *pikago.UserProfile) error {
	processUserProfileMutex.Lock()
	defer processUserProfileMutex.Unlock()

	if userProfile == nil {
		return errors.New("nil user profile")
	}
	tx, err := models.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// complete tasks
	err = task_manager.CompleteTask(
		tx,
		"parse_user_by_id_tasks",
		"pikabu_id",
		userProfile.UserId.Value,
	)
	if err != nil {
		return err
	}

	err = task_manager.CompleteTask(
		tx,
		"parse_user_by_username_tasks",
		"username",
		userProfile.Username,
	)
	if err != nil {
		return err
	}

	print("user profile\n")
	// save results
	return tx.Commit()
}
