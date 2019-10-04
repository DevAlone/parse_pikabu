package resultsprocessor

import (
	"github.com/DevAlone/parse_pikabu/core/config"
	"github.com/DevAlone/parse_pikabu/models"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg"
)

func processStoryNotFoundResults(parsingTimestamp models.TimestampType, res []models.ParserStoryNotFoundResultData) error {
	for _, result := range res {
		err := processStoryNotFoundResult(parsingTimestamp, &result)
		if err != nil {
			// TODO: log
			return err
		}
	}

	return nil
}

func processStoryNotFoundResult(
	parsingTimestamp models.TimestampType,
	res *models.ParserStoryNotFoundResultData,
) error {
	storyLocker.Lock(res.PikabuID)
	defer storyLocker.Unlock(res.PikabuID)

	story := models.PikabuStory{
		PikabuID: res.PikabuID,
	}
	err := models.Db.Select(&story)
	if err != pg.ErrNoRows && err != nil {
		return errors.New(err)
	}
	if err != pg.ErrNoRows {
		story.LastUpdateTimestamp = parsingTimestamp
		story.NextUpdateTimestamp = calculateStoryNextUpdateTimestamp(&story, false)
		story.IsPermanentlyDeleted = true

		err := models.Db.Update(&story)
		if err != nil {
			return errors.New(err)
		}
	}

	// deletedOrNeverExistedStory

	deletedOrNeverExistedStory := models.PikabuDeletedOrNeverExistedStory{
		PikabuID: res.PikabuID,
	}
	err = models.Db.Select(&deletedOrNeverExistedStory)
	if err != pg.ErrNoRows && err != nil {
		return errors.New(err)
	}
	if err != pg.ErrNoRows {
		updatingPeriod := deletedOrNeverExistedStory.NextUpdateTimestamp - deletedOrNeverExistedStory.LastUpdateTimestamp
		if updatingPeriod < 0 {
			updatingPeriod = 0
		}
		if updatingPeriod == 0 {
			updatingPeriod = models.TimestampType(config.Settings.UsersMinUpdatingPeriod)
		} else {
			updatingPeriod = models.TimestampType(float32(updatingPeriod) * 1.5)
		}

		deletedOrNeverExistedStory.NextUpdateTimestamp = parsingTimestamp + updatingPeriod
		deletedOrNeverExistedStory.LastUpdateTimestamp = parsingTimestamp
		err := models.Db.Update(&deletedOrNeverExistedStory)
		if err != nil {
			return errors.New(err)
		}
	}

	return nil
}
