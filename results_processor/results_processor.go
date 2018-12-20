package results_processor

import (
	"bitbucket.org/d3dev/parse_pikabu/config"
	"bitbucket.org/d3dev/parse_pikabu/logger"
	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-errors/errors"
	"github.com/streadway/amqp"
	"gogsweb.2-47.ru/d3dev/pikago"
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
			return errors.New(err)
		}
	}
	logger.Log.Debug("stop waiting for parser results")

	return nil
}

func processMessage(message amqp.Delivery) error {
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

		err = processUserProfile(resp.ParsingTimestamp, resp.Results[0].User)
		if err != nil {
			return errors.New(err)
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
