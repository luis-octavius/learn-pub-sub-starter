package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key, queueType string, handler func(T)) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return fmt.Errorf("")
	}

	deliveryChannel, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Error creating a channel of delivery: %v", err)
	}

	go func() {
		for msg := range deliveryChannel {
			var data T
			err := json.Unmarshal(msg.Body, &data)
			if err != nil {
				log.Fatalf("Error unmarshaling JSON: %v\n", err)
			}

			handler(data)
			msg.Ack(false)
		}
	}()

	return nil
}
