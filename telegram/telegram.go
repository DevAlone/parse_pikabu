package telegram

import (
	"regexp"
	"strings"

	"github.com/DevAlone/parse_pikabu/core/config"
	"github.com/DevAlone/parse_pikabu/modelhooks"
	"github.com/ansel1/merry"

	tb "gopkg.in/tucnak/telebot.v2"
)

var bot *tb.Bot

// RunTelegramNotifier - call to run telegram notifier
func RunTelegramNotifier() error {
	var err error
	bot, err = tb.NewBot(tb.Settings{
		Token: config.Settings.Pikabu18BotToken,
		// Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		return merry.Wrap(err)
	}

	modelhooks.SubscribeToPikabuCommentCreatedEvent(handlePikabuCommentCreatedEvent)
	modelhooks.SubscribeToPikabuCommentChangedEvent(handlePikabuCommentChangedEvent)

	modelhooks.SubscribeToPikabuUserCreatedEvent(handlePikabuUserCreatedEvent)
	modelhooks.SubscribeToPikabuUserChangedEvent(handlePikabuUserChangedEvent)

	// TODO: implement:
	// modelhooks.SubscribeToPikabuStoryCreatedEvent(handlePikabuStoryCreatedEvent)
	// modelhooks.SubscribeToPikabuStoryChangedEvent(handlePikabuStoryChangedEvent)

	return nil
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
