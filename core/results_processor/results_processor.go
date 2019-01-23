package results_processor

import (
	"bitbucket.org/d3dev/parse_pikabu/amqp_helper"
	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg"
	"github.com/streadway/amqp"
	"gogsweb.2-47.ru/d3dev/pikago"
	"reflect"
	"time"
)

func Run() error {
	for true {
		err := startListener()
		if err != nil {
			if e, ok := err.(*errors.Error); ok {
				logger.Log.Error(e.ErrorStack())
			} else {
				logger.Log.Error(err.Error())
			}
			_ = amqp_helper.Cleanup()
		}
		time.Sleep(5 * time.Second)
	}

	return nil
}
func startListener() error {
	logger.Log.Debug("connecting to amqp server...")
	connection, err := amqp_helper.GetAMQPConnection(config.Settings.AMQPAddress)
	if err != nil {
		return errors.New(err)
	}

	ch, err := connection.Channel()
	if err != nil {
		return errors.New(err)
	}
	defer func() { _ = ch.Close() }()

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
		"bitbucket.org/d3dev/parse_pikabu/parser_results",
		true,
		false,
		false,
		false,
		nil,
		/*
			amqp.Table{
				"x-queue-mode": "lazy",
			},
		*/
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
	for message := range messages {
		err = processMessage(message)
		if err != nil {
			return err
		}
		err = message.Ack(false)
		if err != nil {
			return err
		}
	}
	logger.Log.Debug("stop waiting for parser results")

	return nil
}

func processMessage(message amqp.Delivery) error {
	logger.Log.Debugf("got message: %v", string(message.Body))

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
	case "communities_pages":
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