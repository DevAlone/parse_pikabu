package resultsprocessor

import (
	"time"

	"github.com/DevAlone/parse_pikabu/core/logger"
	"github.com/DevAlone/parse_pikabu/helpers"
	"github.com/DevAlone/parse_pikabu/modelhooks"
	"github.com/DevAlone/parse_pikabu/models"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg"
	pikago_models "gogsweb.2-47.ru/d3dev/pikago/models"
)

func processComments(parsingTimestamp models.TimestampType, comments []pikago_models.Comment) error {
	// TODO: make it concurrent
	for _, comment := range comments {
		err := processComment(parsingTimestamp, &comment)
		if err != nil {
			return err
		}
	}
	return nil
}

var commentLocker = helpers.NewIDLocker()

func processComment(parsingTimestamp models.TimestampType, commentData *pikago_models.Comment) error {
	commentLocker.Lock(commentData.ID.Value)
	defer commentLocker.Unlock(commentData.ID.Value)

	err := models.Db.Delete(&models.PikabuDeletedOrNeverExistedComment{
		PikabuID: commentData.ID.Value,
	})
	if err != pg.ErrNoRows && err != nil {
		return errors.New(err)
	}

	images := []models.PikabuCommentImage{}
	for _, image := range commentData.Content.Images {
		animationFormats := map[string]uint64{}
		if image.Animation != nil {
			for key, value := range image.Animation.Formats {
				animationFormats[key] = value.Value
			}
		}
		size := []uint64{}
		for _, value := range image.Size {
			size = append(size, value.Value)
		}
		animationPreviewURL := ""
		animationBaseURL := ""
		if image.Animation != nil {
			animationPreviewURL = image.Animation.PreviewURL
			animationBaseURL = image.Animation.BaseURL
		}

		images = append(images, models.PikabuCommentImage{
			SmallURL:            image.SmallURL,
			LargeURL:            image.LargeURL,
			AnimationPreviewURL: animationPreviewURL,
			AnimationBaseURL:    animationBaseURL,
			AnimationFormats:    animationFormats,
			Size:                size,
		})
	}

	newComment := &models.PikabuComment{
		PikabuID: commentData.ID.Value,

		ParentID:                   commentData.ParentID.Value,
		CreatedAtTimestamp:         models.TimestampType(commentData.CreatedAtTimestamp.Value),
		Text:                       commentData.Content.Text,
		Images:                     images,
		Rating:                     int32(commentData.Rating.Value),
		NumberOfPluses:             int32(commentData.NumberOfPluses.Value),
		NumberOfMinuses:            int32(commentData.NumberOfMinuses.Value),
		StoryID:                    commentData.StoryID.Value,
		StoryURL:                   commentData.StoryURL,
		StoryTitle:                 commentData.StoryTitle,
		AuthorID:                   commentData.AuthorID.Value,
		AuthorUsername:             commentData.AuthorUsername,
		AuthorGender:               int32(commentData.AuthorGender.Value),
		AuthorAvatarURL:            commentData.AuthorAvatarURL,
		IgnoreCode:                 int32(commentData.IgnoreCode.Value),
		IsIgnoredBySomeone:         commentData.IsIgnoredBySomeone,
		IgnoredBy:                  commentData.IgnoredBy,
		IsAuthorProfileDeleted:     commentData.IsAuthorProfileDeleted,
		IsDeleted:                  commentData.IsDeleted,
		IsAuthorCommunityModerator: commentData.IsAuthorCommunityModerator,
		IsAuthorPikabuTeam:         commentData.IsAuthorPikabuTeam,
		IsAuthorOfficial:           commentData.IsAuthorOfficial,
		IsRatingHidden:             commentData.Rating.IsNull,

		AddedTimestamp:       parsingTimestamp,
		LastUpdateTimestamp:  parsingTimestamp,
		NextUpdateTimestamp:  0,
		TaskTakenAtTimestamp: parsingTimestamp,
	}

	newComment.NextUpdateTimestamp = calculateCommentNextUpdateTimestamp(newComment, false)

	comment := &models.PikabuComment{
		PikabuID: commentData.ID.Value,
	}

	err = models.Db.Select(comment)

	if err == pg.ErrNoRows {
		modelhooks.HandleModelCreated(*newComment, parsingTimestamp)
		err := models.Db.Insert(newComment)
		if err != nil {
			return errors.New(err)
		}
		return nil
	} else if err != nil {
		return errors.New(err)
	}

	modelhooks.HandleModelChanged(*comment, *newComment, parsingTimestamp)

	wasDataChanged, err := processModelFieldsVersions(nil, comment, newComment, parsingTimestamp)
	if _, ok := err.(OldParserResultError); ok {
		logger.Log.Warningf("skipping comment %v because of old parsing result", comment)
		return nil
	} else if err != nil {
		return errors.New(err)
	}

	nextUpdateTimestamp := calculateCommentNextUpdateTimestamp(comment, wasDataChanged)
	comment.LastUpdateTimestamp = parsingTimestamp
	comment.NextUpdateTimestamp = nextUpdateTimestamp

	err = models.Db.Update(comment)
	if err != nil {
		return errors.New(err)
	}

	return nil
}

func calculateCommentNextUpdateTimestamp(
	comment *models.PikabuComment,
	wasDataChanged bool,
) models.TimestampType {
	currentTimestamp := models.TimestampType(time.Now().Unix())

	nextUpdateTimestamp := currentTimestamp

	commentAgeInSeconds := currentTimestamp - comment.CreatedAtTimestamp
	for age, updatingPeriod := range map[int64]int64{
		3600 * 24 * 30 * 3: 3600 * 24 * 7,
		3600 * 24 * 30 * 6: 3600 * 24 * 30 * 2,
		4294967296:         3600 * 24 * 30 * 12, // 1 year
	} {
		if commentAgeInSeconds < models.TimestampType(age) {
			nextUpdateTimestamp += models.TimestampType(updatingPeriod)
			break
		}
	}

	return nextUpdateTimestamp
}
