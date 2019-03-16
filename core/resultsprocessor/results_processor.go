package resultsprocessor

import (
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/go-pg/pg/orm"

	"bitbucket.org/d3dev/parse_pikabu/globals"

	"bitbucket.org/d3dev/parse_pikabu/amqphelper"
	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg"
	"github.com/streadway/amqp"
	"gogsweb.2-47.ru/d3dev/pikago"
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
			if !globals.SingleProcessMode {
				_ = amqphelper.Cleanup()
			}
		}
		time.Sleep(5 * time.Second)
	}

	return nil
}
func startListener() error {
	if globals.SingleProcessMode {
		return startListenerChannels()
	}

	return startListenerAMQP()
}

func startListenerChannels() error {
	var wg sync.WaitGroup
	for i := 0; i < config.Settings.NumberOfTasksProcessorsMultiplier*runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func() {
			for message := range globals.ParserResults {
				err := processMessage(message)
				if err != nil {
					logger.Log.Error(err)
				}
			}

			wg.Done()
		}()
	}
	wg.Wait()

	return nil
}

func startListenerAMQP() error {
	logger.Log.Debug("connecting to amqp server...")
	connection, err := amqphelper.GetAMQPConnection(config.Settings.AMQPAddress)
	if err != nil {
		return errors.New(err)
	}

	ch, err := connection.Channel()
	if err != nil {
		return errors.New(err)
	}
	defer func() { _ = ch.Close() }()

	err = amqphelper.DeclareExchanges(ch)
	if err != nil {
		return errors.New(err)
	}

	q, err := amqphelper.DeclareParserResultsQueue(ch)
	if err != nil {
		return errors.New(err)
	}

	err = ch.QueueBind(
		q.Name,
		"",
		"parser_results",
		false,
		nil,
	)
	if err != nil {
		return errors.New(err)
	}

	err = ch.Qos(2, 0, false)
	if err != nil {
		return errors.New(err)
	}

	messages, err := ch.Consume(
		q.Name,
		"", // routing key
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.New(err)
	}

	logger.Log.Debug("start waiting for parser results")
	var wg sync.WaitGroup
	for i := 0; i < config.Settings.NumberOfTasksProcessorsMultiplier*runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func() {
			for message := range messages {
				err = processAMQPMessage(message)
				if err != nil {
					logger.Log.Error(err)
					if e, ok := err.(*errors.Error); ok {
						logger.Log.Error(e.ErrorStack())
					}
					// panic(err)
				}
				err = message.Ack(false)
				if err != nil {
					logger.Log.Error(err)
					if e, ok := err.(*errors.Error); ok {
						logger.Log.Error(e.ErrorStack())
					}
					// panic(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	logger.Log.Debug("stop waiting for parser results")

	return nil
}

func processAMQPMessage(message amqp.Delivery) error {
	logger.Log.Debugf("got message: %v", string(message.Body))

	switch message.RoutingKey {
	case "user_profile":
		var resp models.ParserUserProfileResult
		err := pikago.JsonUnmarshal(message.Body, &resp, true)
		if err != nil {
			return errors.New(err)
		}

		if len(resp.Results) < 1 {
			return errors.Errorf("bad result: %v", resp)
		}

		userProfiles := []*pikago_models.UserProfile{}
		for _, result := range resp.Results {
			userProfiles = append(userProfiles, result.User)
		}

		err = processUserProfiles(resp.ParsingTimestamp, userProfiles)
		if err != nil {
			return err
		}
	case "user_profile_not_found":
		var resp models.ParserUserProfileNotFoundResult
		err := pikago.JsonUnmarshal(message.Body, &resp, true)
		if err != nil {
			return errors.New(err)
		}

		if len(resp.Results) < 1 {
			return errors.Errorf("bad result: %v", resp)
		}

		err = processUserProfileNotFoundResults(resp.ParsingTimestamp, resp.Results)
		if err != nil {
			return err
		}
	case "communities_pages":
		var resp models.ParserCommunitiesPageResult
		err := pikago.JsonUnmarshal(message.Body, &resp, true)
		if err != nil {
			return errors.New(err)
		}

		if len(resp.Results) < 1 {
			return errors.Errorf("bad result: %v", resp)
		}

		err = processCommunitiesPages(resp.ParsingTimestamp, resp.Results)
		if err != nil {
			return errors.New(err)
		}
	default:
		logger.Log.Warningf(
			"Unregistered result type \"%v\". Message: \"%v\"",
			message.RoutingKey,
			string(message.Body),
		)
	}

	return nil
}

func processMessage(message *models.ParserResult) error {
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

func processModelFieldsVersions(
	transaction *pg.Tx,
	oldModelPtr interface{},
	newModelPtr interface{},
	parsingTimestamp models.TimestampType,
) (bool, error) {
	// TODO: consider adding lock here
	wasDataChanged := false

	if reflect.TypeOf(oldModelPtr) != reflect.TypeOf(newModelPtr) {
		return false, errors.New("types should be equal")
	}

	oldModel := reflect.ValueOf(oldModelPtr).Elem()
	newModel := reflect.ValueOf(newModelPtr).Elem()

	oldID := oldModel.FieldByName("PikabuId").Uint()
	newID := newModel.FieldByName("PikabuId").Uint()

	if oldID != newID {
		return false, errors.New("ids should be equal")
	}

	addedTimestamp := models.TimestampType(oldModel.FieldByName("AddedTimestamp").Int())
	lastUpdateTimestamp := models.TimestampType(oldModel.FieldByName("LastUpdateTimestamp").Int())

	if parsingTimestamp <= lastUpdateTimestamp {
		// TODO: find a better way
		return false, OldParserResultError{}
	}

	oldModelType := reflect.TypeOf(oldModel.Interface())

	for i := 0; i < oldModelType.NumField(); i++ {
		fieldType := oldModelType.Field(i)

		_, isVersionedField := fieldType.Tag.Lookup("gen_versions")
		if !isVersionedField {
			continue
		}

		oldField := oldModel.FieldByName(fieldType.Name)
		newField := newModel.FieldByName(fieldType.Name)

		if reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
			continue
		}

		wasDataChanged = true

		// generate versions
		versionTable := models.FieldsVersionTablesMap[oldModelType.Name()+fieldType.Name+"Version"]

		insertVersion := func(
			timestamp models.TimestampType,
			value reflect.Value,
			ignoreIfExists bool,
		) error {
			e := reflect.ValueOf(versionTable).Elem()
			e.FieldByName("ItemId").SetUint(oldID)
			e.FieldByName("Timestamp").Set(reflect.ValueOf(timestamp))
			e.FieldByName("Value").Set(value)

			var err error
			if ignoreIfExists {
				var q *orm.Query
				if transaction == nil {
					q = models.Db.Model(versionTable)
				} else {
					q = transaction.Model(versionTable)
				}
				_, err = q.
					OnConflict("DO NOTHING").
					Insert(versionTable)
			} else {
				if transaction == nil {
					err = models.Db.Insert(versionTable)
				} else {
					err = transaction.Insert(versionTable)
				}
			}
			if err != nil {
				return errors.New(err)
			}

			return nil
		}

		var q *orm.Query
		if transaction == nil {
			q = models.Db.Model(versionTable)
		} else {
			q = transaction.Model(versionTable)
		}
		count, err := q.Where("item_id = ?", oldID).Count()
		if err != nil {
			return false, errors.New(err)
		}

		if count == 0 {
			err := insertVersion(
				addedTimestamp,
				oldField,
				false)
			if err != nil {
				return false, err
			}
		}

		err = insertVersion(
			lastUpdateTimestamp,
			oldField,
			true)
		if err != nil {
			return false, err
		}

		err = insertVersion(
			parsingTimestamp,
			newField,
			false)
		if err != nil {
			return false, err
		}

		// set the field
		if !oldField.CanSet() {
			panic(errors.Errorf("field %v from model %v cannot be set", oldField, oldModel))
		}
		oldField.Set(newField)
	}

	return wasDataChanged, nil
}
