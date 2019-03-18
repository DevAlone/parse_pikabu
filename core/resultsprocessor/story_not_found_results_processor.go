package resultsprocessor

import (
	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/models"
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

	deletedOrNeverExistedStory := models.PikabuDeletedOrNeverExistedStory{
		PikabuID: res.PikabuID,
	}
	err := models.Db.Select(&deletedOrNeverExistedStory)
	if err != pg.ErrNoRows && err != nil {
		return err
	}
	if err != pg.ErrNoRows {
		updatingPeriod := deletedOrNeverExistedStory.NextUpdateTimestamp - deletedOrNeverExistedStory.LastUpdateTimestamp
		if updatingPeriod < 0 {
			updatingPeriod = -updatingPeriod
		}
		if updatingPeriod == 0 {
			updatingPeriod = models.TimestampType(config.Settings.UsersMinUpdatingPeriod)
		} else {
			updatingPeriod = models.TimestampType(float32(updatingPeriod) * 1.5)
		}

		deletedOrNeverExistedStory.NextUpdateTimestamp = parsingTimestamp + updatingPeriod
		deletedOrNeverExistedStory.LastUpdateTimestamp = parsingTimestamp
	}

	story := models.PikabuStory{
		PikabuID: res.PikabuID,
	}
	err = models.Db.Select(&story)
	if err != pg.ErrNoRows && err != nil {
		return err
	}
	if err != pg.ErrNoRows {
		story.IsPermanentlyDeleted = true
		story.LastUpdateTimestamp = parsingTimestamp
		story.NextUpdateTimestamp = calculateStoryNextUpdateTimestamp(&story, true)
		story.TaskTakenAtTimestamp = parsingTimestamp

		err := models.Db.Update(&story)
		if err != nil {
			return err
		}
	}

	return nil
}
