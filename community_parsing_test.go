package main

import (
	"bitbucket.org/d3dev/parse_pikabu/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/results_processor"
	"github.com/go-pg/pg/orm"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestCommunityParsing(t *testing.T) {
	logger.Log.Debug(`start test "community parsing"`)

	// init db, create list of tables
	err := models.InitDb()
	if err != nil {
		handleError(err)
	}

	// clear tables
	for _, table := range models.Tables {
		err := models.Db.DropTable(table, &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			handleError(err)
		}
	}

	// create again
	err = models.InitDb()
	if err != nil {
		handleError(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// start results processor
	go func() {
		err := results_processor.Run()
		if err != nil {
			handleError(err)
		}
		wg.Done()
	}()

	err = pushTaskToQueue("communities_page", []byte(`
{
	"parsing_timestamp": 1,
	"parser_id": "d3dev/parser_id",
	"number_of_results": 1,
	"results": [
		{
			"total": 2,
			"eof": true,
			"list": [
				{
					"id": 1,
					"name": "name1",
					"link_name": "link_name1",
					"url": "url1",
					"avatar_url": "avatar_url1",
					"bg_image_url": "bg_image_url1",
					"tags": [
						"tag1", "tag2", "tag3"
					],
					"stories": 10,
					"subscribers": 11,
					"description": "description1",
					"rules": "rules1",
					"restriction": "restriction1",
					"is_ignored": true,
					"is_subscribed": false,
					"community_admin": {
						"id": 100,
						"name": "doesn't matter'",
						"avatar": "doesn't matter'"
					},
					"community_moderators": [
						{
							"id": 1000,
							"name": "doesn't matter'",
							"avatar": "doesn't matter'"
						},
						{
							"id": 1001,
							"name": "doesn't matter'",
							"avatar": "doesn't matter'"
						},
						{
							"id": 1002,
							"name": "doesn't matter'",
							"avatar": "doesn't matter'"
						}
					]
				}
			]
		}
	]
}
`,
	))
	handleError(err)

	waitForQueueEmpty()

	communities := []models.PikabuCommunity{}
	err = models.Db.Model(&communities).
		Order("pikabu_id").
		Select()
	handleError(err)

	assert.Equal(t, []models.PikabuCommunity{
		{
			PikabuId:            1,
			Name:                "name1",
			LinkName:            "link_name1",
			URL:                 "url1",
			AvatarURL:           "avatar_url1",
			BackgroundImageURL:  "bg_image_url1",
			Tags:                []string{"tag1", "tag2", "tag3"},
			NumberOfStories:     10,
			NumberOfSubscribers: 11,
			Description:         "description1",
			Rules:               "rules1",
			Restrictions:        "restriction1",
			AdminId:             100,
			ModeratorIds:        []uint64{1000, 1001, 1002},
			AddedTimestamp:      1,
			LastUpdateTimestamp: 1,
		},
	}, communities)

	// make some changes
	err = pushTaskToQueue("communities_page", []byte(`
{
	"parsing_timestamp": 2,
	"parser_id": "d3dev/parser_id",
	"number_of_results": 1,
	"results": [
		{
			"total": 2,
			"eof": true,
			"list": [
				{
					"id": 1,
					"name": "name2",
					"link_name": "link_name2",
					"url": "url2",
					"avatar_url": "avatar_url2",
					"bg_image_url": "bg_image_url2",
					"tags": [
						"tag1", "tag3"
					],
					"stories": 10,
					"subscribers": 11,
					"description": "description1",
					"rules": "rules1",
					"restriction": "restriction1",
					"is_ignored": true,
					"is_subscribed": false,
					"community_admin": {
						"id": 100,
						"name": "doesn't matter'",
						"avatar": "doesn't matter'"
					},
					"community_moderators": [
						{
							"id": 1000,
							"name": "doesn't matter'",
							"avatar": "doesn't matter'"
						},
						{
							"id": 1001,
							"name": "doesn't matter'",
							"avatar": "doesn't matter'"
						},
						{
							"id": 1002,
							"name": "doesn't matter'",
							"avatar": "doesn't matter'"
						},
						{
							"id": 1003,
							"name": "doesn't matter'",
							"avatar": "doesn't matter'"
						}
					]
				}
			]
		}
	]
}
`,
	))
	handleError(err)

	waitForQueueEmpty()

	communities = []models.PikabuCommunity{}
	err = models.Db.Model(&communities).Order("pikabu_id").Select()
	handleError(err)

	assert.Equal(t, []models.PikabuCommunity{
		{
			PikabuId:            1,
			Name:                "name2",
			LinkName:            "link_name2",
			URL:                 "url2",
			AvatarURL:           "avatar_url2",
			BackgroundImageURL:  "bg_image_url2",
			Tags:                []string{"tag1", "tag3"},
			NumberOfStories:     10,
			NumberOfSubscribers: 11,
			Description:         "description1",
			Rules:               "rules1",
			Restrictions:        "restriction1",
			AdminId:             100,
			ModeratorIds:        []uint64{1000, 1001, 1002, 1003},
			AddedTimestamp:      1,
			LastUpdateTimestamp: 2,
		},
	}, communities)

	linkNameVersions := []models.PikabuCommunityLinkNameVersion{}
	err = models.Db.Model(&linkNameVersions).
		Where("item_id = ?", 1).
		Order("timestamp").
		Select()
	handleError(err)

	assert.Equal(t, []models.PikabuCommunityLinkNameVersion{
		{ItemId: 1, Timestamp: 1, Value: "link_name1"},
		{ItemId: 1, Timestamp: 2, Value: "link_name2"},
	}, linkNameVersions)

	tagsVersions := []models.PikabuCommunityTagsVersion{}
	err = models.Db.Model(&tagsVersions).
		Where("item_id = ?", 1).
		Order("timestamp").
		Select()
	handleError(err)

	assert.Equal(t, []models.PikabuCommunityTagsVersion{
		{ItemId: 1, Timestamp: 1, Value: []string{"tag1", "tag2", "tag3"}},
		{ItemId: 1, Timestamp: 2, Value: []string{"tag1", "tag3"}},
	}, tagsVersions)

	moderatorIdsVersions := []models.PikabuCommunityModeratorIdsVersion{}
	err = models.Db.Model(&moderatorIdsVersions).
		Where("item_id = ?", 1).
		Order("timestamp").
		Select()
	handleError(err)

	assert.Equal(t, []models.PikabuCommunityModeratorIdsVersion{
		{ItemId: 1, Timestamp: 1, Value: []uint64{1000, 1001, 1002}},
		{ItemId: 1, Timestamp: 2, Value: []uint64{1000, 1001, 1002, 1003}},
	}, moderatorIdsVersions)
}
