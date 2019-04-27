package parser

import (
	"time"

	pikago_models "gogsweb.2-47.ru/d3dev/pikago/models"
)

func (p *Parser) processParseCommunitiesPagesTask() error {
	results := []pikago_models.CommunitiesPage{}

	page := 0
	for {
		communitiesPage, err := p.pikagoClient.CommunitiesGet(page)
		if err != nil {
			return err
		}
		if len(communitiesPage.List) == 0 {
			break
		}
		results = append(results, *communitiesPage)

		page++
		time.Sleep(time.Duration(p.Config.PikagoWaitBetweenProcessingPages) * time.Second)
	}

	return p.PutResultsToQueue("communities_pages", results)
}
