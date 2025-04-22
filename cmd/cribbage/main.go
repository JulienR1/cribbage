package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/julienr1/cribbage/internal/assert"
	"github.com/julienr1/cribbage/internal/game"
)

func main() {
	fmt.Println("cribbage")

	var playerCount int = 2
	if len(os.Args) > 1 {
		count, err := strconv.Atoi(os.Args[1])
		assert.AssertE(err)
		playerCount = count
	}

	players := make(game.Hands, playerCount)
	for i := range playerCount {
		players[i] = game.NewHand()
	}

	game := game.New(players)

	log.Println("--- HANDS ---")
	game.Next()
	log.Println("\n", players.String())

	log.Println("--- CRIB ---")
	game.Next()

}
