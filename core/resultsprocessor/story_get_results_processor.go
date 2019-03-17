package resultsprocessor

import (
	"time"

	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg"
	pikago_models "gogsweb.2-47.ru/d3dev/pikago/models"
)

func processStoryGetResults(parsingTimestamp models.TimestampType, storyGetResults []pikago_models.StoryGetResult) error {
	// TODO: make it concurrent
	for _, storyGetResult := range storyGetResults {
		err := processStoryGetResult(parsingTimestamp, &storyGetResult)
		if err != nil {
			return err
		}
	}
	return nil
}

func processStoryGetResult(parsingTimestamp models.TimestampType, storyGetResult *pikago_models.StoryGetResult) error {
	err := processStoryData(parsingTimestamp, storyGetResult.StoryData)
	if err != nil {
		return err
	}

	err = processComments(parsingTimestamp, storyGetResult.Comments)
	if err != nil {
		return err
	}

	return nil
}

var storyLocker = helpers.NewIDLocker()

func processStoryData(parsingTimestamp models.TimestampType, storyData *pikago_models.Story) error {
	storyLocker.Lock(storyData.StoryID.Value)
	defer storyLocker.Unlock(storyData.StoryID.Value)

	err := models.Db.Delete(&models.PikabuDeletedOrNeverExistedStory{
		PikabuID: storyData.StoryID.Value,
	})
	if err != nil {
		return err
	}

	contentBlocks := []models.PikabuStoryBlock{}
	for _, block := range storyData.ContentBlocks {
		contentBlocks = append(contentBlocks, models.PikabuStoryBlock{
			Type: block.Type,
			Data: block.Data,
		})
	}

	newStory := &models.PikabuStory{
		PikabuID:        storyData.StoryID.Value,
		Rating:          int32(storyData.Rating.Value),
		NumberOfPluses:  int32(storyData.NumberOfPluses.Value),
		NumberOfMinuses: int32(storyData.NumberOfMinuses.Value),
		Title:           storyData.Title,
		// TODO: check that RawData is not bytes
		ContentBlocks:      contentBlocks,
		CreatedAtTimestamp: models.TimestampType(storyData.CreatedAtTimestamp.Value),
		StoryURL:           storyData.StoryURL,
		Tags:               storyData.Tags,
		NumberOfComments:   int32(storyData.NumberOfComments.Value),
		IsDeleted:          storyData.IsDeleted,
		IsRatingHidden:     storyData.IsRatingHidden,
		HasMineTag:         storyData.HasMineTag,
		HasAdultTag:        storyData.HasAdultTag,
		IsLongpost:         storyData.IsLongpost,
		AuthorID:           storyData.AuthorID.Value,
		AuthorUsername:     storyData.AuthorUsername,
		AuthorProfileURL:   storyData.AuthorProfileURL,
		AuthorAvatarURL:    storyData.AuthorAvatarURL,
		CommunityLink:      storyData.CommunityLink,
		CommunityName:      storyData.CommunityName,
		CommentsAreHot:     storyData.CommentsAreHot,

		AddedTimestamp:       parsingTimestamp,
		LastUpdateTimestamp:  parsingTimestamp,
		NextUpdateTimestamp:  0,
		TaskTakenAtTimestamp: parsingTimestamp,
	}
	newStory.NextUpdateTimestamp = calculateStoryNextUpdateTimestamp(newStory, false)

	story := &models.PikabuStory{
		PikabuID: storyData.StoryID.Value,
	}
	err = models.Db.Select(story)

	if err == pg.ErrNoRows {
		err := models.Db.Insert(newStory)
		if err != nil {
			return errors.New(err)
		}
		return nil
	} else if err != nil {
		return err
	}

	wasDataChanged, err := processModelFieldsVersions(nil, story, newStory, parsingTimestamp)
	if _, ok := err.(OldParserResultError); ok {
		logger.Log.Warningf("skipping story %v because of old parsing result", storyData.StoryID.Value)
		return nil
	} else if err != nil {
		return errors.New(err)
	}

	nextUpdateTimestamp := calculateStoryNextUpdateTimestamp(story, wasDataChanged)
	story.LastUpdateTimestamp = parsingTimestamp
	story.NextUpdateTimestamp = nextUpdateTimestamp

	err = models.Db.Update(story)
	if err != nil {
		return errors.New(err)
	}

	return nil
}

func calculateStoryNextUpdateTimestamp(
	story *models.PikabuStory,
	wasDataChanged bool,
) models.TimestampType {
	currentTimestamp := models.TimestampType(time.Now().Unix())

	nextUpdateTimestamp := currentTimestamp

	storyTimeGap := currentTimestamp - story.CreatedAtTimestamp
	for gap, updatingPeriod := range map[int64]int64{
		3600:               1800,
		3600 * 3:           3600,
		3600 * 12:          3600 * 6,
		3600 * 24:          3600 * 12,
		3600 * 24 * 7:      3600 * 24,
		3600 * 24 * 30:     3600 * 24 * 7,
		3600 * 24 * 30 * 3: 3600 * 24 * 30,
		3600 * 24 * 30 * 6: 3600 * 24 * 30 * 2,
		4294967296:         3600 * 24 * 30 * 12, // one year
	} {
		if storyTimeGap < models.TimestampType(gap) {
			nextUpdateTimestamp += models.TimestampType(updatingPeriod)
			break
		}
	}

	return nextUpdateTimestamp
}
