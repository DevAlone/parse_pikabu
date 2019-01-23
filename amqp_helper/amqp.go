package amqp_helper

import (
	"github.com/orcaman/concurrent-map"
	"github.com/streadway/amqp"
)

// TODO: fix
// var amqpConnections = map[string]*amqp.Connection{}
// var amqpConnectionsMutex = sync.RWMutex{}
var amqpConnections = cmap.New()

func Cleanup() error {
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

func GetAMQPConnection(amqpAddress string) (*amqp.Connection, error) {
	// amqpConnectionsMutex.RLock()
	if connection, ok := amqpConnections.Get(amqpAddress); ok {
		// amqpConnectionsMutex.RUnlock()
		return connection.(*amqp.Connection), nil
	}
	// amqpConnectionsMutex.RUnlock()

	connection, err := amqp.Dial(amqpAddress)
	if err != nil {
		return nil, err
	}
	// amqpConnectionsMutex.Lock()
	amqpConnections.Set(amqpAddress, connection)
	// fmt.Println("%v", amqpConnections)
	// amqpConnectionsMutex.Unlock()

	return connection, nil
}
