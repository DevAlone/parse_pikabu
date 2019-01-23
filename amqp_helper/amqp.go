package amqp_helper

import (
	"github.com/streadway/amqp"
	"sync"
)

var amqpConnections = map[string]*amqp.Connection{}
var amqpConnectionsMutex = sync.RWMutex{}

func Cleanup() error {
	amqpConnectionsMutex.Lock()
	defer amqpConnectionsMutex.Unlock()

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
	amqpConnectionsMutex.RLock()
	if connection, ok := amqpConnections[amqpAddress]; ok {
		amqpConnectionsMutex.RUnlock()
		return connection, nil
	}
	amqpConnectionsMutex.RUnlock()

	connection, err := amqp.Dial(amqpAddress)
	if err != nil {
		return nil, err
	}
	amqpConnectionsMutex.Lock()
	amqpConnections[amqpAddress] = connection
	amqpConnectionsMutex.Unlock()

	return connection, nil
}
