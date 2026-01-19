package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	marshaled, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("Error marshaling data: %v", err)
	}

	ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        marshaled,
	})

	return nil
}

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(val)
	if err != nil {
		return fmt.Errorf("Error encoding data to gob: %v", err)
	}

	ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/gob",
		Body:        buf.Bytes(),
	})

	return nil
}
