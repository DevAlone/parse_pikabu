package parser

import (
	"strings"

	"github.com/DevAlone/parse_pikabu/models"
	"github.com/ansel1/merry"
	"gogsweb.2-47.ru/d3dev/pikago"
	pikago_models "gogsweb.2-47.ru/d3dev/pikago/models"
)

func (p *Parser) processParseStoryTask(task *models.ParseStoryTask) error {
	results := []pikago_models.StoryGetResult{}

	for page := uint(0); page < 99; page++ {
		res, err := p.pikagoClient.StoryGet(task.PikabuID, page, 0)
		if pe, ok := err.(*pikago.PikabuErrorRequestedPageNotFound); ok {
			return p.PutResultsToQueue("story_not_found", []models.ParserStoryNotFoundResultData{
				{
					PikabuID:    task.PikabuID,
					PikabuError: pe,
				},
			})
		}
		if pe, ok := err.(*pikago.PikabuError); ok {
			// Пост содержит материалы, запрещенные политикой магазина приложений, вы можете посмотреть его в веб-версии Пикабу
			if strings.Contains(pe.Message, "Google Play Market") ||
				strings.Contains(pe.Message, `\u0437\u0430\u043f\u0440\u0435\u0449\u0435\u043d\u043d\u044b\u0435 \u043f\u043e\u043b\u0438\u0442\u0438\u043a\u043e\u0439 \u043c\u0430\u0433\u0430\u0437\u0438\u043d\u0430 \u043f\u0440\u0438\u043b\u043e\u0436\u0435\u043d\u0438\u0439`) ||
				strings.Contains(pe.Message, "запрещенные политикой магазина приложений") {
				// consider as deleted
				return p.PutResultsToQueue("story_not_found", []models.ParserStoryNotFoundResultData{
					{
						PikabuID:    task.PikabuID,
						PikabuError: pe,
					},
				})
			}
		}

		if err != nil {
			return err
		}

		if res == nil {
			return merry.Errorf("processParseStoryTask(): res is nil. Task is %+v", task)
		}

		results = append(results, *res)
		if !res.HasNextCommentsPage {
			break
		}
	}

	return p.PutResultsToQueue("story_get_result", results)
}
