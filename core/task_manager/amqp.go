package task_manager

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/globals"

	"bitbucket.org/d3dev/parse_pikabu/amqp_helper"
	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-errors/errors"
	"github.com/streadway/amqp"
)

var amqpChannel *amqp.Channel

func Cleanup() error {
	if amqpChannel != nil {
		err := amqpChannel.Close()
		amqpChannel = nil
		if err != nil {
			return err
		}
	}

	return nil
}

func initChannel(connection *amqp.Connection) error {
	if amqpChannel != nil {
		return nil
	}

	var err error
	amqpChannel, err = connection.Channel()
	if err != nil {
		return err
	}

	err = amqp_helper.DeclareExchanges(amqpChannel)
	if err != nil {
		return err
	}

	return nil
}

func PushTaskToQueue(taskPtr interface{}) error {
	if globals.SingleProcessMode {
		switch t := taskPtr.(type) {
		case *models.ParseUserTask:
			globals.ParserParseUserTasks <- t
			return nil
		case *models.SimpleTask:
			globals.ParserSimpleTasks <- t
			return nil
		default:
			return errors.Errorf("trying to push undeclared type of task %v %v", reflect.TypeOf(t), t)
		}
	} else {
		return PushTaskToAMQPQueue(taskPtr)
	}
}

func PushTaskToAMQPQueue(taskPtr interface{}) error {
	// TODO: limit number of tasks in queue
	for i := 0; i < 2; i++ {
		connection, err := amqp_helper.GetAMQPConnection(config.Settings.AMQPAddress)
		if err != nil {
			return errors.New(err)
		}
		var routingKey string
		switch t := taskPtr.(type) {
		case *models.ParseUserTask:
			routingKey = "parse_user"
		case *models.SimpleTask:
			routingKey = t.Name
			taskPtr = map[string]string{} // empty map
		default:
			return errors.Errorf("trying to push undeclared type of task to amqp")
		}

		message, err := json.Marshal(taskPtr)
		if err != nil {
			return errors.New(err)
		}

		err = initChannel(connection)
		if err != nil {
			return errors.New(err)
		}

		err = waitTasksQueueForEmpty(connection)
		if err != nil {
			return err
		}

		err = amqpChannel.Publish(
			"parser_tasks",
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
			_ = Cleanup()
			_ = amqp_helper.Cleanup()
			logger.Log.Error(err)
		}
	}

	return errors.Errorf("Unable to connect to queue")
}

func waitTasksQueueForEmpty(connection *amqp.Connection) error {
	for true {
		parserTasksQueue, err := amqpChannel.QueueInspect("bitbucket.org/d3dev/parse_pikabu/parser_tasks")
		if err != nil {
			// TODO: fix
			if !strings.Contains(err.Error(), "NOT_FOUND - no queue") {
				return err
			} else {
				amqpChannel = nil
				_ = initChannel(connection)
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
