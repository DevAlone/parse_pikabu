package tests

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/core/resultsprocessor"
	"bitbucket.org/d3dev/parse_pikabu/core/taskmanager"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/parser"
	"github.com/go-pg/pg"
	"github.com/stretchr/testify/assert"
	pikago_models "gogsweb.2-47.ru/d3dev/pikago/models"
)

func TestStoryParsing(t *testing.T) {
	initTestEnvironment()
	logger.Log.Debug(`start test "story parsing"`)

	clearAndInitDb()

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

	err := pushStoriesResultToQueue(1, []pikago_models.StoryGetResult{
		pikago_models.StoryGetResult{
			QueryTimestampMs:       pikago_models.Int64{Value: 0},
			HasNextCommentsPage:    false,
			CommentsAreSortedBy:    "",
			MaxCommentsBranchDepth: pikago_models.Int64{Value: 99999},
			DeepCommentsAreHidden:  false,
			StoryData: &pikago_models.Story{
				StoryID: pikago_models.UInt64{Value: 1},
				Rating:  pikago_models.NullableInt64{Value: 1, IsNull: false},
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

	err = pushStoriesResultToQueue(2, []pikago_models.StoryGetResult{
		pikago_models.StoryGetResult{
			QueryTimestampMs:       pikago_models.Int64{Value: 0},
			HasNextCommentsPage:    false,
			CommentsAreSortedBy:    "",
			MaxCommentsBranchDepth: pikago_models.Int64{Value: 99999},
			DeepCommentsAreHidden:  false,
			StoryData: &pikago_models.Story{
				StoryID: pikago_models.UInt64{Value: 1},
				Rating:  pikago_models.NullableInt64{Value: 1, IsNull: false},
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

	err = pushStoriesResultToQueue(3, []pikago_models.StoryGetResult{
		pikago_models.StoryGetResult{
			QueryTimestampMs:       pikago_models.Int64{Value: 0},
			HasNextCommentsPage:    false,
			CommentsAreSortedBy:    "",
			MaxCommentsBranchDepth: pikago_models.Int64{Value: 99999},
			DeepCommentsAreHidden:  false,
			StoryData: &pikago_models.Story{
				StoryID: pikago_models.UInt64{Value: 1},
				Rating:  pikago_models.NullableInt64{Value: 2, IsNull: false},
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

	clearAndInitDb()
}

func TestStoryDataIsNil(t *testing.T) {
	initTestEnvironment()
	logger.Log.Debug(`start test "story data is nil"`)

	clearAndInitDb()

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

	wg.Add(1)
	go func() {
		parser.Main()
		wg.Done()
	}()

	story := &models.PikabuStory{
		PikabuID: 59,
	}

	for {
		/* create task for parsing story number 59 */
		err := taskmanager.ForceAddParseStoryTask(59)
		helpers.PanicOnError(err)

		time.Sleep(1 * time.Second)
		logger.Log.Debug("trying to select a story from db")
		err = models.Db.Select(story)
		if err == pg.ErrNoRows {
			continue
		}
		helpers.PanicOnError(err)
		break
	}

	assert.Equal(t, "Смешной пацан с мороженным", story.Title)
	assert.Equal(t, []models.PikabuStoryBlock{
		models.PikabuStoryBlock{
			Type: "i",
			Data: map[string]interface{}{
				"animation": interface{}(nil),
				"img_size":  interface{}(nil),
				"large":     "https://cs.pikabu.ru/images/big_size_comm/2012-01_6/13275325751151.gif",
				"small":     "https://cs.pikabu.ru/images/big_size_comm/2012-01_6/13275325751151.gif",
			},
		},
	}, story.ContentBlocks)

	/* clean up */
	clearAndInitDb()
}
