package parser

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/globals"

	"bitbucket.org/d3dev/parse_pikabu/amqp_helper"
	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"bitbucket.org/d3dev/parse_pikabu/parser/logger"
	"github.com/go-errors/errors"
	"github.com/streadway/amqp"
)

func (p *Parser) Cleanup() error {
	return p.CleanupAMQP()
}

func (p *Parser) CleanupAMQP() error {
	if p.amqpChannel != nil {
		err := p.amqpChannel.Close()
		p.amqpChannel = nil
		if err != nil {
			return err
		}
	}

	if c, ok := amqp_helper.AmqpConnections.Get(p.Config.AMQPAddress); ok {
		connection := c.(*amqp.Connection)
		err := connection.Close()
		amqp_helper.AmqpConnections.Remove(p.Config.AMQPAddress)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) initChannel(connection *amqp.Connection) error {
	if p.amqpChannel != nil {
		return nil
	}

	var err error
	p.amqpChannel, err = connection.Channel()
	if err != nil {
		return err
	}

	err = amqp_helper.DeclareExchanges(p.amqpChannel)
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) PutResultsToQueue(routingKey string, result interface{}) error {
	if globals.SingleProcessMode {
		numberOfResults := 0
		resultType := reflect.TypeOf(result)
		switch resultType.Kind() {
		case reflect.Slice, reflect.Array:
			numberOfResults = reflect.ValueOf(result).Len()
		default:
			result = []interface{}{result}
			numberOfResults = 1
		}

		var pr models.ParserResult
		pr.ParsingTimestamp = models.TimestampType(time.Now().Unix())
		pr.ParserId = "d3dev/" + p.Config.ParserId
		pr.NumberOfResults = numberOfResults
		pr.Results = result
		globals.ParserResults <- &pr

		return nil
	} else {
		return p.PutResultsToAMQPQueue(routingKey, result)
	}
}

func (p *Parser) PutResultsToAMQPQueue(routingKey string, result interface{}) error {

	for i := 0; i < 2; i++ {
		connection, err := amqp_helper.GetAMQPConnection(p.Config.AMQPAddress)
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
		jsonMessage.ParserId = "d3dev/" + p.Config.ParserId
		jsonMessage.NumberOfResults = numberOfResults
		jsonMessage.Results = result

		message, err := json.Marshal(jsonMessage)
		if err != nil {
			return err
		}

		err = p.initChannel(connection)
		if err != nil {
			return err
		}

		err = p.waitResultsQueueForEmpty(connection)
		if err != nil {
			return err
		}

		err = p.amqpChannel.Publish(
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
			_ = p.CleanupAMQP()
			logger.Log.Error(err)
		}
	}

	return errors.Errorf("Unable to connect to queue")
}

func (p *Parser) waitResultsQueueForEmpty(connection *amqp.Connection) error {
	for true {
		parserTasksQueue, err := p.amqpChannel.QueueInspect("bitbucket.org/d3dev/parse_pikabu/parser_results")
		if err != nil {
			// TODO: fix
			if !strings.Contains(err.Error(), "NOT_FOUND - no queue") {
				return err
			} else {
				p.amqpChannel = nil
				_ = p.initChannel(connection)
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
