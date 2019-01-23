package results_processor

import (
	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg"
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

	_, err = tx.Model().Exec(`
		UPDATE simple_tasks
		SET is_done = true
		WHERE name = 'parse_communities'
	`)
	if err != nil {
		return err
	}

	// TODO: process
	// TODO: process deleted communities
	logger.Log.Debugf("processing community pages. number of pages is %v", len(communitiesPages))
	for _, communitiesPage := range communitiesPages {
		for _, community := range communitiesPage.List {
			err := processCommunity(parsingTimestamp, tx, &community)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func processCommunity(
	parsingTimestamp models.TimestampType,
	tx *pg.Tx, parsedCommunity *pikago.Community,
) error {
	moderatorIds := []uint64{}
	for _, moderator := range parsedCommunity.CommunityModerators {
		moderatorIds = append(moderatorIds, moderator.Id.Value)
	}

	newCommunity := &models.PikabuCommunity{
		PikabuId:            parsedCommunity.Id.Value,
		Name:                parsedCommunity.Name,
		LinkName:            parsedCommunity.Link,
		URL:                 parsedCommunity.URL,
		AvatarURL:           parsedCommunity.AvatarURL,
		BackgroundImageURL:  parsedCommunity.BackgroundImageURL,
		Tags:                parsedCommunity.Tags,
		NumberOfStories:     int32(parsedCommunity.NumberOfStories.Value),
		NumberOfSubscribers: int32(parsedCommunity.NumberOfSubscribers.Value),
		Description:         parsedCommunity.Description,
		Rules:               parsedCommunity.Rules,
		Restrictions:        parsedCommunity.Restrictions,
		AdminId:             parsedCommunity.CommunityAdmin.Id.Value,
		ModeratorIds:        moderatorIds,
		AddedTimestamp:      parsingTimestamp,
		LastUpdateTimestamp: parsingTimestamp,
	}

	community := &models.PikabuCommunity{
		PikabuId: newCommunity.PikabuId,
	}

	err := tx.Select(community)
	if err == pg.ErrNoRows {
		err := tx.Insert(newCommunity)
		if err != nil {
			return errors.New(err)
		}
		return nil
	} else if err != nil {
		return errors.New(err)
	}

	if parsingTimestamp <= community.LastUpdateTimestamp {
		// TODO: find a better way
		logger.Log.Warning("skipping community %v because of old parsing result", community.LinkName)
		return nil
	}

	_, err = processModelFieldsVersions(tx, community, newCommunity, parsingTimestamp)
	if err != nil {
		return err
	}

	community.LastUpdateTimestamp = parsingTimestamp

	err = tx.Update(community)
	if err != nil {
		return errors.New(err)
	}

	return nil
}
