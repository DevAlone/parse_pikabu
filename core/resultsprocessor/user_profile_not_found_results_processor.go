package resultsprocessor

import (
	"github.com/DevAlone/parse_pikabu/core/config"
	"github.com/DevAlone/parse_pikabu/models"
	"github.com/go-errors/errors"
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
	lockUserByID(res.PikabuID)
	defer unlockUserByID(res.PikabuID)

	var user models.PikabuUser
	err := models.Db.Model(&user).
		Where("pikabu_id = ?", res.PikabuID).
		Select()
	if err != nil && err != pg.ErrNoRows {
		return errors.New(err)
	}

	if err != pg.ErrNoRows {
		user.LastUpdateTimestamp = parsingTimestamp
		user.NextUpdateTimestamp = parsingTimestamp + models.TimestampType(config.Settings.UsersMaxUpdatingPeriod)
		user.IsDeleted = true

		err := models.Db.Update(&user)
		if err != nil {
			return errors.New(err)
		}
	}

	// deletedUser

	var deletedUser models.PikabuDeletedOrNeverExistedUser
	_, err = models.Db.Model(&deletedUser).
		Where("pikabu_id = ?", res.PikabuID).
		SelectOrInsert()
	if err != pg.ErrNoRows && err != nil {
		return errors.New(err)
	}
	if err != pg.ErrNoRows {
		updatingPeriod := deletedUser.NextUpdateTimestamp - deletedUser.LastUpdateTimestamp
		if updatingPeriod < 0 {
			updatingPeriod = 0
		}
		if updatingPeriod == 0 {
			updatingPeriod = models.TimestampType(config.Settings.UsersMinUpdatingPeriod)
		} else {
			updatingPeriod = models.TimestampType(float32(updatingPeriod) * 1.5)
		}

		deletedUser.NextUpdateTimestamp = parsingTimestamp + updatingPeriod
		deletedUser.LastUpdateTimestamp = parsingTimestamp
		err := models.Db.Update(&deletedUser)
		if err != nil {
			return errors.New(err)
		}
	}

	return nil
}
