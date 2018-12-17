package results_processor

import (
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/task_manager"
	"errors"
	"fmt"
	"github.com/go-pg/pg"
	"gogsweb.2-47.ru/d3dev/pikago"
	"reflect"
	"sync"
	"time"
)

var processUserProfileMutex = &sync.Mutex{}

func processUserProfile(userProfile *pikago.UserProfile) error {
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
	err = saveUserProfile(tx, userProfile)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func saveUserProfile(tx *pg.Tx, userProfile *pikago.UserProfile) error {
	currentTimestamp := models.TimestampType(time.Now().Unix())

	awardIds, err := createAwardIdsArray(tx, userProfile.Awards)
	if err != nil {
		return err
	}
	communityIds, err := createCommunityIdsArray(tx, userProfile.Communities)
	if err != nil {
		return err
	}
	banHistoryIds, err := createBanHistoryIdsArray(tx, userProfile.BanHistory)
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
		AddedTimestamp:      currentTimestamp,
		LastUpdateTimestamp: currentTimestamp,
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

	wasDataChanged := false

	err = processField(tx,
		&user.Username,
		&newUser.Username,
		"pikabu_user_username_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.Gender,
		&newUser.Gender,
		"pikabu_user_gender_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.Rating,
		&newUser.Rating,
		"pikabu_user_rating_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.NumberOfComments,
		&newUser.NumberOfComments,
		"pikabu_user_number_of_comments_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.NumberOfSubscribers,
		&newUser.NumberOfSubscribers,
		"pikabu_user_number_of_subscribers_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.NumberOfStories,
		&newUser.NumberOfStories,
		"pikabu_user_number_of_stories_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.NumberOfHotStories,
		&newUser.NumberOfHotStories,
		"pikabu_user_number_of_hot_stories_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.NumberOfPluses,
		&newUser.NumberOfPluses,
		"pikabu_user_number_of_pluses_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.NumberOfMinuses,
		&newUser.NumberOfMinuses,
		"pikabu_user_number_of_minuses_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.SignupTimestamp,
		&newUser.SignupTimestamp,
		"pikabu_user_signup_timestamp_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.AvatarURL,
		&newUser.AvatarURL,
		"pikabu_user_avatar_url_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.ApprovedText,
		&newUser.ApprovedText,
		"pikabu_user_approved_text_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.Awards,
		&newUser.Awards,
		"pikabu_user_awards_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.Communities,
		&newUser.Communities,
		"pikabu_user_communities_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.BanHistory,
		&newUser.BanHistory,
		"pikabu_user_ban_history_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.BanEndTimestamp,
		&newUser.BanEndTimestamp,
		"pikabu_user_ban_end_timestamp_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.IsRatingHidden,
		&newUser.IsRatingHidden,
		"pikabu_user_is_rating_hidden_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.IsBanned,
		&newUser.IsBanned,
		"pikabu_user_is_banned_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	err = processField(tx,
		&user.IsPermanentlyBanned,
		&newUser.IsPermanentlyBanned,
		"pikabu_user_is_permanently_banned_versions",
		user,
		currentTimestamp,
		&wasDataChanged)
	if err != nil {
		return err
	}

	nextUpdateTimestamp := calculateNextUpdateTimestamp(tx, user, wasDataChanged)
	user.LastUpdateTimestamp = currentTimestamp
	user.NextUpdateTimestamp = nextUpdateTimestamp

	if wasDataChanged {
		print("data was changed!\n")
	}

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
	currentTimestamp models.TimestampType,
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

	if count == 0 {
		_, err := tx.Exec(`
			INSERT INTO `+versionsTableName+`
			(timestamp, item_id, value)
			VALUES (?, ?, ?);
		`, user.AddedTimestamp, user.PikabuId, fieldPtrI)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(`
			INSERT INTO `+versionsTableName+`
			(timestamp, item_id, value)
			VALUES (?, ?, ?)
			ON CONFLICT(item_id, timestamp) DO NOTHING;
	`, user.LastUpdateTimestamp, user.PikabuId, fieldPtrI)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
			INSERT INTO `+versionsTableName+`
			(timestamp, item_id, value)
			VALUES (?, ?, ?)
			ON CONFLICT(item_id, timestamp) DO NOTHING;
	`, currentTimestamp, user.PikabuId, parsedFieldPtrI)
	if err != nil {
		return err
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
	awards []pikago.UserProfileAward,
) ([]uint64, error) {
	// TODO: implement
	return []uint64{}, nil
}

func createCommunityIdsArray(
	tx *pg.Tx,
	awards []pikago.UserProfileCommunity,
) ([]uint64, error) {
	// TODO: implement
	return []uint64{}, nil
}

func createBanHistoryIdsArray(
	tx *pg.Tx,
	awards []pikago.UserProfileBanHistory,
) ([]uint64, error) {
	// TODO: implement
	return []uint64{}, nil
}

func calculateNextUpdateTimestamp(
	tx *pg.Tx,
	user *models.PikabuUser,
	wasDataChanged bool,
) models.TimestampType {

	// TODO: implement!
	return models.TimestampType(time.Now().Unix() + 3600)
	// TODO: move to settings
	/*if user.SignupTimestamp >= models.TimestampType(time.Now().Unix()-3600*24) {

		return models.TimestampType(time.Now().Unix() + 3600*(24+12))
	}

	previousUpdatingPeriod := models.TimestampType(
		math.Abs(float64(user.NextUpdateTimestamp - user.LastUpdateTimestamp)),
	)*/
}
