package parse_pikabu

import (
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/results_processor"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg/orm"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func handleError(err error) {
	if err, ok := err.(*errors.Error); ok {
		panic(err.ErrorStack())
	}
	panic(err)
}

func TestUserProfileParsing(t *testing.T) {
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

	var wg sync.WaitGroup
	wg.Add(1)

	// start server
	/*go func() {
		err := server.Run()
		if err != nil {
			handleError(err)
		}
		wg.Done()
	}()

	// start task manager
	go func() {
		err := task_manager.Run()
		if err != nil {
			handleError(err)
		}
		wg.Done()
	}()*/

	// start results processor
	go func() {
		err := results_processor.Run()
		if err != nil {
			handleError(err)
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
			"awards": [
        {                                                                                    
          "id": "287578",           
          "user_id": "10080",
          "award_id": "0",                                                                
          "award_title": "Пятничное [Моё]",                                                  
          "award_image": "https://cs10.pikabu.ru/post_img/2018/04/05/8/152293551235288152.png",
          "story_id": "5983144",                                                             
          "story_title": "Впихнуть невпихуемое ",
          "date": "2018-06-25 11:43:55",
          "is_hidden": "0",                                                              
          "comment_id": null,                                                               
          "link": "/story/vpikhnut_nevpikhuemoe_5983144"                                 
        },                                                                                  
        {                                                                   
          "id": "269145",             
          "user_id": "10080",                                                             
          "award_id": "14",                                                                  
          "award_title": "редактирование тегов в 100 и более постах",                     
          "award_image": "https://cs.pikabu.ru/images/awards/2x/100_story_edits_tags.png",   
          "story_id": "0",          
          "story_title": "",
          "date": "2018-05-28 17:00:12",
          "is_hidden": "0",
          "comment_id": null,
          "link": "/edits"
        }
			],
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
		handleError(err)
	}

	// TODO: wait for queue to become empty
	time.Sleep(1 * time.Second)

	user := &models.PikabuUser{
		PikabuId: 2561615,
	}
	err = models.Db.Select(user)
	if err != nil {
		handleError(err)
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
	// assert.Equal(t, true, user.Awards)
	assert.Equal(t, "approved User", user.ApprovedText)
	assert.Equal(t, int32(1001), user.NumberOfSubscribers)
	assert.Equal(t, true, user.IsBanned)
	assert.Equal(t, false, user.IsPermanentlyBanned)
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
			"awards": [
        {                                                                                    
          "id": "287578",           
          "user_id": "10080",
          "award_id": "0",                                                                
          "award_title": "Пятничное [Моё]",                                                  
          "award_image": "https://cs10.pikabu.ru/post_img/2018/04/05/8/152293551235288152.png",
          "story_id": "5983144",                                                             
          "story_title": "Впихнуть невпихуемое ",
          "date": "2018-06-25 11:43:55",
          "is_hidden": "0",                                                              
          "comment_id": null,                                                               
          "link": "/story/vpikhnut_nevpikhuemoe_5983144"                                 
        },                                                                                  
        {                                                                   
          "id": "269145",             
          "user_id": "10080",                                                             
          "award_id": "14",                                                                  
          "award_title": "редактирование тегов в 100 и более постах",                     
          "award_image": "https://cs.pikabu.ru/images/awards/2x/100_story_edits_tags.png",   
          "story_id": "0",          
          "story_title": "",
          "date": "2018-05-28 17:00:12",
          "is_hidden": "0",
          "comment_id": null,
          "link": "/edits"
        },
		{
          "id": "252211",
          "user_id": "10080",
          "award_id": "0",
          "award_title": "Лучший вопрос на Прямой линии",
          "award_image": "https://cs10.pikabu.ru/post_img/2018/04/23/6/152447447033039417.png",
          "story_id": "4563678",
          "story_title": "Прямая линия #7",
          "date": "2018-05-02 21:23:49",
          "is_hidden": "0",
          "comment_id": null,
          "link": "/story/pryamaya_liniya_7_4563678"
        }
			],
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
		handleError(err)
	}

	// TODO: wait for queue to become empty
	time.Sleep(1 * time.Second)

	user = &models.PikabuUser{
		PikabuId: 2561615,
	}
	err = models.Db.Select(user)
	if err != nil {
		handleError(err)
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
	// assert.Equal(t, true, user.Awards)
	assert.Equal(t, "approved User", user.ApprovedText)
	assert.Equal(t, int32(1001), user.NumberOfSubscribers)
	assert.Equal(t, true, user.IsBanned)
	assert.Equal(t, false, user.IsPermanentlyBanned)
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
		handleError(err)
	}

	// TODO: wait for queue to become empty
	time.Sleep(1 * time.Second)

	user = &models.PikabuUser{
		PikabuId: 2561615,
	}
	err = models.Db.Select(user)
	if err != nil {
		handleError(err)
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
	assert.Equal(t, "approved User", user.ApprovedText)
	assert.Equal(t, int32(1001), user.NumberOfSubscribers)
	assert.Equal(t, true, user.IsBanned)
	assert.Equal(t, false, user.IsPermanentlyBanned)
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
		handleError(err)
	}

	// TODO: wait for queue to become empty
	time.Sleep(1 * time.Second)

	user = &models.PikabuUser{
		PikabuId: 2561615,
	}
	err = models.Db.Select(user)
	if err != nil {
		handleError(err)
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
	assert.Equal(t, "approved User", user.ApprovedText)
	assert.Equal(t, int32(1001), user.NumberOfSubscribers)
	assert.Equal(t, true, user.IsBanned)
	assert.Equal(t, false, user.IsPermanentlyBanned)
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
		handleError(err)
	}

	// TODO: wait for queue to become empty
	time.Sleep(1 * time.Second)

	user = &models.PikabuUser{
		PikabuId: 2561615,
	}
	err = models.Db.Select(user)
	if err != nil {
		handleError(err)
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
	assert.Equal(t, "approved User", user.ApprovedText)
	assert.Equal(t, int32(1001), user.NumberOfSubscribers)
	assert.Equal(t, true, user.IsBanned)
	assert.Equal(t, false, user.IsPermanentlyBanned)
	assert.Equal(t, models.TimestampType(100), user.BanEndTimestamp)
	assert.Equal(t, models.TimestampType(100), user.AddedTimestamp)
	assert.Equal(t, models.TimestampType(1500), user.LastUpdateTimestamp)

	usernameVersions := []models.PikabuUserUsernameVersion{}
	err = models.Db.Model(&usernameVersions).
		Where("item_id = ?", user.PikabuId).
		Order("timestamp").
		Select()
	if err != nil {
		handleError(err)
	}

	assert.Equal(t, []models.PikabuUserUsernameVersion{
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
	}, usernameVersions)

	ratingVersions := []models.PikabuUserRatingVersion{}
	err = models.Db.Model(&ratingVersions).
		Where("item_id = ?", user.PikabuId).
		Order("timestamp").
		Select()
	if err != nil {
		handleError(err)
	}

	assert.Equal(t, []models.PikabuUserRatingVersion{
		{models.Float32FieldVersion{models.FieldVersionBase{
			100, user.PikabuId,
		}, -3.5}},
		{models.Float32FieldVersion{models.FieldVersionBase{
			201, user.PikabuId,
		}, 10.5}},
		{models.Float32FieldVersion{models.FieldVersionBase{
			555, user.PikabuId,
		}, 5}},
	}, ratingVersions)

	isRatingHiddenVersions := []models.PikabuUserIsRatingHiddenVersion{}
	err = models.Db.Model(&isRatingHiddenVersions).
		Where("item_id = ?", user.PikabuId).
		Order("timestamp").
		Select()
	if err != nil {
		handleError(err)
	}

	assert.Equal(t, []models.PikabuUserIsRatingHiddenVersion{
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
	}, isRatingHiddenVersions)

	userBanEndTimeVersions := []models.PikabuUserBanEndTimestampVersion{}
	err = models.Db.Model(&userBanEndTimeVersions).
		Where("item_id = ?", user.PikabuId).
		Order("timestamp").
		Select()
	if err != nil {
		handleError(err)
	}

	assert.Equal(t, []models.PikabuUserBanEndTimestampVersion{
		{models.TimestampTypeFieldVersion{models.FieldVersionBase{
			100, user.PikabuId,
		}, models.TimestampType(1545459492)}},
		{models.TimestampTypeFieldVersion{models.FieldVersionBase{
			201, user.PikabuId,
		}, models.TimestampType(1545459492)}},
		{models.TimestampTypeFieldVersion{models.FieldVersionBase{
			555, user.PikabuId,
		}, models.TimestampType(100)}},
	}, userBanEndTimeVersions)

	// TODO: check awards
	assert.Equal(t, []uint64{}, user.Awards)

	awardVersions := []models.PikabuUserAwardsVersion{}
	err = models.Db.Model(&awardVersions).
		Where("item_id = ?", user.PikabuId).
		Order("timestamp").
		Select()
	if err != nil {
		handleError(err)
	}

	assert.Equal(t, []models.PikabuUserAwardsVersion{
		{models.FieldVersionBase{
			100, user.PikabuId,
		}, []uint64{287578, 269145}},
		{models.FieldVersionBase{
			201, user.PikabuId,
		}, []uint64{287578, 269145, 252211}},
		{models.FieldVersionBase{
			555, user.PikabuId,
		}, []uint64{}},
	}, awardVersions)
	// TODO: check communities
	// TODO: check public ban history

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
