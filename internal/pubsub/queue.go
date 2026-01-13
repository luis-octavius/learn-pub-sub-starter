package pubsub

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key, queueType string) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		log.Println("Error creating channel from connection", err)
		return nil, amqp.Queue{}, err
	}

	var durable bool
	var autoDelete bool
	var exclusive bool

	switch queueType {
	case "transient":
		durable = false
		autoDelete = true
		exclusive = true
	case "durable":
		durable = true
		autoDelete = false
		exclusive = false
	}

	queue, err := ch.QueueDeclare(queueName, durable, autoDelete, exclusive, false, nil)
	if err != nil {
		log.Fatal("Error creating the queue: ", err)
		return nil, amqp.Queue{}, err
	}

	ch.QueueBind(queueName, key, exchange, false, nil)
	return ch, queue, err
}
