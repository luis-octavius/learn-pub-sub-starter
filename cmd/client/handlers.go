package main

import (
	"fmt"
	"log"

	"github.com/luis-octavius/learn-pub-sub-starter/internal/gamelogic"
	"github.com/luis-octavius/learn-pub-sub-starter/internal/pubsub"
	"github.com/luis-octavius/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.Acktype {
	return func(p routing.PlayingState) pubsub.Acktype {
		defer fmt.Printf("> ")
		gs.HandlePause(p)
		return pubsub.Ack
	}
}

func handlerMove(gs *gamelogic.GameState, publishCh *amqp.Channel) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(move gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Printf("> ")
		moveOutcome := gs.HandleMove(move)

		switch moveOutcome {
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			fmt.Println("WARRR")
			key := routing.WarRecognitionsPrefix + "." + gs.GetUsername()

			err := pubsub.PublishJSON(publishCh, routing.ExchangePerilTopic, key, gamelogic.RecognitionOfWar{
				Attacker: move.Player,
				Defender: gs.GetPlayerSnap(),
			})
			if err != nil {
				log.Println("PublishJSON - error: ", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NackDiscard
		}
		return pubsub.NackDiscard
	}
}

func handlerWar(gs *gamelogic.GameState, publishCh *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(w gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Printf("> ")

		warOutcome, winner, loser := gs.HandleWar(w)

		switch warOutcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			err := publishGameLog(
				publishCh,
				gs.GetUsername(),
				fmt.Sprintf("%s won a war against %s", winner, loser),
			)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeYouWon:
			err := publishGameLog(
				publishCh,
				gs.GetUsername(),
				fmt.Sprintf("%s won a war against %s", winner, loser),
			)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			err := publishGameLog(
				publishCh,
				gs.GetUsername(),
				fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser),
			)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		}

		log.Println("Error - War Outcome not fit the cases")
		return pubsub.NackDiscard
	}
}
