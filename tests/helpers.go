package tests

import (
	"os"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/globals"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-pg/pg/orm"
	pikago_models "gogsweb.2-47.ru/d3dev/pikago/models"
)

func initLogs() {
	logger.Init()
}

func initTestEnvironment() {
	helpers.PanicOnError(os.Chdir("../"))
	err := config.UpdateSettingsFromFile("core.config.json")
	helpers.PanicOnError(err)
	helpers.PanicOnError(globals.Init())

	initLogs()
}

func clearAndInitDb() {
	err := models.InitDb()
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

}

func pushResultToQueue(parsingTimestamp int64, result interface{}) {

}

func pushStoriesResultToQueue(parsingTimestamp int64, stories []pikago_models.StoryGetResult) error {
	logger.Log.Debug(`pushing result to queue`)
	var pr models.ParserResult
	pr.ParsingTimestamp = models.TimestampType(parsingTimestamp)
	pr.ParserID = "d3dev/parser_id"
	pr.NumberOfResults = len(stories)
	pr.Results = stories
	globals.ParserResults <- &pr

	return nil
}

func waitForTasksQueueEmpty() {
	logger.Log.Debug(`waiting for tasks queue to become empty`)
	for {
		// TODO: fix
		/*
			if len(globals.ParserParseStoryTasks) == 0 {
				time.Sleep(1 * time.Second)
				return
			}
		*/

		time.Sleep(500 * time.Millisecond)
	}
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
