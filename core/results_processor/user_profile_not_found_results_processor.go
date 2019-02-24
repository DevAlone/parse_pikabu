package results_processor

import (
	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/core/task_manager"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-pg/pg"
)

func processUserProfileNotFoundResults(res *models.ParserUserProfileNotFoundResult) error {
	parsingTimestamp := res.ParsingTimestamp
	for _, result := range res.Results {
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

	tx, err := models.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// complete tasks
	err = task_manager.CompleteTask(
		tx,
		"parse_user_tasks",
		"pikabu_id",
		res.PikabuId,
	)
	if err != nil {
		return err
	}

	err = task_manager.CompleteTask(
		tx,
		"parse_user_tasks",
		"username",
		res.Username,
	)
	if err != nil {
		return err
	}

	var user models.PikabuUser
	err = tx.Model(&user).
		Where("pikabu_id = ?", res.PikabuId).Select()
	if err != nil && err != pg.ErrNoRows {
		return err
	}

	if err != pg.ErrNoRows {
		user.NextUpdateTimestamp = user.LastUpdateTimestamp + models.TimestampType(config.Settings.UsersMaxUpdatingPeriod)
		user.IsDeleted = true

		err := tx.Update(&user)
		if err != nil {
			return err
		}
	}

	var deletedUser models.PikabuDeletedOrNeverExistedUser
	_, err = tx.Model(&deletedUser).
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
		err := tx.Update(&deletedUser)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
