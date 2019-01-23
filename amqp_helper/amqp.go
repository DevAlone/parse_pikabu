package amqp_helper

import (
	"github.com/streadway/amqp"
)

var amqpConnections = map[string]*amqp.Connection{}

func Cleanup() error {
	for address, connection := range amqpConnections {
		err := connection.Close()
		delete(amqpConnections, address)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetAMQPConnection(amqpAddress string) (*amqp.Connection, error) {
	if connection, ok := amqpConnections[amqpAddress]; ok {
		return connection, nil
	}

	connection, err := amqp.Dial(amqpAddress)
	if err != nil {
		return nil, err
	}
	amqpConnections[amqpAddress] = connection

	return connection, nil
}
