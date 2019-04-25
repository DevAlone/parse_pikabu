package resultsprocessor

import (
	"reflect"
	"runtime"
	"sync"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/globals"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-errors/errors"
	pikago_models "gogsweb.2-47.ru/d3dev/pikago/models"
)

// Run runs result processing
func Run() error {
	for {
		err := startListener()
		if err != nil {
			if e, ok := err.(*errors.Error); ok {
				logger.Log.Error(e.ErrorStack())
			} else {
				logger.Log.Error(err.Error())
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func startListener() error {
	var wg sync.WaitGroup
	for i := 0; i < config.Settings.NumberOfTasksProcessorsMultiplier*runtime.GOMAXPROCS(0); i++ {
		logger.Log.Debug("started results processor routine")
		wg.Add(1)
		go func() {
			for message := range globals.ParserResults {
				logger.LogError(processMessage(message))
			}
			wg.Done()
		}()
	}
	wg.Wait()

	return nil
}

func processMessage(message *models.ParserResult) error {
	logger.Log.Debugf("got message from parser %v", message)

	switch m := message.Results.(type) {
	case []models.ParserUserProfileResultData:
		userProfiles := []*pikago_models.UserProfile{}
		for _, result := range m {
			userProfiles = append(userProfiles, result.User)
		}
		return processUserProfiles(message.ParsingTimestamp, userProfiles)
	case []models.ParserUserProfileNotFoundResultData:
		return processUserProfileNotFoundResults(message.ParsingTimestamp, m)
	case []pikago_models.CommunitiesPage:
		return processCommunitiesPages(message.ParsingTimestamp, m)
	case []pikago_models.StoryGetResult:
		return processStoryGetResults(message.ParsingTimestamp, m)
	case []models.ParserStoryNotFoundResultData:
		return processStoryNotFoundResults(message.ParsingTimestamp, m)
	default:
		logger.Log.Warningf(
			"processMessage(): Unregistered result type \"%v\". Message: \"%v\". m: \"%v\"",
			reflect.TypeOf(m),
			message,
			m,
		)
	}

	return nil
}

// OldParserResultError shows that result of parser is too
// old and will be ignored
type OldParserResultError struct{}

func (e OldParserResultError) Error() string { return "old parser result error" }
