package tests

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/core/resultsprocessor"
	"bitbucket.org/d3dev/parse_pikabu/globals"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-pg/pg/orm"
	"github.com/stretchr/testify/assert"
	pikago_models "gogsweb.2-47.ru/d3dev/pikago/models"
)

func TestStoryParsing(t *testing.T) {
	err := config.UpdateSettingsFromFile("../core.config.json")
	helpers.PanicOnError(err)
	helpers.PanicOnError(globals.Init())

	initLogs()
	logger.Log.Debug(`start test "story parsing"`)

	err = models.InitDb()
	if err != nil {
		helpers.PanicOnError(err)
	}

	// clear tables
	for _, table := range models.Tables {
		err := models.Db.DropTable(table, &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			helpers.PanicOnError(err)
		}
	}

	// create again
	err = models.InitDb()
	if err != nil {
		helpers.PanicOnError(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// start results processor
	go func() {
		err := resultsprocessor.Run()
		if err != nil {
			helpers.PanicOnError(err)
		}
		wg.Done()
	}()

	err = pushResultToQueue(1, []pikago_models.StoryGetResult{
		pikago_models.StoryGetResult{
			QueryTimestampMs:       pikago_models.Int64{Value: 0},
			HasNextCommentsPage:    false,
			CommentsAreSortedBy:    "",
			MaxCommentsBranchDepth: pikago_models.Int64{Value: 99999},
			DeepCommentsAreHidden:  false,
			StoryData: &pikago_models.Story{
				StoryID: pikago_models.UInt64{Value: 1},
				Rating:  pikago_models.Int64{Value: 1},
				Title:   "test story",
				ContentBlocks: []pikago_models.StoryBlock{
					pikago_models.StoryBlock{
						Type:    "i",
						RawData: &json.RawMessage{},
						Data: pikago_models.Image{
							SmallURL:  "small url",
							LargeURL:  "large url",
							Animation: nil,
							Size: []pikago_models.UInt64{
								pikago_models.UInt64{Value: 10},
								pikago_models.UInt64{Value: 20},
							},
						},
					},
					pikago_models.StoryBlock{
						Type:    "t",
						RawData: &json.RawMessage{},
						Data:    "text block",
					},
				},
			},
			Comments: nil,
		},
	})
	helpers.PanicOnError(err)

	waitForResultsQueueEmpty()

	story := &models.PikabuStory{
		PikabuID: 1,
	}
	err = models.Db.Select(story)
	helpers.PanicOnError(err)

	assert.Equal(t, int32(1), story.Rating)
	assert.Equal(t, "test story", story.Title)
	assert.Equal(t, []models.PikabuStoryBlock{
		models.PikabuStoryBlock{
			Type: "i",
			Data: map[string]interface{}{
				"small":     "small url",
				"large":     "large url",
				"animation": nil,
				"img_size": []interface{}{
					float64(10),
					float64(20),
				},
			},
		},
		models.PikabuStoryBlock{
			Type: "t",
			Data: "text block",
		},
	}, story.ContentBlocks)

	// rating
	ratingVersions := []models.PikabuStoryRatingVersion{}
	err = models.Db.Model(&ratingVersions).
		Where("item_id = ?", story.PikabuID).
		Order("timestamp").
		Select()
	helpers.PanicOnError(err)

	assert.Equal(t, []models.PikabuStoryRatingVersion{}, ratingVersions)

	// title
	titleVersions := []models.PikabuStoryTitleVersion{}
	err = models.Db.Model(&titleVersions).
		Where("item_id = ?", story.PikabuID).
		Order("timestamp").
		Select()
	helpers.PanicOnError(err)

	assert.Equal(t, []models.PikabuStoryTitleVersion{}, titleVersions)

	// contentBlocks
	contentBlocksVersions := []models.PikabuStoryContentBlocksVersion{}
	err = models.Db.Model(&contentBlocksVersions).
		Where("item_id = ?", story.PikabuID).
		Order("timestamp").
		Select()
	helpers.PanicOnError(err)

	assert.Equal(t, []models.PikabuStoryContentBlocksVersion{}, contentBlocksVersions)

	/* push the same result again */

	err = pushResultToQueue(2, []pikago_models.StoryGetResult{
		pikago_models.StoryGetResult{
			QueryTimestampMs:       pikago_models.Int64{Value: 0},
			HasNextCommentsPage:    false,
			CommentsAreSortedBy:    "",
			MaxCommentsBranchDepth: pikago_models.Int64{Value: 99999},
			DeepCommentsAreHidden:  false,
			StoryData: &pikago_models.Story{
				StoryID: pikago_models.UInt64{Value: 1},
				Rating:  pikago_models.Int64{Value: 1},
				Title:   "test story",
				ContentBlocks: []pikago_models.StoryBlock{
					pikago_models.StoryBlock{
						Type:    "i",
						RawData: &json.RawMessage{},
						Data: pikago_models.Image{
							SmallURL:  "small url",
							LargeURL:  "large url",
							Animation: nil,
							Size: []pikago_models.UInt64{
								pikago_models.UInt64{Value: 10},
								pikago_models.UInt64{Value: 20},
							},
						},
					},
					pikago_models.StoryBlock{
						Type:    "t",
						RawData: &json.RawMessage{},
						Data:    "text block",
					},
				},
			},
			Comments: nil,
		},
	})
	helpers.PanicOnError(err)

	waitForResultsQueueEmpty()

	story = &models.PikabuStory{
		PikabuID: 1,
	}
	err = models.Db.Select(story)
	helpers.PanicOnError(err)

	assert.Equal(t, models.TimestampType(1), story.AddedTimestamp)
	assert.Equal(t, models.TimestampType(2), story.LastUpdateTimestamp)
	assert.Equal(t, int32(1), story.Rating)
	assert.Equal(t, "test story", story.Title)
	assert.Equal(t, []models.PikabuStoryBlock{
		models.PikabuStoryBlock{
			Type: "i",
			Data: map[string]interface{}{
				"small":     "small url",
				"large":     "large url",
				"animation": nil,
				"img_size": []interface{}{
					float64(10),
					float64(20),
				},
			},
		},
		models.PikabuStoryBlock{
			Type: "t",
			Data: "text block",
		},
	}, story.ContentBlocks)

	// rating
	ratingVersions = []models.PikabuStoryRatingVersion{}
	err = models.Db.Model(&ratingVersions).
		Where("item_id = ?", story.PikabuID).
		Order("timestamp").
		Select()
	helpers.PanicOnError(err)

	assert.Equal(t, []models.PikabuStoryRatingVersion{}, ratingVersions)

	// title
	titleVersions = []models.PikabuStoryTitleVersion{}
	err = models.Db.Model(&titleVersions).
		Where("item_id = ?", story.PikabuID).
		Order("timestamp").
		Select()
	helpers.PanicOnError(err)

	assert.Equal(t, []models.PikabuStoryTitleVersion{}, titleVersions)

	// contentBlocks
	contentBlocksVersions = []models.PikabuStoryContentBlocksVersion{}
	err = models.Db.Model(&contentBlocksVersions).
		Where("item_id = ?", story.PikabuID).
		Order("timestamp").
		Select()
	helpers.PanicOnError(err)

	assert.Equal(t, []models.PikabuStoryContentBlocksVersion{}, contentBlocksVersions)

	/* change some fields */

	err = pushResultToQueue(3, []pikago_models.StoryGetResult{
		pikago_models.StoryGetResult{
			QueryTimestampMs:       pikago_models.Int64{Value: 0},
			HasNextCommentsPage:    false,
			CommentsAreSortedBy:    "",
			MaxCommentsBranchDepth: pikago_models.Int64{Value: 99999},
			DeepCommentsAreHidden:  false,
			StoryData: &pikago_models.Story{
				StoryID: pikago_models.UInt64{Value: 1},
				Rating:  pikago_models.Int64{Value: 2},
				Title:   "new title",
				ContentBlocks: []pikago_models.StoryBlock{
					pikago_models.StoryBlock{
						Type:    "i",
						RawData: &json.RawMessage{},
						Data: pikago_models.Image{
							SmallURL:  "small url1",
							LargeURL:  "large url",
							Animation: nil,
							Size: []pikago_models.UInt64{
								pikago_models.UInt64{Value: 10},
								pikago_models.UInt64{Value: 20},
							},
						},
					},
				},
			},
			Comments: nil,
		},
	})
	helpers.PanicOnError(err)

	waitForResultsQueueEmpty()

	story = &models.PikabuStory{
		PikabuID: 1,
	}
	err = models.Db.Select(story)
	helpers.PanicOnError(err)

	assert.Equal(t, models.TimestampType(1), story.AddedTimestamp)
	assert.Equal(t, models.TimestampType(3), story.LastUpdateTimestamp)
	assert.Equal(t, int32(2), story.Rating)
	assert.Equal(t, "new title", story.Title)
	assert.Equal(t, []models.PikabuStoryBlock{
		models.PikabuStoryBlock{
			Type: "i",
			Data: map[string]interface{}{
				"small":     "small url1",
				"large":     "large url",
				"animation": nil,
				"img_size": []interface{}{
					float64(10),
					float64(20),
				},
			},
		},
	}, story.ContentBlocks)

	// rating
	ratingVersions = []models.PikabuStoryRatingVersion{}
	err = models.Db.Model(&ratingVersions).
		Where("item_id = ?", story.PikabuID).
		Order("timestamp").
		Select()
	helpers.PanicOnError(err)

	assert.Equal(t, []models.PikabuStoryRatingVersion{
		models.PikabuStoryRatingVersion{
			ItemId:    1,
			Timestamp: 1,
			Value:     1,
		},
		models.PikabuStoryRatingVersion{
			ItemId:    1,
			Timestamp: 2,
			Value:     1,
		},
		models.PikabuStoryRatingVersion{
			ItemId:    1,
			Timestamp: 3,
			Value:     2,
		},
	}, ratingVersions)

	// title
	titleVersions = []models.PikabuStoryTitleVersion{}
	err = models.Db.Model(&titleVersions).
		Where("item_id = ?", story.PikabuID).
		Order("timestamp").
		Select()
	helpers.PanicOnError(err)

	assert.Equal(t, []models.PikabuStoryTitleVersion{
		models.PikabuStoryTitleVersion{
			ItemId:    1,
			Timestamp: 1,
			Value:     "test story",
		},
		models.PikabuStoryTitleVersion{
			ItemId:    1,
			Timestamp: 2,
			Value:     "test story",
		},
		models.PikabuStoryTitleVersion{
			ItemId:    1,
			Timestamp: 3,
			Value:     "new title",
		},
	}, titleVersions)

	// contentBlocks
	contentBlocksVersions = []models.PikabuStoryContentBlocksVersion{}
	err = models.Db.Model(&contentBlocksVersions).
		Where("item_id = ?", story.PikabuID).
		Order("timestamp").
		Select()
	helpers.PanicOnError(err)

	assert.Equal(t, []models.PikabuStoryContentBlocksVersion{
		models.PikabuStoryContentBlocksVersion{
			ItemId:    1,
			Timestamp: 1,
			Value: []models.PikabuStoryBlock{
				models.PikabuStoryBlock{
					Type: "i",
					Data: map[string]interface{}{
						"small":     "small url",
						"large":     "large url",
						"animation": nil,
						"img_size": []interface{}{
							float64(10),
							float64(20),
						},
					},
				},
				models.PikabuStoryBlock{
					Type: "t",
					Data: "text block",
				},
			},
		},
		models.PikabuStoryContentBlocksVersion{
			ItemId:    1,
			Timestamp: 2,
			Value: []models.PikabuStoryBlock{
				models.PikabuStoryBlock{
					Type: "i",
					Data: map[string]interface{}{
						"small":     "small url",
						"large":     "large url",
						"animation": nil,
						"img_size": []interface{}{
							float64(10),
							float64(20),
						},
					},
				},
				models.PikabuStoryBlock{
					Type: "t",
					Data: "text block",
				},
			},
		},
		models.PikabuStoryContentBlocksVersion{
			ItemId:    1,
			Timestamp: 3,
			Value: []models.PikabuStoryBlock{
				models.PikabuStoryBlock{
					Type: "i",
					Data: map[string]interface{}{
						"small":     "small url1",
						"large":     "large url",
						"animation": nil,
						"img_size": []interface{}{
							float64(10),
							float64(20),
						},
					},
				},
			},
		},
	}, contentBlocksVersions)

	// clear tables
	for _, table := range models.Tables {
		err := models.Db.DropTable(table, &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			helpers.PanicOnError(err)
		}
	}
}

func pushResultToQueue(parsingTimestamp int64, stories []pikago_models.StoryGetResult) error {
	logger.Log.Debug(`pushing result to queue`)
	var pr models.ParserResult
	pr.ParsingTimestamp = models.TimestampType(parsingTimestamp)
	pr.ParserId = "d3dev/parser_id"
	pr.NumberOfResults = len(stories)
	pr.Results = stories
	globals.ParserResults <- &pr

	return nil
}

func waitForResultsQueueEmpty() {
	logger.Log.Debug(`waiting for queue to become empty`)
	for {
		if len(globals.ParserResults) == 0 {
			// TODO: check whether the message was actually processed
			time.Sleep(1 * time.Second)
			return
		}

		time.Sleep(500 * time.Millisecond)
	}
}
