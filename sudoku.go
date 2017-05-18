package main

import (
	"github.com/jkeene-NAN/sudoku/game"
	"log"
)

func main() {
	log.Print("commencing")
	var initialGame *game.Game
	initialGame = game.NewGame()
	gameState, iterations, err := game.PlayGame(initialGame, 100000000)
	if err != nil {
		log.Printf("error playing game: %v", err)
	} else {
		log.Printf("iterations: %d, gameState: %v", iterations, *gameState)
	}
	log.Print("done")
}

