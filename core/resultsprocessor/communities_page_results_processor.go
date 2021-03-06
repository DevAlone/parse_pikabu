package resultsprocessor

import (
	"sync"

	"github.com/DevAlone/parse_pikabu/core/logger"
	"github.com/DevAlone/parse_pikabu/modelhooks"
	"github.com/DevAlone/parse_pikabu/models"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg"
	pikago_models "gogsweb.2-47.ru/d3dev/pikago/models"
)

var processCommunitiesPagesMutex = &sync.Mutex{}

func processCommunitiesPages(parsingTimestamp models.TimestampType, communitiesPages []pikago_models.CommunitiesPage) error {
	processCommunitiesPagesMutex.Lock()
	defer processCommunitiesPagesMutex.Unlock()

	_, err := models.Db.Model().Exec(`
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
			err := processCommunity(parsingTimestamp, &community)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func processCommunity(
	parsingTimestamp models.TimestampType,
	parsedCommunity *pikago_models.Community,
) error {
	moderatorIds := []uint64{}
	for _, moderator := range parsedCommunity.CommunityModerators {
		moderatorIds = append(moderatorIds, moderator.ID.Value)
	}

	newCommunity := &models.PikabuCommunity{
		PikabuID:            parsedCommunity.ID.Value,
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
		AdminId:             parsedCommunity.CommunityAdmin.ID.Value,
		ModeratorIds:        moderatorIds,
		AddedTimestamp:      parsingTimestamp,
		LastUpdateTimestamp: parsingTimestamp,
	}

	community := &models.PikabuCommunity{
		PikabuID: newCommunity.PikabuID,
	}

	err := models.Db.Select(community)
	if err == pg.ErrNoRows {
		modelhooks.HandleModelCreated(*newCommunity, parsingTimestamp)
		err := models.Db.Insert(newCommunity)
		if err != nil {
			return errors.New(err)
		}
		return nil
	} else if err != nil {
		return errors.New(err)
	}

	modelhooks.HandleModelChanged(*community, *newCommunity, parsingTimestamp)

	if parsingTimestamp <= community.LastUpdateTimestamp {
		// TODO: find a better way
		logger.Log.Warningf("skipping community %v because of old parsing result", community.LinkName)
		return nil
	}

	_, err = processModelFieldsVersions(nil, community, newCommunity, parsingTimestamp)
	if err != nil {
		return err
	}

	community.LastUpdateTimestamp = parsingTimestamp

	err = models.Db.Update(community)
	if err != nil {
		return errors.New(err)
	}

	return nil
}
