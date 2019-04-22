package resultsprocessor

import (
	"fmt"
	"sync"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/core/taskmanager"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg"
	pikago_models "gogsweb.2-47.ru/d3dev/pikago/models"
)

func processUserProfiles(parsingTimestamp models.TimestampType, userProfiles []*pikago_models.UserProfile) error {
	for _, userProfile := range userProfiles {
		// TODO: make it concurrent
		err := processUserProfile(parsingTimestamp, userProfile)
		if err != nil {
			return err
		}
	}

	return nil
}

var userProfileIdLocks = map[uint64]bool{}
var userProfileIdLocksMutex = sync.Mutex{}

func lockUserById(userId uint64) {
	// lock results with the same id
	found := true
	for found {
		userProfileIdLocksMutex.Lock()
		_, found = userProfileIdLocks[userId]
		if !found {
			userProfileIdLocks[userId] = true
			userProfileIdLocksMutex.Unlock()
			return
		}
		userProfileIdLocksMutex.Unlock()

		time.Sleep(10 * time.Millisecond)
	}
}

func unlockUserById(userId uint64) {
	userProfileIdLocksMutex.Lock()
	delete(userProfileIdLocks, userId)
	userProfileIdLocksMutex.Unlock()
}

func processUserProfile(parsingTimestamp models.TimestampType, userProfile *pikago_models.UserProfile) error {
	lockUserById(userProfile.UserId.Value)
	defer unlockUserById(userProfile.UserId.Value)

	if userProfile == nil {
		return errors.New("nil user profile")
	}
	/*
		tx, err := models.Db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()
	*/

	// complete tasks
	err := taskmanager.CompleteTask(
		nil,
		"parse_user_tasks",
		"pikabu_id",
		userProfile.UserId.Value,
	)
	if err != nil {
		return err
	}

	task := models.ParseUserTask{
		PikabuID: userProfile.UserId.Value,
	}
	err = models.Db.Select(&task)
	if err != pg.ErrNoRows && err != nil {
		return err
	}
	if err != pg.ErrNoRows {
		task.Username = userProfile.Username
		err = models.Db.Update(&task)
		if err != nil {
			return err
		}
	}

	err = taskmanager.CompleteTask(
		nil,
		"parse_user_tasks",
		"username",
		userProfile.Username,
	)
	if err != nil {
		return err
	}

	_, err = models.Db.Model(&models.PikabuDeletedOrNeverExistedUser{
		PikabuID: userProfile.UserId.Value,
	}).WherePK().Delete()
	if err != nil && err != pg.ErrNoRows {
		return err
	}

	// save results
	err = saveUserProfile(parsingTimestamp, userProfile)
	if err != nil {
		return err
	}

	return nil
}

func saveUserProfile(parsingTimestamp models.TimestampType, userProfile *pikago_models.UserProfile) error {
	awardIds, err := CreateAwardIdsArray(userProfile.Awards, parsingTimestamp)
	if err != nil {
		return err
	}
	communityIds, err := createCommunityIdsArray(userProfile.Communities, parsingTimestamp)
	if err != nil {
		return err
	}
	banHistoryIds, err := createBanHistoryIdsArray(userProfile.BanHistory, parsingTimestamp)
	if err != nil {
		return err
	}
	newUser := &models.PikabuUser{
		PikabuID:            userProfile.UserId.Value,
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
		IsDeleted:           false,
		AddedTimestamp:      parsingTimestamp,
		LastUpdateTimestamp: parsingTimestamp,
		NextUpdateTimestamp: 0,
	}
	newUser.NextUpdateTimestamp = calculateNextUpdateTimestamp(newUser, false)

	user := &models.PikabuUser{
		PikabuID: userProfile.UserId.Value,
	}
	err = models.Db.Select(user)

	if err == pg.ErrNoRows {
		err := models.Db.Insert(newUser)
		if err != nil {
			return errors.New(err)
		}
		return nil
	} else if err != nil {
		return err
	}

	wasDataChanged, err := processModelFieldsVersions(nil, user, newUser, parsingTimestamp)
	if _, ok := err.(OldParserResultError); ok {
		logger.Log.Warningf("skipping user %v because of old parsing result", user.Username)
		return nil
	} else if err != nil {
		return errors.New(err)
	}

	nextUpdateTimestamp := calculateNextUpdateTimestamp(user, wasDataChanged)
	user.LastUpdateTimestamp = parsingTimestamp
	user.NextUpdateTimestamp = nextUpdateTimestamp

	err = models.Db.Update(user)
	if err != nil {
		return errors.New(err)
	}

	return nil
}

func CreateAwardIdsArray(
	parsedAwards []pikago_models.UserProfileAward,
	parsingTimestamp models.TimestampType,
) ([]uint64, error) {
	result := []uint64{}

	for _, parsedAward := range parsedAwards {
		award := &models.PikabuUserAward{
			PikabuID:            parsedAward.Id.Value,
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
			PikabuID: parsedAward.Id.Value,
		}
		err := models.Db.Select(awardFromDb)
		if err != pg.ErrNoRows && err != nil {
			return nil, err
		}

		found := err != pg.ErrNoRows

		if found {
			_, err := processModelFieldsVersions(nil, awardFromDb, award, parsingTimestamp)
			if _, ok := err.(OldParserResultError); ok {
				logger.Log.Warningf("skipping item %v because of old parsing result", award)
			} else {
				if err != nil {
					return nil, err
				}
				awardFromDb.LastUpdateTimestamp = parsingTimestamp
				err = models.Db.Update(awardFromDb)
				if err != nil {
					return nil, errors.New(err)
				}
			}
		} else {
			err := models.Db.Insert(award)
			if err != nil {
				return nil, err
			}
		}
		result = append(result, award.PikabuID)
	}
	return result, nil
}

func createCommunityIdsArray(
	parsedCommunities []pikago_models.UserProfileCommunity,
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
		err := models.Db.Model(communityFromDb).
			Where("link = ?", parsedCommunity.Link).
			Select()
		if err != pg.ErrNoRows && err != nil {
			return nil, err
		}

		found := err != pg.ErrNoRows
		if found {
			community.Id = communityFromDb.Id
			community.AddedTimestamp = communityFromDb.AddedTimestamp
			err := models.Db.Update(community)
			if err != nil {
				return nil, errors.New(err)
			}
		} else {
			_, err := models.Db.Model(community).
				Returning("*").
				Insert()
			if err != nil {
				return nil, err
			}
		}
		result = append(result, community.Id)
	}

	return result, nil
}

func createBanHistoryIdsArray(
	parsedBanHistoryItems []pikago_models.UserProfileBanHistory,
	parsingTimestamp models.TimestampType,
) ([]uint64, error) {
	result := []uint64{}
	for _, parsedBanHistoryItem := range parsedBanHistoryItems {
		banHistoryItem := &models.PikabuUserBanHistoryItem{
			PikabuID:                parsedBanHistoryItem.Id.Value,
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
			PikabuID: banHistoryItem.PikabuID,
		}
		err := models.Db.Select(dbBanHistoryItem)
		if err != pg.ErrNoRows && err != nil {
			return nil, err
		}

		found := err != pg.ErrNoRows
		if found {
			_, err := processModelFieldsVersions(nil, dbBanHistoryItem, banHistoryItem, parsingTimestamp)
			if _, ok := err.(OldParserResultError); ok {
				logger.Log.Warningf("skipping item %v because of old parsing result", banHistoryItem)
			} else {
				if err != nil {
					return nil, err
				}
				dbBanHistoryItem.LastUpdateTimestamp = parsingTimestamp
				err = models.Db.Update(dbBanHistoryItem)
				if err != nil {
					return nil, errors.New(err)
				}
			}
		} else {
			err := models.Db.Insert(banHistoryItem)
			if err != nil {
				return nil, err
			}
		}
		result = append(result, banHistoryItem.PikabuID)
	}

	return result, nil
}

func calculateNextUpdateTimestamp(
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
