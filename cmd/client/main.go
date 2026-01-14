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
	fmt.Println("Starting Peril client...")

	connectionStr := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connectionStr)
	if err != nil {
		log.Fatal("Error creating the amqp connection: ", err)
	}

	// ch, err := conn.Channel()
	// if err != nil {
	// 	log.Fatal("Error opening channel: ", err)
	// }

	defer conn.Close()

	fmt.Println("Connection successful")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal("Error calling ClientWelcome: ", err)
	}

	pauseUser := routing.PauseKey + "." + username

	game := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, pauseUser, routing.PauseKey, "transient", handlerPause(game))
	if err != nil {
		log.Fatal("Error - SubscribeJSON: ", err)
	}

	for {
		words := gamelogic.GetInput()

		firstWord := words[0]

		switch firstWord {
		case "spawn":
			err = game.CommandSpawn(words)
			if err != nil {
				log.Println("Error in spawn command, possibly bad input: ", err)
				continue
			}

			log.Println("Success spawning units")
			continue
		case "move":
			armyMove, err := game.CommandMove(words)
			if err != nil {
				log.Println("Error in move command, possibly bad input: ", err)
				continue
			}

			log.Println("Army move successful: ", armyMove)
			continue

		case "status":
			game.CommandStatus()
			continue
		case "help":
			gamelogic.PrintClientHelp()
			continue
		case "spam":
			log.Println("Spamming not allowed yet!")
			continue
		case "quit":
			gamelogic.PrintQuit()
		default:
			log.Println("Did not understand the command...")
			continue
		}

		break
	}

	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan
	//
	fmt.Printf("\nEnding connection and closing...\n")
}
