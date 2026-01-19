package main

import (
	"fmt"
	"log"

	"github.com/luis-octavius/learn-pub-sub-starter/internal/gamelogic"
	"github.com/luis-octavius/learn-pub-sub-starter/internal/pubsub"
	"github.com/luis-octavius/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	connectionStr := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connectionStr)
	if err != nil {
		log.Fatal("Error creating the amqp connection: ", err)
	}

	defer conn.Close()

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("Could not create channel: %v", err)
	}

	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+".*", "durable")
	if err != nil {
		log.Fatal("Error declaring and binding queue: ", err)
	}

	fmt.Println("Connection successful")

	gamelogic.PrintServerHelp()

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		firstWord := words[0]
		if firstWord == "pause" {
			log.Println("Sending pause...")
			err = pubsub.PublishJSON(publishCh, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			continue
		}

		if firstWord == "resume" {
			log.Println("Sending resume...")
			err = pubsub.PublishJSON(publishCh, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			continue
		}

		if firstWord == "quit" {
			log.Println("Sending exit...")
			break
		}

		log.Println("Did not understand the command...")

	}

	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan
	//
	// fmt.Printf("\nEnding connection and closing...\n")
}
