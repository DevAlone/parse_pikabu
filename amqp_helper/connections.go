package amqp_helper

import (
	"sync"

	cmap "github.com/orcaman/concurrent-map"
	"github.com/streadway/amqp"
)

var AmqpConnections = cmap.New()

func Cleanup() error {
	// TODO: complete
	// amqpConnectionsMutex.Lock()
	// defer amqpConnectionsMutex.Unlock()

	/*
		for address, conn := range amqpConnections {
			connection := conn.(*amqp.Connection)
			err := connection.Close()
			delete(amqpConnections, address)
			if err != nil {
				return err
			}
		}
	*/

	return nil
}

var getAMQPConnectionMutex sync.Mutex

func GetAMQPConnection(amqpAddress string) (*amqp.Connection, error) {
	getAMQPConnectionMutex.Lock()
	defer getAMQPConnectionMutex.Unlock()

	if connection, ok := AmqpConnections.Get(amqpAddress); ok {
		return connection.(*amqp.Connection), nil
	}

	connection, err := amqp.Dial(amqpAddress)
	if err != nil {
		return nil, err
	}
	AmqpConnections.Set(amqpAddress, connection)

	return connection, nil
}
