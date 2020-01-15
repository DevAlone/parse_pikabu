package telegram

import (
	"fmt"
	"strings"
	"time"

	"github.com/DevAlone/parse_pikabu/core/config"
	"github.com/DevAlone/parse_pikabu/core/logger"
	"github.com/DevAlone/parse_pikabu/models"
	"github.com/ansel1/merry"
)

func handlePikabuUserCreatedEvent(
	data *models.PikabuUser,
	eventTime models.TimestampType,
) error {
	return nil
}

func handlePikabuUserChangedEvent(
	prevState *models.PikabuUser,
	currState *models.PikabuUser,
	eventTime models.TimestampType,
) error {
	if prevState.IsDeleted != currState.IsDeleted || prevState.IsPermanentlyBanned != currState.IsPermanentlyBanned {
		logger.Log.Debugf("processing deleted user %+v %+v\n", prevState, currState)
		messages := createUserChangedTgMessage(prevState, currState, config.Settings.Pikabu18BotDeletedUsersChat)
		for _, message := range messages {
			recipient, err := bot.ChatByID(config.Settings.Pikabu18BotDeletedUsersChat)
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

func createUserChangedTgMessage(
	prevState *models.PikabuUser,
	currState *models.PikabuUser,
	channelName string,
) []interface{} {
	result := []interface{}{}
	result = append(
		result,
		userToMessages(prevState, channelName)...,
	)
	result = append(
		result,
		userToMessages(currState, channelName)...,
	)

	return result
}

func userToMessages(
	data *models.PikabuUser,
	channelName string,
) []interface{} {
	result := []interface{}{}

	text := "ID на Пикабу: " + fmt.Sprint(data.PikabuID) + "\n"
	text += "Никнейм: " + data.Username + "\n"
	text += "Дата Регистрации: " + fmt.Sprint(time.Unix(int64(data.SignupTimestamp), 0)) + "\n"
	text += "Пол: " + data.Gender + "\n"
	text += "Рейтинг: " + fmt.Sprint(data.Rating) + "\n"
	text += "Количество комментариев: " + fmt.Sprint(data.NumberOfComments) + "\n"
	text += "Количество подписчиков: " + fmt.Sprint(data.NumberOfSubscribers) + "\n"
	text += "Количество постов: " + fmt.Sprint(data.NumberOfStories) + "\n"
	text += "Количество горячих постов: " + fmt.Sprint(data.NumberOfHotStories) + "\n"
	text += "Количество плюсов: " + fmt.Sprint(data.NumberOfPluses) + "\n"
	text += "Количество минусов: " + fmt.Sprint(data.NumberOfMinuses) + "\n"
	text += "Ссылка: https://pikabu.ru/@" + data.Username + "\n"
	text += "Удалён: " + fmt.Sprint(data.IsDeleted) + "\n"
	text += "Забанен: " + fmt.Sprint(data.IsPermanentlyBanned) + "\n"

	result = append(
		result,
		text,
	)

	if len(strings.TrimSpace(data.AvatarURL)) > 0 {
		result = append(
			result,
			URLToTelebotSendable(data.AvatarURL),
		)
	}

	return result
}
