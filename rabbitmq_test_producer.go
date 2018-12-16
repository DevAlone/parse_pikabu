package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"os"
)

func failOnError(err error, message string) {
	if err != nil {
		print(message, "\n")
		panic(err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "unable to connect")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "unable to open channel")
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
	failOnError(err, "unable to create exchange")

	message := map[string]string{
		"response": "1",
		"status":   "ok",
	}
	if len(os.Args) > 1 {
		message["data"] = os.Args[1]
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		failOnError(err, "unable to create a message")
	}

	// TODO: Add a listener with Channel.NotifyReturn to handle any undeliverable message when calling publish with either the mandatory or immediate parameters as true.
	err = ch.Publish(
		"parser_results",
		"user_profile",
		true,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         messageBytes,
		},
	)
	failOnError(err, "unable to publish the message")

	print("sent\n")
}
