package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int

const (
	Ack Acktype = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key, queueType string, handler func(T) Acktype) error {
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

			acktype := handler(data)
			switch acktype {
			case Ack:
				msg.Ack(false)
				log.Println("Ack")
			case NackRequeue:
				msg.Nack(false, true)
				log.Println("NackRequeue")
			case NackDiscard:
				msg.Nack(false, false)
				log.Println("NackDiscard")
			}
		}
	}()

	return nil
}
