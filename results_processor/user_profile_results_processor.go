package results_processor

import (
	"fmt"
	"sync"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/config"
	"bitbucket.org/d3dev/parse_pikabu/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/task_manager"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg"
	"gogsweb.2-47.ru/d3dev/pikago"
)

var processUserProfileMutex = &sync.Mutex{}

func processUserProfiles(parsingTimestamp models.TimestampType, userProfiles []*pikago.UserProfile) error {
	for _, userProfile := range userProfiles {
		// TODO: make it concurrent
		err := processUserProfile(parsingTimestamp, userProfile)
		if err != nil {
			return err
		}
	}

	return nil
}

func processUserProfile(parsingTimestamp models.TimestampType, userProfile *pikago.UserProfile) error {
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

	// save results
	err = saveUserProfile(tx, parsingTimestamp, userProfile)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func saveUserProfile(tx *pg.Tx, parsingTimestamp models.TimestampType, userProfile *pikago.UserProfile) error {
	awardIds, err := CreateAwardIdsArray(tx, userProfile.Awards, parsingTimestamp)
	if err != nil {
		return err
	}
	communityIds, err := createCommunityIdsArray(tx, userProfile.Communities, parsingTimestamp)
	if err != nil {
		return err
	}
	banHistoryIds, err := createBanHistoryIdsArray(tx, userProfile.BanHistory, parsingTimestamp)
	if err != nil {
		return err
	}
	newUser := &models.PikabuUser{
		PikabuId:            userProfile.UserId.Value,
		Username:            userProfile.Username,
		Gender:              fmt.Sprint(userProfile.Gender.Value),
		Rating:              int32(userProfile.Rating.Value),
		NumberOfComments:    int32(userProfile.CommentsCount.Value),
		NumberOfSubscribers: int32(userProfile.SubscribersCount.Value),
		NumberOfStories:     int32(userProfile.StoriesCount.Value),
		NumberOfHotStories:  int32(userProfile.StoriesHotCount.Value),
		NumberOfPluses:      int32(userProfile.PlusesCount.Value),
		NumberOfMinuses:     int32(userProfile.MinusesCount.Value),
		SignupTimestamp:     models.TimestampType(userProfile.SignupTimestamp.Value),
		AvatarURL:           userProfile.AvatarURL,
		ApprovedText:        userProfile.Approved,
		AwardIds:            awardIds,
		CommunityIds:        communityIds,
		BanHistoryItemIds:   banHistoryIds,
		BanEndTimestamp:     models.TimestampType(userProfile.BanEndTimestamp.Value),
		IsRatingHidden:      userProfile.IsRatingBanned,
		IsBanned:            userProfile.IsUserBanned,
		IsPermanentlyBanned: userProfile.IsUserPermanentlyBanned,
		// IsDeleted: false,
		AddedTimestamp:      parsingTimestamp,
		LastUpdateTimestamp: parsingTimestamp,
		NextUpdateTimestamp: 0,
	}
	newUser.NextUpdateTimestamp = calculateNextUpdateTimestamp(tx, newUser, false)

	user := &models.PikabuUser{
		PikabuId: userProfile.UserId.Value,
	}
	err = tx.Select(user)

	if err == pg.ErrNoRows {
		err := tx.Insert(newUser)
		if err != nil {
			return errors.New(err)
		}
		return nil
	} else if err != nil {
		return err
	}

	wasDataChanged, err := processModelFieldsVersions(tx, user, newUser, parsingTimestamp)
	if _, ok := err.(OldParserResultError); ok {
		logger.Log.Warning("skipping user %v because of old parsing result", user.Username)
		return nil
	} else if err != nil {
		return err
	}

	nextUpdateTimestamp := calculateNextUpdateTimestamp(tx, user, wasDataChanged)
	user.LastUpdateTimestamp = parsingTimestamp
	user.NextUpdateTimestamp = nextUpdateTimestamp

	err = tx.Update(user)
	if err != nil {
		return errors.New(err)
	}

	return nil
}

func CreateAwardIdsArray(
	tx *pg.Tx,
	parsedAwards []pikago.UserProfileAward,
	parsingTimestamp models.TimestampType,
) ([]uint64, error) {
	result := []uint64{}

	for _, parsedAward := range parsedAwards {
		award := &models.PikabuUserAward{
			PikabuId:            parsedAward.Id.Value,
			AddedTimestamp:      parsingTimestamp,
			UserId:              parsedAward.UserId.Value,
			AwardId:             parsedAward.AwardId.Value,
			AwardTitle:          parsedAward.AwardTitle,
			AwardImageURL:       parsedAward.AwardImageURL,
			StoryId:             parsedAward.StoryId.Value,
			StoryTitle:          parsedAward.StoryTitle,
			IssuingDate:         parsedAward.IssuingDate,
			IsHidden:            parsedAward.IsHidden.Value != 0,
			CommentId:           parsedAward.CommentId.Value,
			Link:                parsedAward.Link,
			LastUpdateTimestamp: parsingTimestamp,
		}
		awardFromDb := &models.PikabuUserAward{
			PikabuId: parsedAward.Id.Value,
		}
		err := tx.Select(awardFromDb)
		if err != pg.ErrNoRows && err != nil {
			return nil, err
		}

		found := err != pg.ErrNoRows

		if found {
			_, err := processModelFieldsVersions(tx, awardFromDb, award, parsingTimestamp)
			if _, ok := err.(OldParserResultError); ok {
				logger.Log.Warning("skipping item %v because of old parsing result", award)
			} else {
				if err != nil {
					return nil, err
				}
				awardFromDb.LastUpdateTimestamp = parsingTimestamp
				err = tx.Update(awardFromDb)
				if err != nil {
					return nil, errors.New(err)
				}
			}
		} else {
			err := tx.Insert(award)
			if err != nil {
				return nil, err
			}
		}
		result = append(result, award.PikabuId)
	}
	return result, nil
}

func createCommunityIdsArray(
	tx *pg.Tx,
	parsedCommunities []pikago.UserProfileCommunity,
	parsingTimestamp models.TimestampType,
) ([]uint64, error) {
	result := []uint64{}
	for _, parsedCommunity := range parsedCommunities {
		community := &models.PikabuUserCommunity{
			Name:                parsedCommunity.Name,
			Link:                parsedCommunity.Link,
			AvatarURL:           parsedCommunity.AvatarURL,
			AddedTimestamp:      parsingTimestamp,
			LastUpdateTimestamp: parsingTimestamp,
		}
		communityFromDb := &models.PikabuUserCommunity{}
		err := tx.Model(communityFromDb).
			Where("link = ?", parsedCommunity.Link).
			Select()
		if err != pg.ErrNoRows && err != nil {
			return nil, err
		}

		found := err != pg.ErrNoRows
		if found {
			community.Id = communityFromDb.Id
			community.AddedTimestamp = communityFromDb.AddedTimestamp
			err := tx.Update(community)
			if err != nil {
				return nil, errors.New(err)
			}
		} else {
			_, err := tx.Model(community).Returning("*").Insert()
			if err != nil {
				return nil, err
			}
		}
		result = append(result, community.Id)
	}

	return result, nil
}

func createBanHistoryIdsArray(
	tx *pg.Tx,
	parsedBanHistoryItems []pikago.UserProfileBanHistory,
	parsingTimestamp models.TimestampType,
) ([]uint64, error) {
	result := []uint64{}
	for _, parsedBanHistoryItem := range parsedBanHistoryItems {
		banHistoryItem := &models.PikabuUserBanHistoryItem{
			PikabuId:                parsedBanHistoryItem.Id.Value,
			BanStartTimestamp:       models.TimestampType(parsedBanHistoryItem.BanStartTimestamp.Value),
			CommentId:               parsedBanHistoryItem.CommentId.Value,
			CommentHtmlDeleteReason: parsedBanHistoryItem.CommentHtmlDeleteReason,
			StoryId:                 parsedBanHistoryItem.StoryId.Value,
			UserId:                  parsedBanHistoryItem.UserId.Value,
			BanReason:               parsedBanHistoryItem.BanReason,
			BanReasonId:             parsedBanHistoryItem.BanReasonId.Value,
			StoryURL:                parsedBanHistoryItem.StoryURL,
			ModeratorId:             parsedBanHistoryItem.ModeratorId.Value,
			ModeratorName:           parsedBanHistoryItem.ModeratorName,
			ModeratorAvatar:         parsedBanHistoryItem.ModeratorAvatar,
			ReasonsLimit:            parsedBanHistoryItem.ReasonsLimit.Value,
			ReasonCount:             parsedBanHistoryItem.ReasonCount.Value,
			ReasonTitle:             parsedBanHistoryItem.ReasonTitle,

			AddedTimestamp:      parsingTimestamp,
			LastUpdateTimestamp: parsingTimestamp,
		}
		dbBanHistoryItem := &models.PikabuUserBanHistoryItem{
			PikabuId: banHistoryItem.PikabuId,
		}
		err := tx.Select(dbBanHistoryItem)
		if err != pg.ErrNoRows && err != nil {
			return nil, err
		}

		found := err != pg.ErrNoRows
		if found {
			_, err := processModelFieldsVersions(tx, dbBanHistoryItem, banHistoryItem, parsingTimestamp)
			if _, ok := err.(OldParserResultError); ok {
				logger.Log.Warning("skipping item %v because of old parsing result", banHistoryItem)
			} else {
				if err != nil {
					return nil, err
				}
				dbBanHistoryItem.LastUpdateTimestamp = parsingTimestamp
				err = tx.Update(dbBanHistoryItem)
				if err != nil {
					return nil, errors.New(err)
				}
			}
		} else {
			err := tx.Insert(banHistoryItem)
			if err != nil {
				return nil, err
			}
		}
		result = append(result, banHistoryItem.PikabuId)
	}

	return result, nil
}

func calculateNextUpdateTimestamp(
	tx *pg.Tx,
	user *models.PikabuUser,
	wasDataChanged bool,
) models.TimestampType {
	currentTimestamp := models.TimestampType(time.Now().Unix())
	updatingPeriod := user.NextUpdateTimestamp - user.LastUpdateTimestamp
	if updatingPeriod < 0 {
		updatingPeriod = 0
	}
	// update new users frequently
	if user.SignupTimestamp >=
		models.TimestampType(time.Now().Unix()-int64(config.Settings.NewUserTime)) {

		return currentTimestamp + models.TimestampType(config.Settings.NewUsersUpdatingPeriod)
	}

	if wasDataChanged {
		updatingPeriod /= 2
	} else {
		updatingPeriod += models.TimestampType(config.Settings.UsersUpdatingPeriodIncreasingValue)
	}

	if updatingPeriod < models.TimestampType(config.Settings.UsersMinUpdatingPeriod) {
		updatingPeriod = models.TimestampType(config.Settings.UsersMinUpdatingPeriod)
	}
	if updatingPeriod > models.TimestampType(config.Settings.UsersMaxUpdatingPeriod) {
		updatingPeriod = models.TimestampType(config.Settings.UsersMaxUpdatingPeriod)
	}

	return currentTimestamp + models.TimestampType(updatingPeriod)
}
