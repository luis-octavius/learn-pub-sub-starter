package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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

	fmt.Printf("Welcome %s", username)

	pauseUser := routing.PauseKey + "." + username

	pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, pauseUser, routing.PauseKey, "transient")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Printf("\nEnding connection and closing...\n")
}
