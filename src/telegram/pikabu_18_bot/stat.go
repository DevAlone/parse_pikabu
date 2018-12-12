package pikabu_18_bot

import (
	. "models"
	"strconv"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (this *Pikabu18Bot) processGetStat(message *tb.Message) {
	//text := "Количество комментариев: "
	var result struct {
		Count uint64
		MinId uint64
		MaxId uint64
	}
	_, err := Db.QueryOne(&result, `
		SELECT COUNT(*) as count, MIN(id) as min_id, MAX(id) as max_id FROM comments;
	`)
	if err != nil {
		this.Bot.Send(message.Sender, "что-то плохое случилось")
		Log.Error(err)
	}

	text := "Количество комментариев: " + strconv.FormatUint(result.Count, 10) + "\n"
	text += "Первый комментарий: " + strconv.FormatUint(result.MinId, 10) + " /first_comment\n"
	text += "Последний комментарий: " + strconv.FormatUint(result.MaxId, 10) + " /last_comment\n"

	_, err = this.Bot.Send(message.Sender, text)
	if err != nil {
		Log.Error(err)
	}
}
