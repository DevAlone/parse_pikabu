package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

func failOnError(err error, message string) {
	if err != nil {
		print(message, "\n")
		panic(err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "unable to connect")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "unable to create a channel")
	defer ch.Close()

	/*err = ch.ExchangeDeclare(
		"parser_results",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "unable to create exchange")*/

	if len(os.Args) < 2 {
		panic("too few arguments")
	}
	q, err := ch.QueueDeclare(
		os.Args[1],
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "unable to declare a queue")

	err = ch.QueueBind(
		q.Name,
		"",
		"parser_results",
		false,
		nil,
	)
	failOnError(err, "failed to bind queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "unable to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for items. To exit press Ctrl+C")
	<-forever
}
