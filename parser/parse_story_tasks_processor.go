package parser

import (
	"strings"

	"bitbucket.org/d3dev/parse_pikabu/models"
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
			if strings.Contains(pe.Message, "Google Play Market") {
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
		results = append(results, *res)
		if !res.HasNextCommentsPage {
			break
		}
	}

	return p.PutResultsToQueue("story_get_result", results)
}
