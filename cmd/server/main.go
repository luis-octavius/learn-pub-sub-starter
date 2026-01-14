package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	connectionStr := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connectionStr)
	if err != nil {
		log.Fatal("Error creating the amqp connection: ", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Error opening channel: ", err)
	}

	defer conn.Close()

	channel, gameLogs, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, "game_logs", routing.GameLogSlug, "durable")
	if err != nil {
		log.Fatal("Error declaring and binding queue: ", err)
	}

	fmt.Println("Channel: ", channel, "GameLogs: ", gameLogs)

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
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			continue
		}

		if firstWord == "resume" {
			log.Println("Sending resume...")
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
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
