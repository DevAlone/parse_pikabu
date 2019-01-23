package main

import (
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"context"
	"fmt"
	"sync"

	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/old_models"
	"github.com/go-pg/pg"
	"golang.org/x/sync/semaphore"
)

func createCommunitiesIndices() {
	fmt.Println("creating communities indices...")
	printTimeSinceStart()

	processExec := func(_ interface{}, err error) {
		helpers.PanicOnError(err)
	}

	processExec(oldDb.Exec(`
		CREATE INDEX IF NOT EXISTS communities_app_community_url_name 
		ON communities_app_community (url_name);
	`))

	processExec(oldDb.Exec(`
		CREATE INDEX IF NOT EXISTS communities_app_communitycountersentry_timestamp 
		ON communities_app_communitycountersentry (timestamp);
	`))
	processExec(oldDb.Exec(`
		CREATE INDEX IF NOT EXISTS communities_app_communitycountersentry_community_id
		ON communities_app_communitycountersentry (community_id);
	`))
}

func processCommunities() {
	createCommunitiesIndices()

	fmt.Println("processing communities")

	var communities []models.PikabuCommunity

	err := models.Db.Model(&communities).Select()
	helpers.PanicOnError(err)

	var (
		maxWorkers = 64 // 128
		sem        = semaphore.NewWeighted(int64(maxWorkers))
	)
	ctx := context.TODO()
	var wg sync.WaitGroup

	for _, community := range communities {
		helpers.PanicOnError(sem.Acquire(ctx, 1))
		wg.Add(1)
		go func(community models.PikabuCommunity) {
			defer sem.Release(1)
			processCommunity(community)
			wg.Done()
		}(community)
	}

	wg.Wait()
}

func processCommunity(community models.PikabuCommunity) {
	// fmt.Printf("processing community with id %v\n", community.PikabuId)
	var oldCommunity old_models.Community
	err := oldDb.Model(&oldCommunity).
		Where("url_name ilike ?", community.LinkName).
		Select()
	if err == pg.ErrNoRows {
		return
	}
	helpers.PanicOnError(err)

	var counterEntries []old_models.CommunityCountersEntry
	err = oldDb.Model(&counterEntries).
		Where("community_id = ?", oldCommunity.Id).
		Order("timestamp").
		Select()
	helpers.PanicOnError(err)

	for _, counterEntry := range counterEntries {
		numOfSubsVersion := models.PikabuCommunityNumberOfSubscribersVersion{
			Timestamp: models.TimestampType(counterEntry.Timestamp),
			ItemId:    community.PikabuId,
			Value:     counterEntry.SubscribersCount,
		}
		err := models.Db.Insert(&numOfSubsVersion)
		helpers.PanicOnError(err)

		numOfStoriesVersion := models.PikabuCommunityNumberOfStoriesVersion{
			Timestamp: models.TimestampType(counterEntry.Timestamp),
			ItemId:    community.PikabuId,
			Value:     counterEntry.StoriesCount,
		}
		err = models.Db.Insert(&numOfStoriesVersion)
		helpers.PanicOnError(err)

		if models.TimestampType(counterEntry.Timestamp) < community.AddedTimestamp {
			community.AddedTimestamp = models.TimestampType(counterEntry.Timestamp)
			_, err := models.Db.Model(&community).
				Set("added_timestamp = ?added_timestamp").
				Where("pikabu_id = ?pikabu_id").
				Update()
			helpers.PanicOnError(err)
		}
	}
}
