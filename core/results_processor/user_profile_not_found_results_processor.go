package results_processor

import (
	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/core/task_manager"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-pg/pg"
)

func processUserProfileNotFoundResults(parsingTimestamp models.TimestampType, res []models.ParserUserProfileNotFoundResultData) error {
	for _, result := range res {
		err := processUserProfileNotFoundResult(parsingTimestamp, &result)
		if err != nil {
			// TODO: log
			return err
		}
	}

	return nil
}

func processUserProfileNotFoundResult(parsingTimestamp models.TimestampType, res *models.ParserUserProfileNotFoundResultData) error {
	lockUserById(res.PikabuId)
	defer unlockUserById(res.PikabuId)

	// complete tasks
	err := task_manager.CompleteTask(
		nil,
		"parse_user_tasks",
		"pikabu_id",
		res.PikabuId,
	)
	if err != nil {
		return err
	}

	err = task_manager.CompleteTask(
		nil,
		"parse_user_tasks",
		"username",
		res.Username,
	)
	if err != nil {
		return err
	}

	var user models.PikabuUser
	err = models.Db.Model(&user).
		Where("pikabu_id = ?", res.PikabuId).
		Select()
	if err != nil && err != pg.ErrNoRows {
		return err
	}

	if err != pg.ErrNoRows {
		user.LastUpdateTimestamp = parsingTimestamp
		user.NextUpdateTimestamp = user.LastUpdateTimestamp + models.TimestampType(config.Settings.UsersMaxUpdatingPeriod)
		user.IsDeleted = true

		err := models.Db.Update(&user)
		if err != nil {
			return err
		}
	}

	var deletedUser models.PikabuDeletedOrNeverExistedUser
	_, err = models.Db.Model(&deletedUser).
		Where("pikabu_id = ?", res.PikabuId).
		SelectOrInsert()
	if err != pg.ErrNoRows && err != nil {
		return err
	}
	if err != pg.ErrNoRows {
		updatingPeriod := deletedUser.NextUpdateTimestamp - deletedUser.LastUpdateTimestamp
		if updatingPeriod < 0 {
			updatingPeriod = -updatingPeriod
		}
		if updatingPeriod == 0 {
			updatingPeriod = models.TimestampType(config.Settings.UsersMaxUpdatingPeriod)
		} else {
			updatingPeriod = models.TimestampType(float32(updatingPeriod) * 1.5)
		}

		deletedUser.NextUpdateTimestamp = parsingTimestamp + updatingPeriod
		deletedUser.LastUpdateTimestamp = parsingTimestamp
		err := models.Db.Update(&deletedUser)
		if err != nil {
			return err
		}
	}

	return nil
}
