package tests

import (
	"github.com/DevAlone/parse_pikabu/core/logger"
	"github.com/DevAlone/parse_pikabu/core/resultsprocessor"
	"github.com/DevAlone/parse_pikabu/helpers"
	"github.com/DevAlone/parse_pikabu/modelhooks"
	"github.com/DevAlone/parse_pikabu/models"
	"github.com/DevAlone/parse_pikabu/telegram"
	"github.com/stretchr/testify/assert"
	pikago_models "gogsweb.2-47.ru/d3dev/pikago/models"
	"sync"
	"testing"
)

func TestSomeShit(t *testing.T) {

}

func TestDeletedCommentNotifications(t *testing.T) {
	initTestEnvironment()
	logger.Log.Debug(`start test "deleted comment notification"`)

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
	// start modelhooks handler
	go func() {
		err := modelhooks.RunModelHooksHandler()
		helpers.PanicOnError(err)
		wg.Done()
	}()

	wg.Add(1)
	// start telegram notifier
	go func() {
		err := telegram.RunTelegramNotifier()
		helpers.PanicOnError(err)
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
				Title:   "null",
			},
			Comments: []pikago_models.Comment{
				pikago_models.Comment{
					ID: pikago_models.UInt64{Value: 1},
					Content: pikago_models.CommentContent{
						Text: "porn",
					},
				},
			},
		},
	})
	helpers.PanicOnError(err)

	waitForResultsQueueEmpty()

	/* push deleted version of the same result*/
	err = pushStoriesResultToQueue(2, []pikago_models.StoryGetResult{
		pikago_models.StoryGetResult{
			QueryTimestampMs:       pikago_models.Int64{Value: 0},
			HasNextCommentsPage:    false,
			CommentsAreSortedBy:    "",
			MaxCommentsBranchDepth: pikago_models.Int64{Value: 99999},
			DeepCommentsAreHidden:  false,
			StoryData: &pikago_models.Story{
				StoryID: pikago_models.UInt64{Value: 1},
				Title:   "null",
			},
			Comments: []pikago_models.Comment{
				pikago_models.Comment{
					ID: pikago_models.UInt64{Value: 1},
					Content: pikago_models.CommentContent{
						Text: "deleted by @moderator",
					},
					IsDeleted: true,
				},
			},
		},
	})

	helpers.PanicOnError(err)

	waitForResultsQueueEmpty()

	comment := &models.PikabuComment{
		PikabuID: 1,
	}
	err = models.Db.Select(comment)
	helpers.PanicOnError(err)

	assert.Equal(t, models.TimestampType(1), comment.AddedTimestamp)
	assert.Equal(t, true, comment.IsDeleted)

	clearAndInitDb()
}
