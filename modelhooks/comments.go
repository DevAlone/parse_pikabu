package modelhooks

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/DevAlone/parse_pikabu/core/config"
	"github.com/DevAlone/parse_pikabu/core/logger"
	"github.com/DevAlone/parse_pikabu/models"
	"github.com/go-errors/errors"

	tb "gopkg.in/tucnak/telebot.v2"
)

// TODO: rewrite with a new framework

type commentModelChange struct {
	PrevState  models.PikabuComment
	CurrState  models.PikabuComment
	ChangeTime models.TimestampType
}

var commentModelChanges = make(chan *commentModelChange)

type commentModelCreate struct {
	Data       models.PikabuComment
	CreateTime models.TimestampType
}

var commentModelCreates = make(chan *commentModelCreate)

// RunTelegramNotifier - call to run telegram notifier
func RunTelegramNotifier() error {
	// bot, err := tgbotapi.NewBotAPI(config.Settings.Pikabu18BotToken)
	bot, err := tb.NewBot(tb.Settings{
		Token: config.Settings.Pikabu18BotToken,
		// Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		return errors.New(err)
	}

	for {
		select {
		case commentChange := <-commentModelChanges:
			if commentChange.PrevState.IsDeleted != commentChange.CurrState.IsDeleted {
				messages := createCommentsChangedTgMessage(commentChange, config.Settings.Pikabu18BotDeletedChat)
				for _, message := range messages {
					recipient, err := bot.ChatByID(config.Settings.Pikabu18BotDeletedChat)
					if err != nil {
						return errors.New(err)
					}
					_, err = bot.Send(recipient, message)
					if err != nil {
						return errors.New(err)
					}
				}
			}
		case commentCreate := <-commentModelCreates:
			if commentCreate.Data.IsDeleted {
				logger.Log.Debugf("processing new deleted comment %+v\n", commentCreate.Data)
				messages := commentToMessages(&commentCreate.Data, config.Settings.Pikabu18BotDeletedAtFirstParsingChat)
				for _, message := range messages {
					recipient, err := bot.ChatByID(config.Settings.Pikabu18BotDeletedAtFirstParsingChat)
					if err != nil {
						return errors.New(err)
					}
					_, err = bot.Send(recipient, message)
					if err != nil {
						return errors.New(err)
					}
				}
			}
		}
	}
}

func createCommentsChangedTgMessage(
	commentChange *commentModelChange,
	channelName string,
) []interface{} {
	result := []interface{}{}
	result = append(
		result,
		commentToMessages(&commentChange.PrevState, channelName)...,
	)
	result = append(
		result,
		commentToMessages(&commentChange.CurrState, channelName)...,
	)

	return result
}

func commentToMessages(
	comment *models.PikabuComment,
	channelName string,
) []interface{} {
	result := []interface{}{}

	text := "Дата: " + fmt.Sprint(time.Unix(int64(comment.CreatedAtTimestamp), 0)) + "\n"
	text += "Автор: " + comment.AuthorUsername + "\n"
	text += "Текст: " + comment.Text + "\n"
	text += "Ссылка: https://pikabu.ru/story/_" + fmt.Sprint(comment.StoryID) + "?cid=" + fmt.Sprint(comment.PikabuID) + "\n"
	text += "Удалён: " + fmt.Sprint(comment.IsDeleted) + "\n"

	result = append(
		result,
		text,
	)

	for _, image := range comment.Images {
		if len(image.LargeURL) != 0 {
			result = append(
				result,
				URLToTelebotSendable(image.LargeURL),
			)
		} else if len(image.SmallURL) != 0 {
			result = append(
				result,
				URLToTelebotSendable(image.SmallURL),
			)
		}
		if len(image.AnimationBaseURL) != 0 {
			// TODO: iterate through AnimationFormats
			result = append(
				result,
				URLToTelebotSendable(image.AnimationBaseURL+".mp4"),
			)
		}
	}

	return result
}

// URLToTelebotSendable -
func URLToTelebotSendable(url string) interface{} {
	if strings.HasSuffix(url, ".gif") {
		return &tb.Document{
			File: tb.FromURL(url),
		}
	}
	if match, _ := regexp.MatchString(`(\.png|\.jpg|\.jpeg)$`, strings.ToLower(url)); match {
		return &tb.Photo{
			File: tb.FromURL(url),
		}
	}
	if match, _ := regexp.MatchString(`(\.mp4|\.mpeg|\.avi|\.webm)$`, strings.ToLower(url)); match {
		return &tb.Video{
			File: tb.FromURL(url),
		}
	}

	return url
}
