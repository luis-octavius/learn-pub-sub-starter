package main

import (
	"fmt"
	"log"
	"time"

	"github.com/luis-octavius/learn-pub-sub-starter/internal/gamelogic"
	"github.com/luis-octavius/learn-pub-sub-starter/internal/pubsub"
	"github.com/luis-octavius/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

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

	fmt.Println("Connection successful")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal("Error calling ClientWelcome: ", err)
	}

	game := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+game.GetUsername(), routing.PauseKey, "transient", handlerPause(game))
	if err != nil {
		log.Fatal("Error - SubscribeJSON: ", err)
	}

	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+game.GetUsername(), routing.ArmyMovesPrefix+".*", "transient", handlerMove(game, ch))
	if err != nil {
		log.Fatal("Error - SubscribeJSON: ", err)
	}

	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix, routing.WarRecognitionsPrefix+".*", "durable", handlerWar(game, ch))
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
			mv, err := game.CommandMove(words)
			if err != nil {
				log.Println("Error in move command, possibly bad input: ", err)
				continue
			}

			err = pubsub.PublishJSON(ch, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+mv.Player.Username, mv)
			if err != nil {
				log.Println("Error publishing the move: ", err)
				continue
			}

			log.Println("Move published successfully")
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

func publishGameLog(publishCh *amqp.Channel, username, msg string) error {
	return pubsub.PublishGob(
		publishCh,
		routing.ExchangePerilTopic,
		routing.GameLogSlug+"."+username,
		routing.GameLog{
			Username:    username,
			CurrentTime: time.Now(),
			Message:     msg,
		},
	)
}
