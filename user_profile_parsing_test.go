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
`,
	))
	if err != nil {
		panic(err)
	}

	// TODO: wait for queue to be empty
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
