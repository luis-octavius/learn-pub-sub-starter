package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(p routing.PlayingState) {
		defer fmt.Printf("> ")
		gs.HandlePause(p)
	}
}

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) {
	return func(a gamelogic.ArmyMove) {
		defer fmt.Printf("> ")
		gs.HandleMove(a)
	}
}
