package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/julienr1/cribbage/internal/assert"
	"github.com/julienr1/cribbage/internal/deck"
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

	var players []*game.Player
	for range playerCount {
		players = append(players, game.NewPlayer())
	}

	for _, p := range players {
		p.On(game.RECEIVE_HAND, func(_ []uint8) []uint8 {
			log.Println("RECEIVE_HAND:", p.Hand)
			return []uint8{}
		})

		p.On(game.REQUEST_CRIB_CARD, func(data []uint8) []uint8 {
			assert.Assert(len(data) == 1, "expected REQUEST_CRIB_CARD to specify how many cards to give")
			count := data[0]

			// TODO: some correct implementation

			// TODO: remove this implementation
			var response []uint8
			for i := range count {
				response = append(response, uint8(p.Hand[i]))
			}
			p.Hand = p.Hand[count:]

			return response
		})

		p.On(game.FLIP_EXTRA, func(data []uint8) []uint8 {
			assert.Assert(len(data) == 1, "expected FLIP_EXTRA to receive a deck.Card as uint8")
			card := deck.Card(data[0])
			log.Println("FLIP_EXTRA: card is", card.String())
			return []uint8{}
		})

		p.On(game.SCORE_CHANGED, func(data []uint8) []uint8 {
			log.Println("SCORE_CHANGED", data)
			return []uint8{}
		})
	}

	ctx, cancel := context.WithCancel(context.Background())
	for _, p := range players {
		go p.Listen(ctx)
	}

	game := game.New(players)

	log.Println("--- HANDS ---")
	game.Next()

	log.Println("--- CRIB ---")
	game.Next()

	log.Println("--- EXTRA ---")
	game.Next()

	log.Println("--- PLAYING ---")

	cancel()
}
