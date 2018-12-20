package results_processor

import (
	"bitbucket.org/d3dev/parse_pikabu/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/task_manager"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg"
	"gogsweb.2-47.ru/d3dev/pikago"
	"reflect"
	"sync"
	"time"
)

var processUserProfileMutex = &sync.Mutex{}

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
	awardIds, err := createAwardIdsArray(tx, userProfile.Awards, parsingTimestamp)
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
		Rating:              userProfile.Rating.Value,
		NumberOfComments:    int32(userProfile.CommentsCount.Value),
		NumberOfSubscribers: int32(userProfile.SubscribersCount.Value),
		NumberOfStories:     int32(userProfile.StoriesCount.Value),
		NumberOfHotStories:  int32(userProfile.StoriesHotCount.Value),
		NumberOfPluses:      int32(userProfile.PlusesCount.Value),
		NumberOfMinuses:     int32(userProfile.MinusesCount.Value),
		SignupTimestamp:     models.TimestampType(userProfile.SignupTimestamp.Value),
		AvatarURL:           userProfile.AvatarURL,
		ApprovedText:        userProfile.Approved,
		Awards:              awardIds,
		Communities:         communityIds,
		BanHistory:          banHistoryIds,
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
		return tx.Insert(newUser)
	} else if err != nil {
		return err
	}

	if parsingTimestamp <= user.LastUpdateTimestamp {
		// TODO: find a better way
		logger.Log.Warning("skipping user %v because of old parsing result", user.Username)
		return nil
	}

	wasDataChanged := false

	err = processField(tx,
		&user.Username,
		&newUser.Username,
		"pikabu_user_username_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.Gender,
		&newUser.Gender,
		"pikabu_user_gender_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.Rating,
		&newUser.Rating,
		"pikabu_user_rating_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.NumberOfComments,
		&newUser.NumberOfComments,
		"pikabu_user_number_of_comments_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.NumberOfSubscribers,
		&newUser.NumberOfSubscribers,
		"pikabu_user_number_of_subscribers_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.NumberOfStories,
		&newUser.NumberOfStories,
		"pikabu_user_number_of_stories_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.NumberOfHotStories,
		&newUser.NumberOfHotStories,
		"pikabu_user_number_of_hot_stories_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.NumberOfPluses,
		&newUser.NumberOfPluses,
		"pikabu_user_number_of_pluses_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.NumberOfMinuses,
		&newUser.NumberOfMinuses,
		"pikabu_user_number_of_minuses_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.SignupTimestamp,
		&newUser.SignupTimestamp,
		"pikabu_user_signup_timestamp_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.AvatarURL,
		&newUser.AvatarURL,
		"pikabu_user_avatar_url_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.ApprovedText,
		&newUser.ApprovedText,
		"pikabu_user_approved_text_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.Awards,
		&newUser.Awards,
		"pikabu_user_awards_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.Communities,
		&newUser.Communities,
		"pikabu_user_communities_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.BanHistory,
		&newUser.BanHistory,
		"pikabu_user_ban_history_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.BanEndTimestamp,
		&newUser.BanEndTimestamp,
		"pikabu_user_ban_end_timestamp_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.IsRatingHidden,
		&newUser.IsRatingHidden,
		"pikabu_user_is_rating_hidden_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.IsBanned,
		&newUser.IsBanned,
		"pikabu_user_is_banned_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.IsPermanentlyBanned,
		&newUser.IsPermanentlyBanned,
		"pikabu_user_is_permanently_banned_versions",
		user,
		parsingTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	nextUpdateTimestamp := calculateNextUpdateTimestamp(tx, user, wasDataChanged)
	user.LastUpdateTimestamp = parsingTimestamp
	user.NextUpdateTimestamp = nextUpdateTimestamp

	err = tx.Update(user)
	if err != nil {
		return err
	}

	return nil
}

func processField(
	tx *pg.Tx,
	fieldPtrI interface{},
	parsedFieldPtrI interface{},
	versionsTableName string,
	user *models.PikabuUser,
	parsingTimestamp models.TimestampType,
	wasDataChanged *bool,
) error {
	if reflect.DeepEqual(fieldPtrI, parsedFieldPtrI) {
		return nil
	}
	*wasDataChanged = true

	var count int
	_, err := tx.QueryOne(pg.Scan(&count), `
		SELECT COUNT(*) 
		FROM `+versionsTableName+`
		WHERE item_id = ?
	`, user.PikabuId)
	if err != nil {
		return err
	}

	insertVersion := func(
		timestamp models.TimestampType,
		valuePtr interface{},
		ignoreIfExists bool,
	) error {
		var err error
		// TODO: refactor somehow.
		// Need this shit because go-pg serialize slices as jsonb not as arrays by default
		switch versionsTableName {
		case "pikabu_user_awards_versions", "pikabu_user_communities_versions", "pikabu_user_ban_history_versions":
			var version interface{}

			switch versionsTableName {
			case "pikabu_user_awards_versions":
				version = &models.PikabuUserAwardsVersion{
					FieldVersionBase: models.FieldVersionBase{
						Timestamp: timestamp,
						ItemId:    user.PikabuId,
					},
					Value: *valuePtr.(*[]uint64)}
			case "pikabu_user_communities_versions":
				version = &models.PikabuUserCommunitiesVersion{
					FieldVersionBase: models.FieldVersionBase{
						Timestamp: timestamp,
						ItemId:    user.PikabuId,
					},
					Value: *valuePtr.(*[]uint64)}
			case "pikabu_user_ban_history_versions":
				version = &models.PikabuUserBanHistoryVersion{
					FieldVersionBase: models.FieldVersionBase{
						Timestamp: timestamp,
						ItemId:    user.PikabuId,
					},
					Value: *valuePtr.(*[]uint64)}
			default:
				return errors.New("processField(): bad version table")
			}

			if ignoreIfExists {
				_, err = tx.Model(version).
					OnConflict("DO NOTHING").
					Insert(version)
			} else {
				err = tx.Insert(version)
			}
		default:
			queryPostfix := ""
			if ignoreIfExists {
				queryPostfix = "ON CONFLICT (timestamp, item_id) DO NOTHING"
			}
			_, err = tx.Exec(`
				INSERT INTO `+versionsTableName+`
				(timestamp, item_id, value)
				VALUES (?, ?, ?)
				`+queryPostfix+`;
			`, timestamp, user.PikabuId, valuePtr)
		}
		return err
	}

	if count == 0 {
		err := insertVersion(
			user.AddedTimestamp,
			fieldPtrI,
			false)
		if err != nil {
			return errors.New(err)
		}
	}

	err = insertVersion(
		user.LastUpdateTimestamp,
		fieldPtrI,
		true)
	if err != nil {
		return errors.New(err)
	}

	err = insertVersion(
		parsingTimestamp,
		parsedFieldPtrI,
		false)
	if err != nil {
		return errors.New(err)
	}

	// TODO: refactor this shit
	switch fieldPtr := fieldPtrI.(type) {
	case *string:
		*fieldPtr = *parsedFieldPtrI.(*string)
	case *int:
		*fieldPtr = *parsedFieldPtrI.(*int)
	case *int8:
		*fieldPtr = *parsedFieldPtrI.(*int8)
	case *int16:
		*fieldPtr = *parsedFieldPtrI.(*int16)
	case *int32:
		*fieldPtr = *parsedFieldPtrI.(*int32)
	case *int64:
		*fieldPtr = *parsedFieldPtrI.(*int64)
	case *uint8:
		*fieldPtr = *parsedFieldPtrI.(*uint8)
	case *uint16:
		*fieldPtr = *parsedFieldPtrI.(*uint16)
	case *uint32:
		*fieldPtr = *parsedFieldPtrI.(*uint32)
	case *uint64:
		*fieldPtr = *parsedFieldPtrI.(*uint64)
	case *float32:
		*fieldPtr = *parsedFieldPtrI.(*float32)
	case *float64:
		*fieldPtr = *parsedFieldPtrI.(*float64)
	case *models.TimestampType:
		*fieldPtr = *parsedFieldPtrI.(*models.TimestampType)
	case *bool:
		*fieldPtr = *parsedFieldPtrI.(*bool)
	case *[]uint64:
		*fieldPtr = []uint64{}
		for _, item := range *parsedFieldPtrI.(*[]uint64) {
			*fieldPtr = append(*fieldPtr, item)
		}
	default:
		panic(fmt.Sprintf(
			"processField(): bad type %v",
			reflect.TypeOf(fieldPtr),
		))
	}

	return nil
}

func createAwardIdsArray(
	tx *pg.Tx,
	parsedAwards []pikago.UserProfileAward,
	parsingTimestamp models.TimestampType,
) ([]uint64, error) {
	result := []uint64{}

	for _, parsedAward := range parsedAwards {
		award := &models.PikabuUserAward{
			PikabuId:      parsedAward.Id.Value,
			Timestamp:     parsingTimestamp,
			UserId:        parsedAward.UserId.Value,
			AwardId:       parsedAward.AwardId.Value,
			AwardTitle:    parsedAward.AwardTitle,
			AwardImageURL: parsedAward.AwardImageURL,
			StoryId:       parsedAward.StoryId.Value,
			StoryTitle:    parsedAward.StoryTitle,
			IssuingDate:   parsedAward.IssuingDate,
			IsHidden:      parsedAward.IsHidden.Value != 0,
			CommentId:     parsedAward.CommentId.Value,
			Link:          parsedAward.Link,
		}
		awardFromDb := &models.PikabuUserAward{
			PikabuId: parsedAward.Id.Value,
		}
		err := tx.Select(awardFromDb)
		found := err != pg.ErrNoRows
		if found && err != nil {
			return nil, err
		}
		if found {
			// really bad workaround
			award.Timestamp = awardFromDb.Timestamp
		}
		if found && !reflect.DeepEqual(award, awardFromDb) {
			return nil, errors.New(fmt.Sprintf(
				"award with id %v has been changed. Old state %v, new state %v",
				parsedAward.Id.Value,
				awardFromDb,
				award,
			))
		}
		award.Timestamp = parsingTimestamp
		if !found {
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
			Timestamp: parsingTimestamp,
			Name:      parsedCommunity.Name,
			Link:      parsedCommunity.Link,
			AvatarURL: parsedCommunity.AvatarURL,
		}
		communityFromDb := &models.PikabuUserCommunity{}
		err := tx.Model(communityFromDb).
			Where("link = ?", parsedCommunity.Link).
			Select()
		found := err != pg.ErrNoRows
		if found && err != nil {
			return nil, err
		}

		if found {
			community.Id = communityFromDb.Id
			community.Timestamp = communityFromDb.Timestamp
			if !reflect.DeepEqual(community, communityFromDb) {
				return nil, errors.New(fmt.Sprintf(
					"community with link %v has been changed. Old state %v, new state %v",
					community.Link,
					communityFromDb,
					community,
				))
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
			Timestamp:               parsingTimestamp,
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
		}
		dbBanHistoryItem := &models.PikabuUserBanHistoryItem{
			PikabuId: banHistoryItem.PikabuId,
		}
		err := tx.Select(dbBanHistoryItem)
		found := err != pg.ErrNoRows
		if found && err != nil {
			return nil, err
		}
		if found {
			banHistoryItem.Timestamp = dbBanHistoryItem.Timestamp
		}
		if found && !reflect.DeepEqual(banHistoryItem, dbBanHistoryItem) {
			return nil, errors.New(fmt.Sprintf(
				"ban history item with id %v has been changed. Old state %v, new state %v",
				banHistoryItem.PikabuId,
				dbBanHistoryItem,
				banHistoryItem,
			))
		}
		banHistoryItem.Timestamp = parsingTimestamp
		if !found {
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

	// TODO: implement!
	return models.TimestampType(time.Now().Unix() + 60)
	// TODO: move to settings
	/*if user.SignupTimestamp >= models.TimestampType(time.Now().Unix()-3600*24) {

		return models.TimestampType(time.Now().Unix() + 3600*(24+12))
	}

	previousUpdatingPeriod := models.TimestampType(
		math.Abs(float64(user.NextUpdateTimestamp - user.LastUpdateTimestamp)),
	)*/
}
