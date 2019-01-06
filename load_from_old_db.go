package main

import (
	"fmt"
	"strings"

	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/old_models"
	"bitbucket.org/d3dev/parse_pikabu/results_processor"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/pkg/errors"
	"gogsweb.2-47.ru/d3dev/pikago"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func loadFromOldDb() {
	fmt.Println("loading from db pikabot_graphs...")

	err := models.InitDb()
	if err != nil {
		handleError(err)
	}

	// clear tables
	for _, table := range models.Tables {
		err := models.Db.DropTable(table, &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			handleError(err)
		}
	}

	// create again
	err = models.InitDb()
	if err != nil {
		handleError(err)
	}

	oldDb := pg.Connect(&pg.Options{
		Database: "pikabot_graphs",
		User:     "pikabot_graphs",
		Password: "pikabot_graphs",
	})

	processUsers(oldDb)
}

var skipUsers []old_models.User

func processUsers(oldDb *pg.DB) {
	// select pikabu_id from core_user group by pikabu_id having count(*) > 1;
	err := oldDb.Model(&skipUsers).
		Column("pikabu_id").
		Group("pikabu_id").
		Having("count(*) > 1").
		Select()

	if err != pg.ErrNoRows {
		panicOnError(err)
	}

	offset := 0
	limit := 1024
	for true {
		var users []old_models.User
		err := oldDb.Model(&users).
			Order("pikabu_id").
			Offset(offset).
			Limit(limit).
			Select()

		if err == pg.ErrNoRows {
			break
		}
		panicOnError(err)
		for _, oldUser := range users {
			processUser(oldDb, &oldUser)
		}

		offset += limit
	}
}

func processUser(oldDb *pg.DB, oldUser *old_models.User) {
	for _, skipUser := range skipUsers {
		if skipUser.PikabuId == oldUser.PikabuId {
			fmt.Printf("skipping user with id %v", skipUser.PikabuId)
			return
		}
	}

	awardIds := extractAwardIds(oldUser)

	if oldUser.PikabuId <= 0 {
		panicOnError(errors.Errorf("Bad id %v", oldUser.PikabuId))
	}

	user := &models.PikabuUser{
		PikabuId:            uint64(oldUser.PikabuId),
		Username:            oldUser.Username,
		Gender:              oldUser.Gender,
		Rating:              oldUser.Rating,
		NumberOfComments:    oldUser.CommentsCount,
		NumberOfSubscribers: oldUser.SubscribersCount,
		NumberOfStories:     oldUser.PostsCount,
		NumberOfHotStories:  oldUser.HotPostsCount,
		NumberOfPluses:      oldUser.PlusesCount,
		NumberOfMinuses:     oldUser.MinusesCount,
		SignupTimestamp:     models.TimestampType(oldUser.SignupTimestamp),
		AvatarURL:           oldUser.AvatarURL,
		ApprovedText:        oldUser.Approved,
		AwardIds:            awardIds,
		CommunityIds:        []uint64{},
		BanHistoryItemIds:   []uint64{},
		BanEndTimestamp:     0,
		IsRatingHidden:      oldUser.IsRatingBan,
		IsBanned:            false,
		IsPermanentlyBanned: false,

		// ?
		// IsDeleted bool `sql:",notnull,default:false"`

		AddedTimestamp:      0,
		LastUpdateTimestamp: models.TimestampType(oldUser.LastUpdateTimestamp),
		NextUpdateTimestamp: 0,
	}

	count, err := models.Db.Model(user).
		Where("pikabu_id = ?pikabu_id").
		Count()
	panicOnError(err)
	if count > 0 {
		panicOnError(errors.Errorf("AAA, PANIC!!!!"))
	}

	err = models.Db.Insert(user)
	panicOnError(err)

	processUserVersionsFields(oldDb, oldUser, user)
}

func extractAwardIds(oldUser *old_models.User) []uint64 {
	awardsStr := oldUser.Awards
	awardsStr = strings.TrimSpace(awardsStr)
	if len(awardsStr) == 0 {
		return []uint64{}
	}

	var awards []pikago.UserProfileAward
	err := pikago.JsonUnmarshal([]byte(awardsStr), &awards)
	if err != nil {
		panicOnError(errors.Errorf("unable to unmarshal %v", awardsStr))
	}

	tx, err := models.Db.Begin()
	panicOnError(err)
	defer tx.Rollback()

	results, err := results_processor.CreateAwardIdsArray(
		tx,
		awards,
		models.TimestampType(oldUser.LastUpdateTimestamp),
	)
	panicOnError(err)

	panicOnError(tx.Commit())

	return results
}

func processUserVersionsFields(
	oldDb *pg.DB,
	oldUser *old_models.User,
	user *models.PikabuUser,
) {
	processUserCountersEntryBase(
		oldDb, "core_userratingentry", "pikabu_user_rating_versions", oldUser, user)
	processUserCountersEntryBase(
		oldDb, "core_usersubscriberscountentry", "pikabu_user_number_of_subscribers_versions", oldUser, user)
	processUserCountersEntryBase(
		oldDb, "core_usercommentscountentry", "pikabu_user_number_of_comments_versions", oldUser, user)
	processUserCountersEntryBase(
		oldDb, "core_userpostscountentry", "pikabu_user_number_of_stories_versions", oldUser, user)
	processUserCountersEntryBase(
		oldDb, "core_userhotpostscountentry", "pikabu_user_number_of_hot_stories_versions", oldUser, user)
	processUserCountersEntryBase(
		oldDb, "core_userplusescountentry", "pikabu_user_number_of_pluses_versions", oldUser, user)
	processUserCountersEntryBase(
		oldDb, "core_userminusescountentry", "pikabu_user_number_of_minuses_versions", oldUser, user)
	// TODO: avatar url
}

func processUserCountersEntryBase(
	oldDb *pg.DB,
	tableName string,
	newTableName string,
	oldUser *old_models.User,
	user *models.PikabuUser,
) {
	var result []old_models.CountersEntryBase
	_, err := oldDb.Query(&result, `
		SELECT * FROM `+tableName+` 
		WHERE user_id = ?
		ORDER BY timestamp;
	`, oldUser.Id)
	panicOnError(err)

	for _, item := range result {
		_, err := models.Db.Exec(`
			INSERT INTO `+newTableName+` 
			(timestamp, item_id, value)
			VALUES (?, ?, ?);
		`, models.TimestampType(item.Timestamp), user.PikabuId, item.Value)
		panicOnError(err)
	}
}
