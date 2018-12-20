package results_processor

import (
	"bitbucket.org/d3dev/parse_pikabu/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"gogsweb.2-47.ru/d3dev/pikago"
	"sync"
)

var processCommunitiesPagesMutex = &sync.Mutex{}

func processCommunitiesPages(parsingTimestamp models.TimestampType, communitiesPages []pikago.CommunitiesPage) error {
	processCommunitiesPagesMutex.Lock()
	defer processCommunitiesPagesMutex.Unlock()

	tx, err := models.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// TODO: process
	logger.Log.Debugf("processing community pages. number of pages is %v", len(communitiesPages))

	return tx.Commit()
}
