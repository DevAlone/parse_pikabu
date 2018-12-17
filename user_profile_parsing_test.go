package parse_pikabu

import (
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/results_processor"
	"github.com/go-pg/pg/orm"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestUserProfileParsing(t *testing.T) {
	err := models.InitDb()
	if err != nil {
		panic(err)
	}

	// clear tables
	for _, table := range models.Tables {
		err := models.Db.DropTable(table, &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			panic(err)
		}
	}

	// create again
	err = models.InitDb()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// start server
	/*go func() {
		err := server.Run()
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	// start task manager
	go func() {
		err := task_manager.Run()
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()*/

	// start results processor
	go func() {
		err := results_processor.Run()
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	err = pushTaskToQueue([]byte(`
{
	"parsing_timestamp": 100,
	"data": {
		"user": {
			"current_user_id": 0,
			"user_id": "2561615",
			"user_name": "Pisacavtor",
			"rating": "-3.5",
			"gender": "6",
			"comments_count": 3,
			"stories_count": 2,
			"stories_hot_count": "1",
			"pluses_count": 5,
			"minuses_count": 9,
			"signup_date": "1544846469",
			"is_rating_ban": true,
			"avatar": "https://cs8.pikabu.ru/avatars/2561/x2561615-512432259.png",
			"awards": [],
			"is_subscribed": false,
			"is_ignored": false,
			"note": null,
			"approved": "approved User",
			"communities": [],
			"subscribers_count": 1001,
			"is_user_banned": true,
			"is_user_fully_banned": false,
			"public_ban_history": [
			  {
				"id": "151513",
				"date": 1544854692,
				"moderator_id": "1836690",
				"comment_id": "0",
				"comment_desc": "",
				"story_id": "6354471",
				"user_id": "2561615",
				"reason": "Отсутствие пруфа или неподтверждённая/искажённая информация (вброс)",
				"reason_id": "94",
				"story_url": "https://pikabu.ru/story/3_chasa_pyitok_6354471",
				"moderator_name": "depotato",
				"moderator_avatar": "https://cs5.pikabu.ru/avatars/1836/s1836690-1399622318.png",
				"reason_limit": null,
				"reason_count": null,
				"reason_title": null
			  }
			],
			"user_ban_time": 1545459492
		}
	}
}
`,
	))
	if err != nil {
		panic(err)
	}

	// TODO: wait for queue to become empty
	time.Sleep(1 * time.Second)

	user := &models.PikabuUser{
		PikabuId: 2561615,
	}
	err = models.Db.Select(user)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "Pisacavtor", user.Username)
	assert.Equal(t, float32(-3.5), user.Rating)
	assert.Equal(t, "6", user.Gender)
	assert.Equal(t, int32(3), user.NumberOfComments)
	assert.Equal(t, int32(2), user.NumberOfStories)
	assert.Equal(t, int32(1), user.NumberOfHotStories)
	assert.Equal(t, int32(5), user.NumberOfPluses)
	assert.Equal(t, int32(9), user.NumberOfMinuses)
	assert.Equal(t, models.TimestampType(1544846469), user.SignupTimestamp)
	assert.Equal(t, true, user.IsRatingHidden)
	assert.Equal(t, "https://cs8.pikabu.ru/avatars/2561/x2561615-512432259.png", user.AvatarURL)
	// TODO: check awards
	// assert.Equal(t, true, user.Awards)
	assert.Equal(t, "approved User", user.ApprovedText)
	// TODO: check communities
	assert.Equal(t, int32(1001), user.NumberOfSubscribers)
	assert.Equal(t, true, user.IsBanned)
	assert.Equal(t, false, user.IsPermanentlyBanned)
	// TODO: check public ban history
	assert.Equal(t, models.TimestampType(1545459492), user.BanEndTimestamp)
	assert.Equal(t, models.TimestampType(100), user.AddedTimestamp)
	assert.Equal(t, models.TimestampType(100), user.LastUpdateTimestamp)

	// change some fields
	err = pushTaskToQueue([]byte(`
{
	"parsing_timestamp": 201,
	"data": {
		"user": {
			"current_user_id": 0,
			"user_id": "2561615",
			"user_name": "Pisacavtor1",
			"rating": "10.5",
			"gender": "6",
			"comments_count": 9,
			"stories_count": 2,
			"stories_hot_count": "1",
			"pluses_count": 5,
			"minuses_count": 9,
			"signup_date": "1544846469",
			"is_rating_ban": false,
			"avatar": "https://cs8.pikabu.ru/avatars/2561/x2561615-512432259.png",
			"awards": [],
			"is_subscribed": false,
			"is_ignored": false,
			"note": null,
			"approved": "approved User",
			"communities": [],
			"subscribers_count": 1001,
			"is_user_banned": true,
			"is_user_fully_banned": false,
			"public_ban_history": [
			  {
				"id": "151513",
				"date": 1544854692,
				"moderator_id": "1836690",
				"comment_id": "0",
				"comment_desc": "",
				"story_id": "6354471",
				"user_id": "2561615",
				"reason": "Отсутствие пруфа или неподтверждённая/искажённая информация (вброс)",
				"reason_id": "94",
				"story_url": "https://pikabu.ru/story/3_chasa_pyitok_6354471",
				"moderator_name": "depotato",
				"moderator_avatar": "https://cs5.pikabu.ru/avatars/1836/s1836690-1399622318.png",
				"reason_limit": null,
				"reason_count": null,
				"reason_title": null
			  }
			],
			"user_ban_time": 1545459492
		}
	}
}
`,
	))
	if err != nil {
		panic(err)
	}

	// TODO: wait for queue to become empty
	time.Sleep(1 * time.Second)

	user = &models.PikabuUser{
		PikabuId: 2561615,
	}
	err = models.Db.Select(user)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "Pisacavtor1", user.Username)
	assert.Equal(t, float32(10.5), user.Rating)
	assert.Equal(t, "6", user.Gender)
	assert.Equal(t, int32(9), user.NumberOfComments)
	assert.Equal(t, int32(2), user.NumberOfStories)
	assert.Equal(t, int32(1), user.NumberOfHotStories)
	assert.Equal(t, int32(5), user.NumberOfPluses)
	assert.Equal(t, int32(9), user.NumberOfMinuses)
	assert.Equal(t, models.TimestampType(1544846469), user.SignupTimestamp)
	assert.Equal(t, false, user.IsRatingHidden)
	assert.Equal(t, "https://cs8.pikabu.ru/avatars/2561/x2561615-512432259.png", user.AvatarURL)
	// TODO: check awards
	// assert.Equal(t, true, user.Awards)
	assert.Equal(t, "approved User", user.ApprovedText)
	// TODO: check communities
	assert.Equal(t, int32(1001), user.NumberOfSubscribers)
	assert.Equal(t, true, user.IsBanned)
	assert.Equal(t, false, user.IsPermanentlyBanned)
	// TODO: check public ban history
	assert.Equal(t, models.TimestampType(1545459492), user.BanEndTimestamp)
	assert.Equal(t, models.TimestampType(100), user.AddedTimestamp)
	assert.Equal(t, models.TimestampType(201), user.LastUpdateTimestamp)

	err = pushTaskToQueue([]byte(`
{
	"parsing_timestamp": 555,
	"data": {
		"user": {
			"current_user_id": 0,
			"user_id": "2561615",
			"user_name": "Pisacavtor",
			"rating": "5",
			"gender": "6",
			"comments_count": 9,
			"stories_count": 2,
			"stories_hot_count": "1",
			"pluses_count": 5,
			"minuses_count": 9,
			"signup_date": "1544846469",
			"is_rating_ban": false,
			"avatar": "https://cs8.pikabu.ru/avatars/2561/x2561615-512432259.png",
			"awards": [],
			"is_subscribed": false,
			"is_ignored": false,
			"note": null,
			"approved": "approved User",
			"communities": [],
			"subscribers_count": 1001,
			"is_user_banned": true,
			"is_user_fully_banned": false,
			"public_ban_history": [
			  {
				"id": "151513",
				"date": 1544854692,
				"moderator_id": "1836690",
				"comment_id": "0",
				"comment_desc": "",
				"story_id": "6354471",
				"user_id": "2561615",
				"reason": "Отсутствие пруфа или неподтверждённая/искажённая информация (вброс)",
				"reason_id": "94",
				"story_url": "https://pikabu.ru/story/3_chasa_pyitok_6354471",
				"moderator_name": "depotato",
				"moderator_avatar": "https://cs5.pikabu.ru/avatars/1836/s1836690-1399622318.png",
				"reason_limit": null,
				"reason_count": null,
				"reason_title": null
			  }
			],
			"user_ban_time": 100
		}
	}
}
`,
	))
	if err != nil {
		panic(err)
	}

	// TODO: wait for queue to become empty
	time.Sleep(1 * time.Second)

	user = &models.PikabuUser{
		PikabuId: 2561615,
	}
	err = models.Db.Select(user)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "Pisacavtor", user.Username)
	assert.Equal(t, float32(5), user.Rating)
	assert.Equal(t, "6", user.Gender)
	assert.Equal(t, int32(9), user.NumberOfComments)
	assert.Equal(t, int32(2), user.NumberOfStories)
	assert.Equal(t, int32(1), user.NumberOfHotStories)
	assert.Equal(t, int32(5), user.NumberOfPluses)
	assert.Equal(t, int32(9), user.NumberOfMinuses)
	assert.Equal(t, models.TimestampType(1544846469), user.SignupTimestamp)
	assert.Equal(t, false, user.IsRatingHidden)
	assert.Equal(t, "https://cs8.pikabu.ru/avatars/2561/x2561615-512432259.png", user.AvatarURL)
	// TODO: check awards
	// assert.Equal(t, true, user.Awards)
	assert.Equal(t, "approved User", user.ApprovedText)
	// TODO: check communities
	assert.Equal(t, int32(1001), user.NumberOfSubscribers)
	assert.Equal(t, true, user.IsBanned)
	assert.Equal(t, false, user.IsPermanentlyBanned)
	// TODO: check public ban history
	assert.Equal(t, models.TimestampType(100), user.BanEndTimestamp)
	assert.Equal(t, models.TimestampType(100), user.AddedTimestamp)
	assert.Equal(t, models.TimestampType(555), user.LastUpdateTimestamp)

	err = pushTaskToQueue([]byte(`
{
	"parsing_timestamp": 1000,
	"data": {
		"user": {
			"current_user_id": 0,
			"user_id": "2561615",
			"user_name": "Pisacavtor",
			"rating": "5",
			"gender": "6",
			"comments_count": 9,
			"stories_count": 2,
			"stories_hot_count": "1",
			"pluses_count": 5,
			"minuses_count": 9,
			"signup_date": "1544846469",
			"is_rating_ban": false,
			"avatar": "https://cs8.pikabu.ru/avatars/2561/x2561615-512432259.png",
			"awards": [],
			"is_subscribed": false,
			"is_ignored": false,
			"note": null,
			"approved": "approved User",
			"communities": [],
			"subscribers_count": 1001,
			"is_user_banned": true,
			"is_user_fully_banned": false,
			"public_ban_history": [
			  {
				"id": "151513",
				"date": 1544854692,
				"moderator_id": "1836690",
				"comment_id": "0",
				"comment_desc": "",
				"story_id": "6354471",
				"user_id": "2561615",
				"reason": "Отсутствие пруфа или неподтверждённая/искажённая информация (вброс)",
				"reason_id": "94",
				"story_url": "https://pikabu.ru/story/3_chasa_pyitok_6354471",
				"moderator_name": "depotato",
				"moderator_avatar": "https://cs5.pikabu.ru/avatars/1836/s1836690-1399622318.png",
				"reason_limit": null,
				"reason_count": null,
				"reason_title": null
			  }
			],
			"user_ban_time": 100
		}
	}
}
`,
	))
	if err != nil {
		panic(err)
	}

	// TODO: wait for queue to become empty
	time.Sleep(1 * time.Second)

	user = &models.PikabuUser{
		PikabuId: 2561615,
	}
	err = models.Db.Select(user)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "Pisacavtor", user.Username)
	assert.Equal(t, float32(5), user.Rating)
	assert.Equal(t, "6", user.Gender)
	assert.Equal(t, int32(9), user.NumberOfComments)
	assert.Equal(t, int32(2), user.NumberOfStories)
	assert.Equal(t, int32(1), user.NumberOfHotStories)
	assert.Equal(t, int32(5), user.NumberOfPluses)
	assert.Equal(t, int32(9), user.NumberOfMinuses)
	assert.Equal(t, models.TimestampType(1544846469), user.SignupTimestamp)
	assert.Equal(t, false, user.IsRatingHidden)
	assert.Equal(t, "https://cs8.pikabu.ru/avatars/2561/x2561615-512432259.png", user.AvatarURL)
	// TODO: check awards
	// assert.Equal(t, true, user.Awards)
	assert.Equal(t, "approved User", user.ApprovedText)
	// TODO: check communities
	assert.Equal(t, int32(1001), user.NumberOfSubscribers)
	assert.Equal(t, true, user.IsBanned)
	assert.Equal(t, false, user.IsPermanentlyBanned)
	// TODO: check public ban history
	assert.Equal(t, models.TimestampType(100), user.BanEndTimestamp)
	assert.Equal(t, models.TimestampType(100), user.AddedTimestamp)
	assert.Equal(t, models.TimestampType(1000), user.LastUpdateTimestamp)

	err = pushTaskToQueue([]byte(`
{
	"parsing_timestamp": 1500,
	"data": {
		"user": {
			"current_user_id": 0,
			"user_id": "2561615",
			"user_name": "Pisacavtor",
			"rating": "5",
			"gender": "6",
			"comments_count": 9,
			"stories_count": 2,
			"stories_hot_count": "1",
			"pluses_count": 5,
			"minuses_count": 9,
			"signup_date": "1544846469",
			"is_rating_ban": true,
			"avatar": "https://cs8.pikabu.ru/avatars/2561/x2561615-512432259.png",
			"awards": [],
			"is_subscribed": false,
			"is_ignored": false,
			"note": null,
			"approved": "approved User",
			"communities": [],
			"subscribers_count": 1001,
			"is_user_banned": true,
			"is_user_fully_banned": false,
			"public_ban_history": [
			  {
				"id": "151513",
				"date": 1544854692,
				"moderator_id": "1836690",
				"comment_id": "0",
				"comment_desc": "",
				"story_id": "6354471",
				"user_id": "2561615",
				"reason": "Отсутствие пруфа или неподтверждённая/искажённая информация (вброс)",
				"reason_id": "94",
				"story_url": "https://pikabu.ru/story/3_chasa_pyitok_6354471",
				"moderator_name": "depotato",
				"moderator_avatar": "https://cs5.pikabu.ru/avatars/1836/s1836690-1399622318.png",
				"reason_limit": null,
				"reason_count": null,
				"reason_title": null
			  }
			],
			"user_ban_time": 100
		}
	}
}
`,
	))
	if err != nil {
		panic(err)
	}

	// TODO: wait for queue to become empty
	time.Sleep(1 * time.Second)

	user = &models.PikabuUser{
		PikabuId: 2561615,
	}
	err = models.Db.Select(user)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "Pisacavtor", user.Username)
	assert.Equal(t, float32(5), user.Rating)
	assert.Equal(t, "6", user.Gender)
	assert.Equal(t, int32(9), user.NumberOfComments)
	assert.Equal(t, int32(2), user.NumberOfStories)
	assert.Equal(t, int32(1), user.NumberOfHotStories)
	assert.Equal(t, int32(5), user.NumberOfPluses)
	assert.Equal(t, int32(9), user.NumberOfMinuses)
	assert.Equal(t, models.TimestampType(1544846469), user.SignupTimestamp)
	assert.Equal(t, true, user.IsRatingHidden)
	assert.Equal(t, "https://cs8.pikabu.ru/avatars/2561/x2561615-512432259.png", user.AvatarURL)
	// TODO: check awards
	// assert.Equal(t, true, user.Awards)
	assert.Equal(t, "approved User", user.ApprovedText)
	// TODO: check communities
	assert.Equal(t, int32(1001), user.NumberOfSubscribers)
	assert.Equal(t, true, user.IsBanned)
	assert.Equal(t, false, user.IsPermanentlyBanned)
	// TODO: check public ban history
	assert.Equal(t, models.TimestampType(100), user.BanEndTimestamp)
	assert.Equal(t, models.TimestampType(100), user.AddedTimestamp)
	assert.Equal(t, models.TimestampType(1500), user.LastUpdateTimestamp)

	usernameVersions := []models.PikabuUserUsernameVersion{}
	err = models.Db.Model(&usernameVersions).
		Where("item_id = ?", user.PikabuId).
		Order("timestamp").
		Select()
	if err != nil {
		panic(err)
	}

	assert.Equal(t, usernameVersions, []models.PikabuUserUsernameVersion{
		{models.StringFieldVersion{
			models.FieldVersionBase{
				100,
				user.PikabuId,
			},
			"Pisacavtor"}},
		{models.StringFieldVersion{
			models.FieldVersionBase{
				201,
				user.PikabuId,
			},
			"Pisacavtor1"}},
		{models.StringFieldVersion{
			models.FieldVersionBase{
				555,
				user.PikabuId,
			},
			"Pisacavtor"}},
	})

	ratingVersions := []models.PikabuUserRatingVersion{}
	err = models.Db.Model(&ratingVersions).
		Where("item_id = ?", user.PikabuId).
		Order("timestamp").
		Select()
	if err != nil {
		panic(err)
	}

	assert.Equal(t, ratingVersions, []models.PikabuUserRatingVersion{
		{models.Float32FieldVersion{models.FieldVersionBase{
			100, user.PikabuId,
		}, -3.5}},
		{models.Float32FieldVersion{models.FieldVersionBase{
			201, user.PikabuId,
		}, 10.5}},
		{models.Float32FieldVersion{models.FieldVersionBase{
			555, user.PikabuId,
		}, 5}},
	})

	isRatingHiddenVersions := []models.PikabuUserIsRatingHiddenVersion{}
	err = models.Db.Model(&isRatingHiddenVersions).
		Where("item_id = ?", user.PikabuId).
		Order("timestamp").
		Select()
	if err != nil {
		panic(err)
	}

	assert.Equal(t, isRatingHiddenVersions, []models.PikabuUserIsRatingHiddenVersion{
		{models.BoolFieldVersion{models.FieldVersionBase{
			100, user.PikabuId,
		}, true}},
		{models.BoolFieldVersion{models.FieldVersionBase{
			201, user.PikabuId,
		}, false}},
		{models.BoolFieldVersion{models.FieldVersionBase{
			1000, user.PikabuId,
		}, false}},
		{models.BoolFieldVersion{models.FieldVersionBase{
			1500, user.PikabuId,
		}, true}},
	})

	userBanEndTimeVersions := []models.PikabuUserBanEndTimestampVersion{}
	err = models.Db.Model(&userBanEndTimeVersions).
		Where("item_id = ?", user.PikabuId).
		Order("timestamp").
		Select()
	if err != nil {
		panic(err)
	}

	assert.Equal(t, userBanEndTimeVersions, []models.PikabuUserBanEndTimestampVersion{
		{models.TimestampTypeFieldVersion{models.FieldVersionBase{
			100, user.PikabuId,
		}, models.TimestampType(1545459492)}},
		{models.TimestampTypeFieldVersion{models.FieldVersionBase{
			201, user.PikabuId,
		}, models.TimestampType(1545459492)}},
		{models.TimestampTypeFieldVersion{models.FieldVersionBase{
			555, user.PikabuId,
		}, models.TimestampType(100)}},
	})

	// TODO: test pikago.UserProfile serialization
}

func pushTaskToQueue(message []byte) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"parser_results",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"parser_results",
		"user_profile",
		true,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         message,
		},
	)
	return err
}
