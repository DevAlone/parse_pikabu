package results_processor

import (
	"bitbucket.org/d3dev/parse_pikabu/config"
	"bitbucket.org/d3dev/parse_pikabu/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg"
	"github.com/streadway/amqp"
	"gogsweb.2-47.ru/d3dev/pikago"
	"os"
	"reflect"
	"time"
)

func Run() error {
	for true {
		err := startListener()
		if err != nil {
			if e, ok := err.(*errors.Error); ok {
				logger.ParserLog.Error(e.ErrorStack())
			} else {
				logger.ParserLog.Error(err.Error())
			}
		}
		time.Sleep(5 * time.Second)
	}

	return nil
}
func startListener() error {
	logger.Log.Debug("connecting to amqp server...")
	conn, err := amqp.Dial(config.Settings.AMQPAddress)
	if err != nil {
		return errors.New(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return errors.New(err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"parser_results",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.New(err)
	}

	q, err := ch.QueueDeclare(
		"bitbucket.org/d3dev/parse_pikabu",
		true,
		false,
		false,
		false,
		nil,
	)
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

	messages, err := ch.Consume(
		q.Name,
		"", // routing key
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.New(err)
	}

	logger.Log.Debug("start waiting for parser results")
	for message := range messages {
		logger.Log.Debugf("got parser result: %v", string(message.Body))
		err = processMessage(message)
		if err != nil {
			return err
		}
	}
	logger.Log.Debug("stop waiting for parser results")

	return nil
}

func processMessage(message amqp.Delivery) error {
	// TODO: check if result is from the future
	writeToFile := func(str string) {
		f, err := os.OpenFile("results.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		_, err = f.WriteString(str + "\n")
		if err != nil {
			panic(err)
		}
		err = f.Sync()
		if err != nil {
			panic(err)
		}
	}
	writeToFile(string(message.Body))

	switch message.RoutingKey {
	case "user_profile":
		var resp models.ParserUserProfileResult
		err := pikago.JsonUnmarshal(message.Body, &resp)
		if err != nil {
			return errors.New(err)
		}

		if len(resp.Results) < 1 {
			return errors.Errorf("bad result: %v", resp)
		}

		userProfiles := []*pikago.UserProfile{}
		for _, result := range resp.Results {
			userProfiles = append(userProfiles, result.User)
		}

		err = processUserProfiles(resp.ParsingTimestamp, userProfiles)
		if err != nil {
			return err
		}
	case "communities_page":
		var resp models.ParserCommunitiesPageResult
		err := pikago.JsonUnmarshal(message.Body, &resp)
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

type OldParserResultError struct{}

func (this OldParserResultError) Error() string { return "old parser result error" }

func processModelFieldsVersions(
	tx *pg.Tx,
	oldModelPtr interface{},
	newModelPtr interface{},
	parsingTimestamp models.TimestampType,
) (bool, error) {
	wasDataChanged := false

	if reflect.TypeOf(oldModelPtr) != reflect.TypeOf(newModelPtr) {
		return false, errors.New("types should be equal")
	}

	oldModel := reflect.ValueOf(oldModelPtr).Elem()
	newModel := reflect.ValueOf(newModelPtr).Elem()

	oldId := oldModel.FieldByName("PikabuId").Uint()
	newId := newModel.FieldByName("PikabuId").Uint()

	if oldId != newId {
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
			e.FieldByName("ItemId").SetUint(oldId)
			e.FieldByName("Timestamp").Set(reflect.ValueOf(timestamp))
			e.FieldByName("Value").Set(value)

			var err error
			if ignoreIfExists {
				_, err = tx.Model(versionTable).
					OnConflict("DO NOTHING").
					Insert(versionTable)
			} else {
				err = tx.Insert(versionTable)
			}

			return err
		}

		count, err := tx.Model(versionTable).Where("item_id = ?", oldId).Count()
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
