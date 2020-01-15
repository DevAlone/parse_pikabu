package telegram

import (
	"fmt"
	"time"

	"github.com/DevAlone/parse_pikabu/core/config"
	"github.com/DevAlone/parse_pikabu/core/logger"
	"github.com/DevAlone/parse_pikabu/models"
	"github.com/ansel1/merry"
)

func handlePikabuCommentCreatedEvent(
	data *models.PikabuComment,
	eventTime models.TimestampType,
) error {
	if data.IsDeleted {
		logger.Log.Debugf("processing new deleted comment %+v\n", data)
		messages := commentToMessages(data, config.Settings.Pikabu18BotDeletedAtFirstParsingChat)
		for _, message := range messages {
			recipient, err := bot.ChatByID(config.Settings.Pikabu18BotDeletedAtFirstParsingChat)
			if err != nil {
				return merry.Wrap(err)
			}
			_, err = bot.Send(recipient, message)
			if err != nil {
				return merry.Wrap(err)
			}
		}
	}
	return nil
}

func handlePikabuCommentChangedEvent(
	prevState *models.PikabuComment,
	currState *models.PikabuComment,
	eventTime models.TimestampType,
) error {
	if prevState.IsDeleted != currState.IsDeleted {
		messages := createCommentsChangedTgMessage(prevState, currState, config.Settings.Pikabu18BotDeletedChat)
		for _, message := range messages {
			recipient, err := bot.ChatByID(config.Settings.Pikabu18BotDeletedChat)
			if err != nil {
				return merry.Wrap(err)
			}
			_, err = bot.Send(recipient, message)
			if err != nil {
				return merry.Wrap(err)
			}
		}
	}
	return nil
}

func createCommentsChangedTgMessage(
	prevState *models.PikabuComment,
	currState *models.PikabuComment,
	channelName string,
) []interface{} {
	result := []interface{}{}
	result = append(
		result,
		commentToMessages(prevState, channelName)...,
	)
	result = append(
		result,
		commentToMessages(currState, channelName)...,
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
