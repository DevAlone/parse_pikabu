package results_processor

import (
	"bitbucket.org/d3dev/parse_pikabu/config"
	"bitbucket.org/d3dev/parse_pikabu/logging"
	"github.com/streadway/amqp"
	"gogsweb.2-47.ru/d3dev/pikago"
	"time"
)

func Run() error {
	for true {
		err := startListener()
		// TODO: handle connection refused
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}
func startListener() error {
	logging.Log.Debug("connecting to amqp server...")
	conn, err := amqp.Dial(config.Settings.AMQPAddress)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
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
		return err
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
		return err
	}

	err = ch.QueueBind(
		q.Name,
		"",
		"parser_results",
		false,
		nil,
	)
	if err != nil {
		return err
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
		return err
	}

	logging.Log.Debug("start waiting for parser results")
	for message := range messages {
		logging.Log.Debugf("got parser result: %v", string(message.Body))
		err = processMessage(message)
		if err != nil {
			return err
		}
	}
	logging.Log.Debug("stop waiting for parser results")

	return nil
}

func processMessage(message amqp.Delivery) error {
	switch message.RoutingKey {
	case "user_profile":
		var resp struct {
			User *pikago.UserProfile `json:"user"`
		}
		err := pikago.JsonUnmarshal(message.Body, &resp)
		if err != nil {
			return err
		}

		err = processUserProfile(resp.User)
		if err != nil {
			return err
		}
	default:
		logging.Log.Warningf(
			"Unregistered result type \"%v\". Message: \"%v\"",
			message.RoutingKey,
			string(message.Body),
		)
	}

	return nil
}
