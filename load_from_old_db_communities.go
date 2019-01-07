package main

import (
	"fmt"

	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/old_models"
)

func processCommunities() {
	fmt.Println("processing communities")

	var communities []models.PikabuCommunity

	err := models.Db.Model(&communities).Select()
	panicOnError(err)

	for _, item := range communities {
		var oldCommunity old_models.Community
		err := oldDb.Model(&oldCommunity).
			Where("url_name ilike ?", item.LinkName).
			Select()
		panicOnError(err)

		var counterEntries []old_models.CommunityCountersEntry
		err = oldDb.Model(&counterEntries).
			Where("community_id = ?", oldCommunity.Id).
			Order("timestamp").
			Select()
		panicOnError(err)

		for _, counterEntry := range counterEntries {
			numOfSubsVersion := models.PikabuCommunityNumberOfSubscribersVersion{
				Timestamp: models.TimestampType(counterEntry.Timestamp),
				ItemId:    item.PikabuId,
				Value:     counterEntry.SubscribersCount,
			}
			err := models.Db.Insert(&numOfSubsVersion)
			panicOnError(err)

			numOfStoriesVersion := models.PikabuCommunityNumberOfStoriesVersion{
				Timestamp: models.TimestampType(counterEntry.Timestamp),
				ItemId:    item.PikabuId,
				Value:     counterEntry.StoriesCount,
			}
			err = models.Db.Insert(&numOfStoriesVersion)
			panicOnError(err)
		}
	}
}
