package parser

import (
	"bitbucket.org/d3dev/parse_pikabu/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"encoding/json"
	"github.com/go-errors/errors"
	"github.com/streadway/amqp"
	"reflect"
	"time"
)

var amqpConnections = map[string]*amqp.Connection{}

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

	if connection, ok := amqpConnections[this.Config.AMQPAddress]; ok {
		err := connection.Close()
		delete(amqpConnections, this.Config.AMQPAddress)
		if err != nil {
			return err
		}
	}

	return nil
}

func getAMQPConnection(amqpAddress string) (*amqp.Connection, error) {
	if connection, ok := amqpConnections[amqpAddress]; ok {
		return connection, nil
	}

	connection, err := amqp.Dial(amqpAddress)
	if err != nil {
		return nil, err
	}

	return connection, nil
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

	err = this.amqpChannel.ExchangeDeclare(
		"parser_results",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func (this *Parser) PutResultsToQueue(routingKey string, result interface{}) error {
	for i := 0; i < 2; i++ {
		connection, err := getAMQPConnection(this.Config.AMQPAddress)
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

		err = this.amqpChannel.Publish(
			"parser_results",
			routingKey,
			true,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent,
				Body:         message,
			},
		)
		if err == nil {
			return nil
		} else {
			_ = this.CleanupAMQP()
			logger.ParserLog.Error(err)
		}
	}

	return errors.Errorf("Unable to connect to queue")
}
