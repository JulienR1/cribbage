package main

import (
	"context"
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

	var players game.Hands
	for range playerCount {
		players = append(players, game.NewHand())
	}

	for _, p := range players {
		p.On(game.REQUEST_CRIB_CARD, func(data []uint8) []uint8 {
			assert.Assert(len(data) == 1, "expected REQUEST_CRIB_CARD to specify how many cards to give")
			count := data[0]

			// TODO: some correct implementation

			// TODO: remove this implementation
			var response []uint8
			for i := range count {
				response = append(response, uint8(p.Cards[i]))
			}
			p.Cards = p.Cards[count:]

			return response
		})
	}

	ctx, cancel := context.WithCancel(context.Background())
	for _, p := range players {
		go p.Listen(ctx)
	}

	game := game.New(players)

	log.Println("--- HANDS ---")
	game.Next()
	log.Println("\n", players.String())

	log.Println("--- CRIB ---")
	game.Next()

	log.Println("--- EXTRA ---")
	game.Next()
	cancel()
}
