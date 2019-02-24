package amqp_helper

import (
	"github.com/streadway/amqp"
)

func DeclareExchanges(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		"parser_results",
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = ch.ExchangeDeclare(
		"parser_tasks",
		"fanout",
		false,
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

func DeclareParserResultsQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"bitbucket.org/d3dev/parse_pikabu/parser_results",
		false,
		false,
		false,
		false,
		nil,
	)
}

func DeclareParserTasksQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"bitbucket.org/d3dev/parse_pikabu/parser_tasks",
		false,
		false,
		false,
		false,
		nil,
	)
}
