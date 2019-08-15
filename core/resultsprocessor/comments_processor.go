package resultsprocessor

import (
	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/modelhooks"
	"bitbucket.org/d3dev/parse_pikabu/models"
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
		modelhooks.HandlePikabuCommentCreate(*newComment, parsingTimestamp)
		err := models.Db.Insert(newComment)
		if err != nil {
			return errors.New(err)
		}
		return nil
	} else if err != nil {
		return errors.New(err)
	}

	modelhooks.HandlePikabuCommentChange(*comment, *newComment, parsingTimestamp)

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
	return models.TimestampType(86400)
}
