package parser

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/amqp_helper"
	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/parser/logger"
	"github.com/go-errors/errors"
	"github.com/streadway/amqp"
)

func (this *Parser) Cleanup() error {
	return this.CleanupAMQP()
}

func (this *Parser) CleanupAMQP() error {
	if this.amqpChannel != nil {
		err := this.amqpChannel.Close()
		this.amqpChannel = nil
		if err != nil {
			return err
		}
	}

	if c, ok := amqp_helper.AmqpConnections.Get(this.Config.AMQPAddress); ok {
		connection := c.(*amqp.Connection)
		err := connection.Close()
		amqp_helper.AmqpConnections.Remove(this.Config.AMQPAddress)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *Parser) initChannel(connection *amqp.Connection) error {
	if this.amqpChannel != nil {
		return nil
	}

	var err error
	this.amqpChannel, err = connection.Channel()
	if err != nil {
		return err
	}

	err = amqp_helper.DeclareExchanges(this.amqpChannel)
	if err != nil {
		return err
	}

	return nil
}

func (this *Parser) PutResultsToQueue(routingKey string, result interface{}) error {
	for i := 0; i < 2; i++ {
		connection, err := amqp_helper.GetAMQPConnection(this.Config.AMQPAddress)
		if err != nil {
			return err
		}

		numberOfResults := 0
		resultType := reflect.TypeOf(result)
		switch resultType.Kind() {
		case reflect.Slice, reflect.Array:
			numberOfResults = reflect.ValueOf(result).Len()
		default:
			result = []interface{}{result}
			numberOfResults = 1
		}

		var jsonMessage models.ParserResult
		jsonMessage.ParsingTimestamp = models.TimestampType(time.Now().Unix())
		jsonMessage.ParserId = "d3dev/" + this.Config.ParserId
		jsonMessage.NumberOfResults = numberOfResults
		jsonMessage.Results = result

		message, err := json.Marshal(jsonMessage)
		if err != nil {
			return err
		}

		err = this.initChannel(connection)
		if err != nil {
			return err
		}

		err = this.waitResultsQueueForEmpty(connection)
		if err != nil {
			return err
		}

		err = this.amqpChannel.Publish(
			"parser_results",
			routingKey,
			true,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Transient,
				Body:         message,
			},
		)
		if err == nil {
			return nil
		} else {
			_ = this.CleanupAMQP()
			logger.Log.Error(err)
		}
	}

	return errors.Errorf("Unable to connect to queue")
}

func (this *Parser) waitResultsQueueForEmpty(connection *amqp.Connection) error {
	for true {
		parserTasksQueue, err := this.amqpChannel.QueueInspect("bitbucket.org/d3dev/parse_pikabu/parser_results")
		if err != nil {
			// TODO: fix
			if !strings.Contains(err.Error(), "NOT_FOUND - no queue") {
				return err
			} else {
				this.amqpChannel = nil
				_ = this.initChannel(connection)
				return errors.New(err)
			}
		}

		if parserTasksQueue.Messages <= config.Settings.MaxNumberOfTasksInQueue {
			return nil
		}

		time.Sleep(50 * time.Millisecond)
	}
	return nil
}
