// consider using https://github.com/tucnak/telebot

package pikabu_18_bot

import (
	. "config"
	//"encoding/json"
	//"errors"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/op/go-logging"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Pikabu18BotRequestError struct {
	Code int
	Text string
}

func (this *Pikabu18BotRequestError) Error() string {
	return strconv.FormatInt(int64(this.Code), 10) + this.Text
}

func NewError(code int, text string) error {
	return &Pikabu18BotRequestError{
		Code: code,
		Text: text,
	}
}

type Pikabu18Bot struct {
	Bot *tb.Bot
}

var Bot *Pikabu18Bot

var Log = logging.MustGetLogger("parse_pikabu/telegram/pikabu_18_bot")
var logFormat = logging.MustStringFormatter(
	// `%{color}%{level:.5s} %{time:15:04:05.000} %{shortfunc} ▶ %{id:03x}%{color:reset} %{message}`,
	`%{color}%{level:.5s} %{time:2006-01-02T15:04:05.999Z-07:00} %{shortfunc} ▶ %{id:03x}%{color:reset} %{message}`,
)

func RunPikabu18Bot() {
	file, err := os.OpenFile("logs/telegram__pikabu_18_bot.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	loggingBackend := logging.NewLogBackend(file, "", 0)
	loggingBackendFormatter := logging.NewBackendFormatter(loggingBackend, logFormat)

	logging.SetBackend(loggingBackend, loggingBackendFormatter)
	poller := &tb.LongPoller{Timeout: 10 * time.Second}
	middleware := tb.NewMiddlewarePoller(poller, func(update *tb.Update) bool {
		jsonBytes, err := json.Marshal(update)
		if err != nil {
			Log.Error(err)
		} else {
			Log.Info(":raw_request ->", string(jsonBytes))
		}

		if update.EditedMessage != nil {
			update.Message = update.EditedMessage
		} else if update.Callback != nil {
			//if strings.HasPrefix(update.Callback.Data, "/comment") {
			update.Message = &tb.Message{
				Sender: update.Callback.Sender,
				Text:   update.Callback.Data,
			}
			update.Callback = nil
			//}
		}

		if update.Message != nil {
			Bot.Bot.Notify(update.Message.Sender, tb.Typing)
			Log.Info(
				":request ->",
				update.Message.Sender.ID,
				"|",
				update.Message.Sender.Username,
				"|",
				update.Message.Text,
			)
		}

		return true
	})

	bot, err := tb.NewBot(tb.Settings{
		Token:  Settings.Pikabu18BotToken,
		Poller: middleware,
	})
	Bot = &Pikabu18Bot{
		Bot: bot,
	}

	if err != nil {
		log.Fatal(err.Error() + ". Did you forget token?")
		panic(err)
	}

	bot.Handle("/start", func(m *tb.Message) {
		Bot.processHelp(m)
	})
	bot.Handle("/help", func(m *tb.Message) {
		Bot.processHelp(m)
	})
	bot.Handle("/comment", func(m *tb.Message) {
		m.Payload = strings.TrimSpace(m.Payload)
		if strings.Contains(m.Payload, "'") {
			// sql injection :)
			m.Payload = strings.ToLower(m.Payload)
			if strings.Contains(m.Payload, "drop database") {
				bot.Send(m.Sender, "Deleting database...")
				time.Sleep(3 * time.Second)
				bot.Send(m.Sender, "Hacking attempt found!")
				time.Sleep(1 * time.Second)
				bot.Send(m.Sender, "Determining user's IP...")
				time.Sleep(3 * time.Second)
				bot.Send(m.Sender, "Calling FSB...")
				return
			}
			_, err = bot.Send(m.Sender, `ERROR:  unterminated quoted string at or near "'`+m.Payload+`';"
LINE 1: select * from comments where id = '`+m.Payload+"';")
			if err != nil {
				Log.Error(err)
			}
			return
		}

		id, err := parseNumber(m.Payload)
		if err != nil || id < 0 {
			_, err := bot.Send(
				m.Sender,
				"плохой запрос, пример хорошего - `/comment 113`",
				tb.ModeMarkdown,
			)
			if err != nil {
				Log.Error(err)
			}
			return
		}
		println(m.Payload, id)

		Bot.processGetComment(m, uint64(id))
	})
	bot.Handle("/deleted", func(m *tb.Message) {
		Bot.processGetDeleted(m, false)
	})
	bot.Handle("/porn", func(m *tb.Message) {
		Bot.processGetDeleted(m, true)
	})
	bot.Handle("/porno", func(m *tb.Message) {
		Bot.processGetDeleted(m, true)
	})
	bot.Handle("/first_comment", func(m *tb.Message) {
		if strings.ToLower(m.Sender.Username) == "devalone" {
			Bot.processGetComment(m, uint64(0))
		}
	})
	bot.Handle("/last_comment", func(m *tb.Message) {
		if strings.ToLower(m.Sender.Username) == "devalone" {
			Bot.processGetComment(m, uint64(9223372036854775807))
		}
	})
	bot.Handle("/user", func(m *tb.Message) {
		Bot.Bot.Send(m.Sender, "пока недоступно")
	})
	bot.Handle("/story", func(m *tb.Message) {
		Bot.Bot.Send(m.Sender, "пока недоступно")
	})
	bot.Handle("/community", func(m *tb.Message) {
		Bot.Bot.Send(m.Sender, "пока недоступно")
	})
	bot.Handle("/stat", func(m *tb.Message) {
		if strings.ToLower(m.Sender.Username) == "devalone" {
			Bot.processGetStat(m)
		}
	})
	bot.Handle(tb.OnText, func(m *tb.Message) {
		Bot.processCommonRequest(m)
	})

	bot.Start()

	return
}

func (this *Pikabu18Bot) processHelp(message *tb.Message) {
	_, err := this.Bot.Send(message.Sender, `
Просто введи id комментария или ссылку на него.
Либо используй команду /comment.
Примеры валидных запросов:

`+"```"+`
pikabu.ru/story/_59?cid=113
/comment 112248140
/comment id=112248140
`+"```"+`

Если бот вас обидел, или вы нашли баг, или хотите предложить улучшение функционала, или вам просто не с кем поговорить, пишите сюда -> pikabu18bot@gmail.com
		`, tb.ModeMarkdown)
	if err != nil {
		Log.Error(err)
	}
}

func parseNumber(text string) (int64, error) {
	numberRegex := regexp.MustCompile(`([0-9]+)[^0-9]*?$`)
	values := numberRegex.FindStringSubmatch(text)
	if values == nil {
		return 0, errors.New("regex error")
	}
	strId := values[1]

	return strconv.ParseInt(strId, 10, 64)
}

func (this *Pikabu18Bot) SendFile(recipient *tb.User, data []byte) error {
	tmpFileName := tempFileName("comment_", ".txt")
	defer os.Remove(tmpFileName)

	err := ioutil.WriteFile(tmpFileName, data, 0644)
	if err != nil {
		return err
	}
	_, err = this.Bot.Send(recipient, &tb.Document{File: tb.FromDisk(tmpFileName)})

	return err
}

func tempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}

func (this *Pikabu18Bot) processCommonRequest(message *tb.Message) {
	s := message.Text
	s = strings.ToLower(s)

	urlRegex := regexp.MustCompile(
		`.*?pikabu.ru/(@(?P<Username>[a-zA-Z0-9_.-]+))|(story/.*?_(?P<StoryId>[0-9]+)(\?[^0-9]+(?P<CommentId>[0-9]+))?).*`)
	values := urlRegex.FindStringSubmatch(s)

	if values != nil {
		keys :=
			urlRegex.SubexpNames()
		matches :=
			map[string]string{}

		for i := 0; i < len(keys); i++ {
			matches[keys[i]] = values[i]

		}

		/*
			if username, ok := matches["Username"]; ok && len(username) > 0 {
				return "/get_user username=" + username, nil
			} else if commentId, ok := matches["CommentId"]; ok && len(commentId) > 0 {
				return "/get_comment id=" + commentId, nil
			} else if storyId, ok := matches["StoryId"]; ok && len(storyId) > 0 {
				return "/get_story id=" + storyId, nil
			}
		*/

		if commentId, ok := matches["CommentId"]; ok && len(commentId) > 0 {
			id, err := strconv.ParseUint(commentId, 10, 64)
			if err == nil {
				this.processGetComment(message, id)
				return
			}
		}
	}

	numberRegex := regexp.MustCompile(`([0-9]+)[^0-9]*?$`)
	values = numberRegex.FindStringSubmatch(s)

	if values != nil {
		id, err := strconv.ParseUint(values[1], 10, 64)
		if err == nil {
			this.processGetComment(message, id)
			return
		}
	}

	this.processHelp(message)
}
