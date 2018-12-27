package main

import (
	"bitbucket.org/d3dev/parse_pikabu/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/results_processor"
	"github.com/go-pg/pg/orm"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestUserProfileParsing(t *testing.T) {
	logger.Log.Debug(`start test "user profile parsing"`)

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

	// start results processor
	go func() {
		err := results_processor.Run()
		if err != nil {
			handleError(err)
		}
		wg.Done()
	}()

	err = pushTaskToQueue("user_profile", []byte(`
{
	"parsing_timestamp": 100,
	"parser_id": "d3dev/parser_id",
	"number_of_results": 1,
	"results": [{
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
          "user_id": "2561615",
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
          "user_id": "2561615",                                                             
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
	}]
}
`,
	))
	if err != nil {
		handleError(err)
	}

	waitForQueueEmpty()

	user := &models.PikabuUser{
		PikabuId: 2561615,
	}
	err = models.Db.Select(user)
	if err != nil {
		handleError(err)
	}
	assert.Equal(t, "Pisacavtor", user.Username)
	assert.Equal(t, int32(-3), user.Rating)
	assert.Equal(t, "6", user.Gender)
	assert.Equal(t, int32(3), user.NumberOfComments)
	assert.Equal(t, int32(2), user.NumberOfStories)
	assert.Equal(t, int32(1), user.NumberOfHotStories)
	assert.Equal(t, int32(5), user.NumberOfPluses)
	assert.Equal(t, int32(9), user.NumberOfMinuses)
	assert.Equal(t, models.TimestampType(1544846469), user.SignupTimestamp)
	assert.Equal(t, true, user.IsRatingHidden)
	assert.Equal(t, "https://cs8.pikabu.ru/avatars/2561/x2561615-512432259.png", user.AvatarURL)
	// assert.Equal(t, true, user.AwardIds)
	assert.Equal(t, "approved User", user.ApprovedText)
	assert.Equal(t, int32(1001), user.NumberOfSubscribers)
	assert.Equal(t, true, user.IsBanned)
	assert.Equal(t, false, user.IsPermanentlyBanned)
	assert.Equal(t, models.TimestampType(1545459492), user.BanEndTimestamp)
	assert.Equal(t, models.TimestampType(100), user.AddedTimestamp)
	assert.Equal(t, models.TimestampType(100), user.LastUpdateTimestamp)

	// change some fields
	err = pushTaskToQueue("user_profile", []byte(`
{
	"parsing_timestamp": 201,
	"parser_id": "d3dev/parser_id",
	"number_of_results": 1,
	"results": [{
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
          "user_id": "2561615",
          "award_id": "0",                                                                
          "award_title": "Пятничное [твоё]",                                                  
          "award_image": "https://cs10.pikabu.ru/post_img/2018/04/05/8/152293551235288152.png",
          "story_id": "9983144",                                                             
          "story_title": "Впихнуть невпихуемое ",
          "date": "2018-06-25 11:43:55",
          "is_hidden": "0",                                                              
          "comment_id": null,                                                               
          "link": "/story/vpikhnut_nevpikhuemoe_5983144"                                 
        },                                                                                  
        {                                                                   
          "id": "269145",             
          "user_id": "2561615",                                                             
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
          "user_id": "2561615",
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
	}]
}
`,
	))
	if err != nil {
		handleError(err)
	}

	waitForQueueEmpty()

	user = &models.PikabuUser{
		PikabuId: 2561615,
	}
	err = models.Db.Select(user)
	if err != nil {
		handleError(err)
	}
	assert.Equal(t, "Pisacavtor1", user.Username)
	assert.Equal(t, int32(10), user.Rating)
	assert.Equal(t, "6", user.Gender)
	assert.Equal(t, int32(9), user.NumberOfComments)
	assert.Equal(t, int32(2), user.NumberOfStories)
	assert.Equal(t, int32(1), user.NumberOfHotStories)
	assert.Equal(t, int32(5), user.NumberOfPluses)
	assert.Equal(t, int32(9), user.NumberOfMinuses)
	assert.Equal(t, models.TimestampType(1544846469), user.SignupTimestamp)
	assert.Equal(t, false, user.IsRatingHidden)
	assert.Equal(t, "https://cs8.pikabu.ru/avatars/2561/x2561615-512432259.png", user.AvatarURL)
	// assert.Equal(t, true, user.AwardIds)
	assert.Equal(t, "approved User", user.ApprovedText)
	assert.Equal(t, int32(1001), user.NumberOfSubscribers)
	assert.Equal(t, true, user.IsBanned)
	assert.Equal(t, false, user.IsPermanentlyBanned)
	assert.Equal(t, models.TimestampType(1545459492), user.BanEndTimestamp)
	assert.Equal(t, models.TimestampType(100), user.AddedTimestamp)
	assert.Equal(t, models.TimestampType(201), user.LastUpdateTimestamp)

	err = pushTaskToQueue("user_profile", []byte(`
{
	"parsing_timestamp": 555,
	"parser_id": "d3dev/parser_id",
	"number_of_results": 1,
	"results": [{
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
	}]
}
`,
	))
	if err != nil {
		handleError(err)
	}

	waitForQueueEmpty()

	user = &models.PikabuUser{
		PikabuId: 2561615,
	}
	err = models.Db.Select(user)
	if err != nil {
		handleError(err)
	}
	assert.Equal(t, "Pisacavtor", user.Username)
	assert.Equal(t, int32(5), user.Rating)
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

	err = pushTaskToQueue("user_profile", []byte(`
{
	"parsing_timestamp": 1000,
	"parser_id": "d3dev/parser_id",
	"number_of_results": 1,
	"results": [{
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
	}]
}
`,
	))
	if err != nil {
		handleError(err)
	}

	waitForQueueEmpty()

	user = &models.PikabuUser{
		PikabuId: 2561615,
	}
	err = models.Db.Select(user)
	if err != nil {
		handleError(err)
	}
	assert.Equal(t, "Pisacavtor", user.Username)
	assert.Equal(t, int32(5), user.Rating)
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

	err = pushTaskToQueue("user_profile", []byte(`
{
	"parsing_timestamp": 1500,
	"parser_id": "d3dev/parser_id",
	"number_of_results": 1,
	"results": [{
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
			"communities": [
        {                                                                                                                                                                                       
          "name": "Cynic Mansion",                                                                                                                                                              
          "link": "cynicmansion",                                                                                                                                                               
          "avatar": "https://cs6.pikabu.ru/images/community/1031/1502225712241040050.png",                                                                                                      
          "avatar_url": "https://cs6.pikabu.ru/images/community/1031/1502225712241040050.png"                                                                                                   
        },                
        {                    
          "name": "Пикабу головного мозга",
          "link": "p_g_m",                 
          "avatar": "https://cs7.pikabu.ru/images/community/1360/1538729487212641089.png",     
          "avatar_url": "https://cs7.pikabu.ru/images/community/1360/1538729487212641089.png"
        },                                       
        {                               
          "name": "Кофе мой друг",
          "link": "Coffee",  
          "avatar": "https://cs8.pikabu.ru/images/community/729/1493440472283550654.png",
          "avatar_url": "https://cs8.pikabu.ru/images/community/729/1493440472283550654.png"
        }
			],
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
			  },
			  {
				"id": "151514",
				"date": 1544854693,
				"moderator_id": "1836691",
				"comment_id": "15",
				"comment_desc": "",
				"story_id": "6354471",
				"user_id": "2561615",
				"reason": "Отсутствие пруфа или неподтверждённая/искажённая информация (вброс)",
				"reason_id": "94",
				"story_url": "https://pikabu.ru/story/3_chasa_pyitok_6354471",
				"moderator_name": "nepotato",
				"moderator_avatar": "https://cs5.pikabu.ru/avatars/1836/s1836690-1399622318.png",
				"reason_limit": null,
				"reason_count": null,
				"reason_title": null
			  }
			],
			"user_ban_time": 100
		}
	}]
}
`,
	))
	if err != nil {
		handleError(err)
	}

	waitForQueueEmpty()

	user = &models.PikabuUser{
		PikabuId: 2561615,
	}
	err = models.Db.Select(user)
	if err != nil {
		handleError(err)
	}
	assert.Equal(t, "Pisacavtor", user.Username)
	assert.Equal(t, int32(5), user.Rating)
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
		{
			Timestamp: 100,
			ItemId:    user.PikabuId,
			Value:     "Pisacavtor",
		},
		{
			Timestamp: 201,
			ItemId:    user.PikabuId,
			Value:     "Pisacavtor1",
		},
		{
			Timestamp: 555,
			ItemId:    user.PikabuId,
			Value:     "Pisacavtor",
		},
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
		{Timestamp: 100, ItemId: user.PikabuId, Value: -3},
		{Timestamp: 201, ItemId: user.PikabuId, Value: 10},
		{Timestamp: 555, ItemId: user.PikabuId, Value: 5},
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
		{Timestamp: 100, ItemId: user.PikabuId, Value: true},
		{Timestamp: 201, ItemId: user.PikabuId, Value: false},
		{Timestamp: 1000, ItemId: user.PikabuId, Value: false},
		{Timestamp: 1500, ItemId: user.PikabuId, Value: true},
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
		{Timestamp: 100, ItemId: user.PikabuId, Value: 1545459492},
		{Timestamp: 201, ItemId: user.PikabuId, Value: 1545459492},
		{Timestamp: 555, ItemId: user.PikabuId, Value: 100},
	}, userBanEndTimeVersions)

	// check awards
	assert.Equal(t, []uint64{}, user.AwardIds)

	awards := []models.PikabuUserAward{}
	err = models.Db.Model(&awards).
		Where("user_id = ?", user.PikabuId).
		Order("pikabu_id").
		Select()
	if err != nil {
		handleError(err)
	}

	assert.Equal(t, []models.PikabuUserAward{
		{
			PikabuId:            252211,
			AddedTimestamp:      201,
			UserId:              2561615,
			AwardId:             0,
			AwardTitle:          "Лучший вопрос на Прямой линии",
			AwardImageURL:       "https://cs10.pikabu.ru/post_img/2018/04/23/6/152447447033039417.png",
			StoryId:             4563678,
			StoryTitle:          "Прямая линия #7",
			IssuingDate:         "2018-05-02 21:23:49",
			IsHidden:            false,
			CommentId:           0,
			Link:                "/story/pryamaya_liniya_7_4563678",
			LastUpdateTimestamp: 201,
		},
		{
			PikabuId:            269145,
			AddedTimestamp:      100,
			UserId:              2561615,
			AwardId:             14,
			AwardTitle:          "редактирование тегов в 100 и более постах",
			AwardImageURL:       "https://cs.pikabu.ru/images/awards/2x/100_story_edits_tags.png",
			StoryId:             0,
			StoryTitle:          "",
			IssuingDate:         "2018-05-28 17:00:12",
			IsHidden:            false,
			CommentId:           0,
			Link:                "/edits",
			LastUpdateTimestamp: 201,
		},
		{
			PikabuId:            287578,
			AddedTimestamp:      100,
			UserId:              2561615,
			AwardId:             0,
			AwardTitle:          "Пятничное [твоё]",
			AwardImageURL:       "https://cs10.pikabu.ru/post_img/2018/04/05/8/152293551235288152.png",
			StoryId:             9983144,
			StoryTitle:          "Впихнуть невпихуемое ",
			IssuingDate:         "2018-06-25 11:43:55",
			IsHidden:            false,
			CommentId:           0,
			Link:                "/story/vpikhnut_nevpikhuemoe_5983144",
			LastUpdateTimestamp: 201,
		},
	}, awards)

	awardVersions := []models.PikabuUserAwardIdsVersion{}
	err = models.Db.Model(&awardVersions).
		Where("item_id = ?", user.PikabuId).
		Order("timestamp").
		Select()
	if err != nil {
		handleError(err)
	}

	assert.Equal(t, []models.PikabuUserAwardIdsVersion{
		{Timestamp: 100, ItemId: user.PikabuId, Value: []uint64{287578, 269145}},
		{Timestamp: 201, ItemId: user.PikabuId, Value: []uint64{287578, 269145, 252211}},
		{Timestamp: 555, ItemId: user.PikabuId, Value: []uint64{}},
	}, awardVersions)

	// award fields versions
	awardTitleVersions := []models.PikabuUserAwardAwardTitleVersion{}
	err = models.Db.Model(&awardTitleVersions).
		Where("item_id = ?", 287578).
		Select()
	handleError(err)

	assert.Equal(t, []models.PikabuUserAwardAwardTitleVersion{
		{Timestamp: 100, ItemId: 287578, Value: "Пятничное [Моё]"},
		{Timestamp: 201, ItemId: 287578, Value: "Пятничное [твоё]"},
	}, awardTitleVersions)

	awardStoryIdVersions := []models.PikabuUserAwardStoryIdVersion{}
	err = models.Db.Model(&awardStoryIdVersions).
		Where("item_id = ?", 287578).
		Select()
	handleError(err)

	assert.Equal(t, []models.PikabuUserAwardStoryIdVersion{
		{Timestamp: 100, ItemId: 287578, Value: 5983144},
		{Timestamp: 201, ItemId: 287578, Value: 9983144},
	}, awardStoryIdVersions)

	// check communities
	communities := []models.PikabuUserCommunity{}
	err = models.Db.Model(&communities).Select()
	if err != nil {
		handleError(err)
	}

	assert.Equal(t, []models.PikabuUserCommunity{
		{
			Id:                  1,
			Name:                "Cynic Mansion",
			Link:                "cynicmansion",
			AvatarURL:           "https://cs6.pikabu.ru/images/community/1031/1502225712241040050.png",
			AddedTimestamp:      1500,
			LastUpdateTimestamp: 1500,
		},
		{
			Id:                  2,
			Name:                "Пикабу головного мозга",
			Link:                "p_g_m",
			AvatarURL:           "https://cs7.pikabu.ru/images/community/1360/1538729487212641089.png",
			AddedTimestamp:      1500,
			LastUpdateTimestamp: 1500,
		},
		{
			Id:                  3,
			Name:                "Кофе мой друг",
			Link:                "Coffee",
			AvatarURL:           "https://cs8.pikabu.ru/images/community/729/1493440472283550654.png",
			AddedTimestamp:      1500,
			LastUpdateTimestamp: 1500,
		},
	}, communities)

	communityVersions := []models.PikabuUserCommunityIdsVersion{}
	err = models.Db.Model(&communityVersions).
		Where("item_id = ?", user.PikabuId).
		Order("timestamp").
		Select()
	if err != nil {
		handleError(err)
	}

	assert.Equal(t, []models.PikabuUserCommunityIdsVersion{
		{Timestamp: 100, ItemId: user.PikabuId, Value: []uint64{}},
		{Timestamp: 1000, ItemId: user.PikabuId, Value: []uint64{}},
		{Timestamp: 1500, ItemId: user.PikabuId, Value: []uint64{1, 2, 3}},
	}, communityVersions)

	// check public ban history
	banHistoryItems := []models.PikabuUserBanHistoryItem{}
	err = models.Db.Model(&banHistoryItems).
		Where("user_id = ?", user.PikabuId).
		Order("added_timestamp").
		Select()
	if err != nil {
		handleError(err)
	}

	assert.Equal(t, []models.PikabuUserBanHistoryItem{
		{
			PikabuId:                151513,
			BanStartTimestamp:       1544854692,
			CommentId:               0,
			CommentHtmlDeleteReason: "",
			StoryId:                 6354471,
			UserId:                  2561615,
			BanReason:               "Отсутствие пруфа или неподтверждённая/искажённая информация (вброс)",
			BanReasonId:             94,
			StoryURL:                "https://pikabu.ru/story/3_chasa_pyitok_6354471",
			ModeratorId:             1836690,
			ModeratorName:           "depotato",
			ModeratorAvatar:         "https://cs5.pikabu.ru/avatars/1836/s1836690-1399622318.png",
			ReasonsLimit:            0,
			ReasonCount:             0,
			ReasonTitle:             "",
			AddedTimestamp:          100,
			LastUpdateTimestamp:     1500,
		},
		{
			PikabuId:                151514,
			BanStartTimestamp:       1544854693,
			CommentId:               15,
			CommentHtmlDeleteReason: "",
			StoryId:                 6354471,
			UserId:                  2561615,
			BanReason:               "Отсутствие пруфа или неподтверждённая/искажённая информация (вброс)",
			BanReasonId:             94,
			StoryURL:                "https://pikabu.ru/story/3_chasa_pyitok_6354471",
			ModeratorId:             1836691,
			ModeratorName:           "nepotato",
			ModeratorAvatar:         "https://cs5.pikabu.ru/avatars/1836/s1836690-1399622318.png",
			ReasonsLimit:            0,
			ReasonCount:             0,
			ReasonTitle:             "",
			AddedTimestamp:          1500,
			LastUpdateTimestamp:     1500,
		},
	}, banHistoryItems)

	banHistoryItemVersions := []models.PikabuUserBanHistoryItemIdsVersion{}
	err = models.Db.Model(&banHistoryItemVersions).
		Where("item_id = ?", user.PikabuId).
		Select()
	if err != nil {
		handleError(err)
	}

	assert.Equal(t, []models.PikabuUserBanHistoryItemIdsVersion{
		{Timestamp: 100, ItemId: user.PikabuId, Value: []uint64{151513}},
		{Timestamp: 1000, ItemId: user.PikabuId, Value: []uint64{151513}},
		{Timestamp: 1500, ItemId: user.PikabuId, Value: []uint64{151513, 151514}},
	}, banHistoryItemVersions)

	err = pushTaskToQueue("user_profile", []byte(`
{
	"parsing_timestamp": 1501,
	"parser_id": "d3dev/parser_id",
	"number_of_results": 1,
	"results": [{
		"user": {
			"current_user_id": 0,
			"user_id": "2561615",
			"user_name": "Pisacavtor",
			"rating": "5.5",
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
			"communities": [
        {                                                                                                                                                                                       
          "name": "Cynic Mansion",                                                                                                                                                              
          "link": "cynicmansion",                                                                                                                                                               
          "avatar": "https://cs6.pikabu.ru/images/community/1031/1502225712241040050.png",                                                                                                      
          "avatar_url": "https://cs6.pikabu.ru/images/community/1031/1502225712241040050.png"                                                                                                   
        },                
        {                    
          "name": "Пикабу головного мозга",
          "link": "p_g_m",                 
          "avatar": "https://cs7.pikabu.ru/images/community/1360/1538729487212641089.png",     
          "avatar_url": "https://cs7.pikabu.ru/images/community/1360/1538729487212641089.png"
        },                                       
        {                               
          "name": "Кофе мой друг",
          "link": "Coffee",  
          "avatar": "https://cs8.pikabu.ru/images/community/729/1493440472283550654.png",
          "avatar_url": "https://cs8.pikabu.ru/images/community/729/1493440472283550654.png"
        }
			],
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
			  },
			  {
				"id": "151514",
				"date": 1544854693,
				"moderator_id": "1836691",
				"comment_id": "15",
				"comment_desc": "",
				"story_id": "6354471",
				"user_id": "2561615",
				"reason": "Отсутствие пруфа или неподтверждённая/искажённая информация (вброс)",
				"reason_id": "94",
				"story_url": "https://pikabu.ru/story/3_chasa_pyitok_6354471",
				"moderator_name": "nepotato",
				"moderator_avatar": "https://cs5.pikabu.ru/avatars/1836/s1836690-1399622318.png",
				"reason_limit": null,
				"reason_count": null,
				"reason_title": null
			  }
			],
			"user_ban_time": 100
		}
	}]
}
`,
	))
	if err != nil {
		handleError(err)
	}

	waitForQueueEmpty()

	user = &models.PikabuUser{
		PikabuId: 2561615,
	}
	err = models.Db.Select(user)
	if err != nil {
		handleError(err)
	}
	assert.Equal(t, "Pisacavtor", user.Username)
	assert.Equal(t, int32(5), user.Rating)
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
	assert.Equal(t, models.TimestampType(1501), user.LastUpdateTimestamp)

	ratingVersions = []models.PikabuUserRatingVersion{}
	err = models.Db.Model(&ratingVersions).
		Where("item_id = ?", user.PikabuId).
		Order("timestamp").
		Select()
	if err != nil {
		handleError(err)
	}

	assert.Equal(t, []models.PikabuUserRatingVersion{
		{Timestamp: 100, ItemId: user.PikabuId, Value: -3},
		{Timestamp: 201, ItemId: user.PikabuId, Value: 10},
		{Timestamp: 555, ItemId: user.PikabuId, Value: 5},
	}, ratingVersions)

	err = pushTaskToQueue("user_profile", []byte(`
{
	"parsing_timestamp": 1502,
	"parser_id": "d3dev/parser_id",
	"number_of_results": 1,
	"results": [{
		"user": {
			"current_user_id": 0,
			"user_id": "2561615",
			"user_name": "Pisacavtor",
			"rating": "4.5",
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
			"communities": [
        {                                                                                                                                                                                       
          "name": "Cynic Mansion",                                                                                                                                                              
          "link": "cynicmansion",                                                                                                                                                               
          "avatar": "https://cs6.pikabu.ru/images/community/1031/1502225712241040050.png",                                                                                                      
          "avatar_url": "https://cs6.pikabu.ru/images/community/1031/1502225712241040050.png"                                                                                                   
        },                
        {                    
          "name": "Пикабу головного мозга",
          "link": "p_g_m",                 
          "avatar": "https://cs7.pikabu.ru/images/community/1360/1538729487212641089.png",     
          "avatar_url": "https://cs7.pikabu.ru/images/community/1360/1538729487212641089.png"
        },                                       
        {                               
          "name": "Кофе мой друг",
          "link": "Coffee",  
          "avatar": "https://cs8.pikabu.ru/images/community/729/1493440472283550654.png",
          "avatar_url": "https://cs8.pikabu.ru/images/community/729/1493440472283550654.png"
        }
			],
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
			  },
			  {
				"id": "151514",
				"date": 1544854693,
				"moderator_id": "1836691",
				"comment_id": "15",
				"comment_desc": "",
				"story_id": "6354471",
				"user_id": "2561615",
				"reason": "Отсутствие пруфа или неподтверждённая/искажённая информация (вброс)",
				"reason_id": "94",
				"story_url": "https://pikabu.ru/story/3_chasa_pyitok_6354471",
				"moderator_name": "nepotato",
				"moderator_avatar": "https://cs5.pikabu.ru/avatars/1836/s1836690-1399622318.png",
				"reason_limit": null,
				"reason_count": null,
				"reason_title": null
			  }
			],
			"user_ban_time": 100
		}
	}]
}
`,
	))
	if err != nil {
		handleError(err)
	}

	waitForQueueEmpty()

	user = &models.PikabuUser{
		PikabuId: 2561615,
	}
	err = models.Db.Select(user)
	if err != nil {
		handleError(err)
	}
	assert.Equal(t, "Pisacavtor", user.Username)
	assert.Equal(t, int32(4), user.Rating)
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
	assert.Equal(t, models.TimestampType(1502), user.LastUpdateTimestamp)

	ratingVersions = []models.PikabuUserRatingVersion{}
	err = models.Db.Model(&ratingVersions).
		Where("item_id = ?", user.PikabuId).
		Order("timestamp").
		Select()
	if err != nil {
		handleError(err)
	}

	assert.Equal(t, []models.PikabuUserRatingVersion{
		{Timestamp: 100, ItemId: user.PikabuId, Value: -3},
		{Timestamp: 201, ItemId: user.PikabuId, Value: 10},
		{Timestamp: 555, ItemId: user.PikabuId, Value: 5},
		{Timestamp: 1501, ItemId: user.PikabuId, Value: 5},
		{Timestamp: 1502, ItemId: user.PikabuId, Value: 4},
	}, ratingVersions)
	// TODO: check pikabu user ban history items versions

	// TODO: test pikago.UserProfile serialization

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
}

func pushTaskToQueue(key string, message []byte) error {
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
		key,
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

func waitForQueueEmpty() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		handleError(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		handleError(err)
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
		handleError(err)
	}

	q, err := ch.QueueDeclare(
		"bitbucket.org/d3dev/parse_pikabu",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		handleError(err)
	}
	for q.Messages > 0 {
		time.Sleep(1 * time.Second)
	}
	time.Sleep(2 * time.Second)
}

func TestConcurrentUserProfileParsing(t *testing.T) {
	logger.Log.Debug(`start test "concurrent user profile parsing"`)

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

	// start results processor
	go func() {
		err := results_processor.Run()
		if err != nil {
			handleError(err)
		}
		wg.Done()
	}()
	data := `
{"parsing_timestamp":1545920253,"parser_id":"d3dev/test_proxy_parser_copy_95","number_of_results":1,"results":[{"user":{"current_user_id":0,"user_id":297508,"user_name":"moderator","rating":120951.5,"gender":0,"comments_count":73338,"stories_count":0,"stories_hot_count":0,"pluses_count":192,"minuses_count":383,"signup_date":1367078405,"is_rating_ban":true,"avatar":"https://cs6.pikabu.ru/avatars/297/x297508-739228827.png","awards":[{"id":249466,"user_id":297508,"award_id":13,"award_title":"5 лет на Пикабу","award_image":"https://cs.pikabu.ru/images/awards/2x/5_let_pikabushnik.png","story_id":0,"story_title":"","date":"2018-04-27 20:00:21","is_hidden":0,"comment_id":0,"link":""},{"id":1746,"user_id":297508,"award_id":8,"award_title":"более 1000 подписчиков","award_image":"https://cs.pikabu.ru/images/awards/2x/1000_subs.png","story_id":0,"story_title":"","date":"2015-07-05 06:00:06","is_hidden":0,"comment_id":0,"link":""}],"is_subscribed":false,"is_ignored":false,"note":"","approved":"Команда Пикабу","communities":[{"name":"Фейки","link":"fakes","avatar":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png","avatar_url":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png"},{"name":"Вопросы по модерации","link":"moderator","avatar":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png","avatar_url":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png"},{"name":"Баяны","link":"reposts","avatar":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png","avatar_url":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png"}],"subscribers_count":5511,"is_user_banned":false,"is_user_fully_banned":false,"public_ban_history":[],"user_ban_time":-62169993079}}]}
{"parsing_timestamp":1545920253,"parser_id":"d3dev/test_proxy_parser_copy_149","number_of_results":1,"results":[{"user":{"current_user_id":0,"user_id":297508,"user_name":"moderator","rating":120951.5,"gender":0,"comments_count":73338,"stories_count":0,"stories_hot_count":0,"pluses_count":192,"minuses_count":383,"signup_date":1367078405,"is_rating_ban":true,"avatar":"https://cs6.pikabu.ru/avatars/297/x297508-739228827.png","awards":[{"id":249466,"user_id":297508,"award_id":13,"award_title":"5 лет на Пикабу","award_image":"https://cs.pikabu.ru/images/awards/2x/5_let_pikabushnik.png","story_id":0,"story_title":"","date":"2018-04-27 20:00:21","is_hidden":0,"comment_id":0,"link":""},{"id":1746,"user_id":297508,"award_id":8,"award_title":"более 1000 подписчиков","award_image":"https://cs.pikabu.ru/images/awards/2x/1000_subs.png","story_id":0,"story_title":"","date":"2015-07-05 06:00:06","is_hidden":0,"comment_id":0,"link":""}],"is_subscribed":false,"is_ignored":false,"note":"","approved":"Команда Пикабу","communities":[{"name":"Фейки","link":"fakes","avatar":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png","avatar_url":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png"},{"name":"Вопросы по модерации","link":"moderator","avatar":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png","avatar_url":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png"},{"name":"Баяны","link":"reposts","avatar":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png","avatar_url":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png"}],"subscribers_count":5511,"is_user_banned":false,"is_user_fully_banned":false,"public_ban_history":[],"user_ban_time":-62169993079}}]}
{"parsing_timestamp":1545920253,"parser_id":"d3dev/test_proxy_parser_copy_42","number_of_results":1,"results":[{"user":{"current_user_id":0,"user_id":297508,"user_name":"moderator","rating":120951.5,"gender":0,"comments_count":73338,"stories_count":0,"stories_hot_count":0,"pluses_count":192,"minuses_count":383,"signup_date":1367078405,"is_rating_ban":true,"avatar":"https://cs6.pikabu.ru/avatars/297/x297508-739228827.png","awards":[{"id":249466,"user_id":297508,"award_id":13,"award_title":"5 лет на Пикабу","award_image":"https://cs.pikabu.ru/images/awards/2x/5_let_pikabushnik.png","story_id":0,"story_title":"","date":"2018-04-27 20:00:21","is_hidden":0,"comment_id":0,"link":""},{"id":1746,"user_id":297508,"award_id":8,"award_title":"более 1000 подписчиков","award_image":"https://cs.pikabu.ru/images/awards/2x/1000_subs.png","story_id":0,"story_title":"","date":"2015-07-05 06:00:06","is_hidden":0,"comment_id":0,"link":""}],"is_subscribed":false,"is_ignored":false,"note":"","approved":"Команда Пикабу","communities":[{"name":"Фейки","link":"fakes","avatar":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png","avatar_url":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png"},{"name":"Вопросы по модерации","link":"moderator","avatar":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png","avatar_url":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png"},{"name":"Баяны","link":"reposts","avatar":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png","avatar_url":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png"}],"subscribers_count":5511,"is_user_banned":false,"is_user_fully_banned":false,"public_ban_history":[],"user_ban_time":-62169993079}}]}
{"parsing_timestamp":1545920253,"parser_id":"d3dev/test_proxy_parser_copy_56","number_of_results":1,"results":[{"user":{"current_user_id":0,"user_id":297508,"user_name":"moderator","rating":120951.5,"gender":0,"comments_count":73338,"stories_count":0,"stories_hot_count":0,"pluses_count":192,"minuses_count":383,"signup_date":1367078405,"is_rating_ban":true,"avatar":"https://cs6.pikabu.ru/avatars/297/x297508-739228827.png","awards":[{"id":249466,"user_id":297508,"award_id":13,"award_title":"5 лет на Пикабу","award_image":"https://cs.pikabu.ru/images/awards/2x/5_let_pikabushnik.png","story_id":0,"story_title":"","date":"2018-04-27 20:00:21","is_hidden":0,"comment_id":0,"link":""},{"id":1746,"user_id":297508,"award_id":8,"award_title":"более 1000 подписчиков","award_image":"https://cs.pikabu.ru/images/awards/2x/1000_subs.png","story_id":0,"story_title":"","date":"2015-07-05 06:00:06","is_hidden":0,"comment_id":0,"link":""}],"is_subscribed":false,"is_ignored":false,"note":"","approved":"Команда Пикабу","communities":[{"name":"Фейки","link":"fakes","avatar":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png","avatar_url":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png"},{"name":"Вопросы по модерации","link":"moderator","avatar":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png","avatar_url":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png"},{"name":"Баяны","link":"reposts","avatar":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png","avatar_url":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png"}],"subscribers_count":5511,"is_user_banned":false,"is_user_fully_banned":false,"public_ban_history":[],"user_ban_time":-62169993079}}]}
{"parsing_timestamp":1545920253,"parser_id":"d3dev/test_proxy_parser_copy_111","number_of_results":1,"results":[{"user":{"current_user_id":0,"user_id":297508,"user_name":"moderator","rating":120951.5,"gender":0,"comments_count":73338,"stories_count":0,"stories_hot_count":0,"pluses_count":192,"minuses_count":383,"signup_date":1367078405,"is_rating_ban":true,"avatar":"https://cs6.pikabu.ru/avatars/297/x297508-739228827.png","awards":[{"id":249466,"user_id":297508,"award_id":13,"award_title":"5 лет на Пикабу","award_image":"https://cs.pikabu.ru/images/awards/2x/5_let_pikabushnik.png","story_id":0,"story_title":"","date":"2018-04-27 20:00:21","is_hidden":0,"comment_id":0,"link":""},{"id":1746,"user_id":297508,"award_id":8,"award_title":"более 1000 подписчиков","award_image":"https://cs.pikabu.ru/images/awards/2x/1000_subs.png","story_id":0,"story_title":"","date":"2015-07-05 06:00:06","is_hidden":0,"comment_id":0,"link":""}],"is_subscribed":false,"is_ignored":false,"note":"","approved":"Команда Пикабу","communities":[{"name":"Фейки","link":"fakes","avatar":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png","avatar_url":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png"},{"name":"Вопросы по модерации","link":"moderator","avatar":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png","avatar_url":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png"},{"name":"Баяны","link":"reposts","avatar":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png","avatar_url":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png"}],"subscribers_count":5511,"is_user_banned":false,"is_user_fully_banned":false,"public_ban_history":[],"user_ban_time":-62169993079}}]}
{"parsing_timestamp":1545920255,"parser_id":"d3dev/test_proxy_parser_copy_89","number_of_results":1,"results":[{"user":{"current_user_id":0,"user_id":297508,"user_name":"moderator","rating":120951.5,"gender":0,"comments_count":73338,"stories_count":0,"stories_hot_count":0,"pluses_count":192,"minuses_count":383,"signup_date":1367078405,"is_rating_ban":true,"avatar":"https://cs6.pikabu.ru/avatars/297/x297508-739228827.png","awards":[{"id":249466,"user_id":297508,"award_id":13,"award_title":"5 лет на Пикабу","award_image":"https://cs.pikabu.ru/images/awards/2x/5_let_pikabushnik.png","story_id":0,"story_title":"","date":"2018-04-27 20:00:21","is_hidden":0,"comment_id":0,"link":""},{"id":1746,"user_id":297508,"award_id":8,"award_title":"более 1000 подписчиков","award_image":"https://cs.pikabu.ru/images/awards/2x/1000_subs.png","story_id":0,"story_title":"","date":"2015-07-05 06:00:06","is_hidden":0,"comment_id":0,"link":""}],"is_subscribed":false,"is_ignored":false,"note":"","approved":"Команда Пикабу","communities":[{"name":"Фейки","link":"fakes","avatar":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png","avatar_url":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png"},{"name":"Вопросы по модерации","link":"moderator","avatar":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png","avatar_url":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png"},{"name":"Баяны","link":"reposts","avatar":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png","avatar_url":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png"}],"subscribers_count":5511,"is_user_banned":false,"is_user_fully_banned":false,"public_ban_history":[],"user_ban_time":-62169993079}}]}
{"parsing_timestamp":1545920256,"parser_id":"d3dev/test_proxy_parser_copy_28","number_of_results":1,"results":[{"user":{"current_user_id":0,"user_id":297508,"user_name":"moderator","rating":120951.5,"gender":0,"comments_count":73338,"stories_count":0,"stories_hot_count":0,"pluses_count":192,"minuses_count":383,"signup_date":1367078405,"is_rating_ban":true,"avatar":"https://cs6.pikabu.ru/avatars/297/x297508-739228827.png","awards":[{"id":249466,"user_id":297508,"award_id":13,"award_title":"5 лет на Пикабу","award_image":"https://cs.pikabu.ru/images/awards/2x/5_let_pikabushnik.png","story_id":0,"story_title":"","date":"2018-04-27 20:00:21","is_hidden":0,"comment_id":0,"link":""},{"id":1746,"user_id":297508,"award_id":8,"award_title":"более 1000 подписчиков","award_image":"https://cs.pikabu.ru/images/awards/2x/1000_subs.png","story_id":0,"story_title":"","date":"2015-07-05 06:00:06","is_hidden":0,"comment_id":0,"link":""}],"is_subscribed":false,"is_ignored":false,"note":"","approved":"Команда Пикабу","communities":[{"name":"Фейки","link":"fakes","avatar":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png","avatar_url":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png"},{"name":"Вопросы по модерации","link":"moderator","avatar":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png","avatar_url":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png"},{"name":"Баяны","link":"reposts","avatar":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png","avatar_url":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png"}],"subscribers_count":5511,"is_user_banned":false,"is_user_fully_banned":false,"public_ban_history":[],"user_ban_time":-62169993079}}]}
{"parsing_timestamp":1545920256,"parser_id":"d3dev/test_proxy_parser_copy_9","number_of_results":1,"results":[{"user":{"current_user_id":0,"user_id":297508,"user_name":"moderator","rating":120951.5,"gender":0,"comments_count":73338,"stories_count":0,"stories_hot_count":0,"pluses_count":192,"minuses_count":383,"signup_date":1367078405,"is_rating_ban":true,"avatar":"https://cs6.pikabu.ru/avatars/297/x297508-739228827.png","awards":[{"id":249466,"user_id":297508,"award_id":13,"award_title":"5 лет на Пикабу","award_image":"https://cs.pikabu.ru/images/awards/2x/5_let_pikabushnik.png","story_id":0,"story_title":"","date":"2018-04-27 20:00:21","is_hidden":0,"comment_id":0,"link":""},{"id":1746,"user_id":297508,"award_id":8,"award_title":"более 1000 подписчиков","award_image":"https://cs.pikabu.ru/images/awards/2x/1000_subs.png","story_id":0,"story_title":"","date":"2015-07-05 06:00:06","is_hidden":0,"comment_id":0,"link":""}],"is_subscribed":false,"is_ignored":false,"note":"","approved":"Команда Пикабу","communities":[{"name":"Фейки","link":"fakes","avatar":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png","avatar_url":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png"},{"name":"Вопросы по модерации","link":"moderator","avatar":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png","avatar_url":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png"},{"name":"Баяны","link":"reposts","avatar":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png","avatar_url":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png"}],"subscribers_count":5511,"is_user_banned":false,"is_user_fully_banned":false,"public_ban_history":[],"user_ban_time":-62169993079}}]}
{"parsing_timestamp":1545920256,"parser_id":"d3dev/test_proxy_parser_copy_45","number_of_results":1,"results":[{"user":{"current_user_id":0,"user_id":297508,"user_name":"moderator","rating":120951.5,"gender":0,"comments_count":73338,"stories_count":0,"stories_hot_count":0,"pluses_count":192,"minuses_count":383,"signup_date":1367078405,"is_rating_ban":true,"avatar":"https://cs6.pikabu.ru/avatars/297/x297508-739228827.png","awards":[{"id":249466,"user_id":297508,"award_id":13,"award_title":"5 лет на Пикабу","award_image":"https://cs.pikabu.ru/images/awards/2x/5_let_pikabushnik.png","story_id":0,"story_title":"","date":"2018-04-27 20:00:21","is_hidden":0,"comment_id":0,"link":""},{"id":1746,"user_id":297508,"award_id":8,"award_title":"более 1000 подписчиков","award_image":"https://cs.pikabu.ru/images/awards/2x/1000_subs.png","story_id":0,"story_title":"","date":"2015-07-05 06:00:06","is_hidden":0,"comment_id":0,"link":""}],"is_subscribed":false,"is_ignored":false,"note":"","approved":"Команда Пикабу","communities":[{"name":"Фейки","link":"fakes","avatar":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png","avatar_url":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png"},{"name":"Вопросы по модерации","link":"moderator","avatar":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png","avatar_url":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png"},{"name":"Баяны","link":"reposts","avatar":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png","avatar_url":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png"}],"subscribers_count":5511,"is_user_banned":false,"is_user_fully_banned":false,"public_ban_history":[],"user_ban_time":-62169993079}}]}
{"parsing_timestamp":1545920256,"parser_id":"d3dev/test_proxy_parser_copy_100","number_of_results":1,"results":[{"user":{"current_user_id":0,"user_id":297508,"user_name":"moderator","rating":120951.5,"gender":0,"comments_count":73338,"stories_count":0,"stories_hot_count":0,"pluses_count":192,"minuses_count":383,"signup_date":1367078405,"is_rating_ban":true,"avatar":"https://cs6.pikabu.ru/avatars/297/x297508-739228827.png","awards":[{"id":249466,"user_id":297508,"award_id":13,"award_title":"5 лет на Пикабу","award_image":"https://cs.pikabu.ru/images/awards/2x/5_let_pikabushnik.png","story_id":0,"story_title":"","date":"2018-04-27 20:00:21","is_hidden":0,"comment_id":0,"link":""},{"id":1746,"user_id":297508,"award_id":8,"award_title":"более 1000 подписчиков","award_image":"https://cs.pikabu.ru/images/awards/2x/1000_subs.png","story_id":0,"story_title":"","date":"2015-07-05 06:00:06","is_hidden":0,"comment_id":0,"link":""}],"is_subscribed":false,"is_ignored":false,"note":"","approved":"Команда Пикабу","communities":[{"name":"Фейки","link":"fakes","avatar":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png","avatar_url":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png"},{"name":"Вопросы по модерации","link":"moderator","avatar":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png","avatar_url":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png"},{"name":"Баяны","link":"reposts","avatar":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png","avatar_url":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png"}],"subscribers_count":5511,"is_user_banned":false,"is_user_fully_banned":false,"public_ban_history":[],"user_ban_time":-62169993079}}]}
{"parsing_timestamp":1545920260,"parser_id":"d3dev/test_proxy_parser_copy_40","number_of_results":1,"results":[{"user":{"current_user_id":0,"user_id":297508,"user_name":"moderator","rating":120951.5,"gender":0,"comments_count":73338,"stories_count":0,"stories_hot_count":0,"pluses_count":192,"minuses_count":383,"signup_date":1367078405,"is_rating_ban":true,"avatar":"https://cs6.pikabu.ru/avatars/297/x297508-739228827.png","awards":[{"id":249466,"user_id":297508,"award_id":13,"award_title":"5 лет на Пикабу","award_image":"https://cs.pikabu.ru/images/awards/2x/5_let_pikabushnik.png","story_id":0,"story_title":"","date":"2018-04-27 20:00:21","is_hidden":0,"comment_id":0,"link":""},{"id":1746,"user_id":297508,"award_id":8,"award_title":"более 1000 подписчиков","award_image":"https://cs.pikabu.ru/images/awards/2x/1000_subs.png","story_id":0,"story_title":"","date":"2015-07-05 06:00:06","is_hidden":0,"comment_id":0,"link":""}],"is_subscribed":false,"is_ignored":false,"note":"","approved":"Команда Пикабу","communities":[{"name":"Фейки","link":"fakes","avatar":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png","avatar_url":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png"},{"name":"Вопросы по модерации","link":"moderator","avatar":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png","avatar_url":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png"},{"name":"Баяны","link":"reposts","avatar":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png","avatar_url":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png"}],"subscribers_count":5511,"is_user_banned":false,"is_user_fully_banned":false,"public_ban_history":[],"user_ban_time":-62169993079}}]}
{"parsing_timestamp":1545920260,"parser_id":"d3dev/test_proxy_parser_copy_115","number_of_results":1,"results":[{"user":{"current_user_id":0,"user_id":297508,"user_name":"moderator","rating":120951.5,"gender":0,"comments_count":73338,"stories_count":0,"stories_hot_count":0,"pluses_count":192,"minuses_count":383,"signup_date":1367078405,"is_rating_ban":true,"avatar":"https://cs6.pikabu.ru/avatars/297/x297508-739228827.png","awards":[{"id":249466,"user_id":297508,"award_id":13,"award_title":"5 лет на Пикабу","award_image":"https://cs.pikabu.ru/images/awards/2x/5_let_pikabushnik.png","story_id":0,"story_title":"","date":"2018-04-27 20:00:21","is_hidden":0,"comment_id":0,"link":""},{"id":1746,"user_id":297508,"award_id":8,"award_title":"более 1000 подписчиков","award_image":"https://cs.pikabu.ru/images/awards/2x/1000_subs.png","story_id":0,"story_title":"","date":"2015-07-05 06:00:06","is_hidden":0,"comment_id":0,"link":""}],"is_subscribed":false,"is_ignored":false,"note":"","approved":"Команда Пикабу","communities":[{"name":"Фейки","link":"fakes","avatar":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png","avatar_url":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png"},{"name":"Вопросы по модерации","link":"moderator","avatar":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png","avatar_url":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png"},{"name":"Баяны","link":"reposts","avatar":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png","avatar_url":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png"}],"subscribers_count":5511,"is_user_banned":false,"is_user_fully_banned":false,"public_ban_history":[],"user_ban_time":-62169993079}}]}
{"parsing_timestamp":1545920261,"parser_id":"d3dev/test_proxy_parser_copy_70","number_of_results":1,"results":[{"user":{"current_user_id":0,"user_id":297508,"user_name":"moderator","rating":120951.5,"gender":0,"comments_count":73338,"stories_count":0,"stories_hot_count":0,"pluses_count":192,"minuses_count":383,"signup_date":1367078405,"is_rating_ban":true,"avatar":"https://cs6.pikabu.ru/avatars/297/x297508-739228827.png","awards":[{"id":249466,"user_id":297508,"award_id":13,"award_title":"5 лет на Пикабу","award_image":"https://cs.pikabu.ru/images/awards/2x/5_let_pikabushnik.png","story_id":0,"story_title":"","date":"2018-04-27 20:00:21","is_hidden":0,"comment_id":0,"link":""},{"id":1746,"user_id":297508,"award_id":8,"award_title":"более 1000 подписчиков","award_image":"https://cs.pikabu.ru/images/awards/2x/1000_subs.png","story_id":0,"story_title":"","date":"2015-07-05 06:00:06","is_hidden":0,"comment_id":0,"link":""}],"is_subscribed":false,"is_ignored":false,"note":"","approved":"Команда Пикабу","communities":[{"name":"Фейки","link":"fakes","avatar":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png","avatar_url":"https://cs8.pikabu.ru/images/community/1/1511452065290441807.png"},{"name":"Вопросы по модерации","link":"moderator","avatar":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png","avatar_url":"https://cs10.pikabu.ru/images/community/298/1540135481227521303.png"},{"name":"Баяны","link":"reposts","avatar":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png","avatar_url":"https://cs8.pikabu.ru/images/community/298/1478189420255384259.png"}],"subscribers_count":5511,"is_user_banned":false,"is_user_fully_banned":false,"public_ban_history":[],"user_ban_time":-62169993079}}]}
`
	for _, item := range strings.Split(data, "\n") {
		item = strings.TrimSpace(item)
		if len(item) == 0 {
			continue
		}
		err = pushTaskToQueue("user_profile", []byte(item))
		if err != nil {
			handleError(err)
		}
	}

	waitForQueueEmpty()
	time.Sleep(5 * time.Second)

	user := &models.PikabuUser{
		PikabuId: 297508,
	}
	err = models.Db.Select(user)
	if err != nil {
		handleError(err)
	}
	assert.Equal(t, "moderator", user.Username)
	assert.Equal(t, int32(120951), user.Rating)
	assert.Equal(t, "0", user.Gender)
	assert.Equal(t, int32(73338), user.NumberOfComments)
	assert.Equal(t, int32(0), user.NumberOfStories)
	assert.Equal(t, int32(0), user.NumberOfHotStories)
	assert.Equal(t, int32(192), user.NumberOfPluses)
	assert.Equal(t, int32(383), user.NumberOfMinuses)
	assert.Equal(t, models.TimestampType(1367078405), user.SignupTimestamp)
	assert.Equal(t, true, user.IsRatingHidden)
	assert.Equal(t, "https://cs6.pikabu.ru/avatars/297/x297508-739228827.png", user.AvatarURL)
	// assert.Equal(t, true, user.AwardIds)
	assert.Equal(t, "Команда Пикабу", user.ApprovedText)
	assert.Equal(t, int32(5511), user.NumberOfSubscribers)
	assert.Equal(t, false, user.IsBanned)
	assert.Equal(t, false, user.IsPermanentlyBanned)
	assert.Equal(t, models.TimestampType(-62169993079), user.BanEndTimestamp)
	assert.Equal(t, models.TimestampType(1545920253), user.AddedTimestamp)
	assert.Equal(t, models.TimestampType(1545920261), user.LastUpdateTimestamp)

	ratingVersions := []models.PikabuUserRatingVersion{}
	err = models.Db.Model(&ratingVersions).
		Where("item_id = ?", user.PikabuId).
		Order("timestamp").
		Select()
	handleError(err)

	assert.Equal(t, []models.PikabuUserRatingVersion{}, ratingVersions)

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
}
