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
	"github.com/julienr1/cribbage/internal/web"
)

func main() {
	fmt.Println("cribbage")

	if len(os.Args) > 1 && os.Args[1] == "web" {
		fmt.Println("web mode")
		web.Run()
		return
	}

	var playerCount int = 2
	if len(os.Args) > 1 {
		count, err := strconv.Atoi(os.Args[1])
		assert.AssertE(err)
		playerCount = count
	}

	var players game.Players
	for range playerCount {
		player := game.NewPlayer(players)
		players = append(players, player)
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

			return response
		})

		p.On(game.FLIP_EXTRA, func(data []uint8) []uint8 {
			assert.Assert(len(data) == 1, "expected FLIP_EXTRA to receive a deck.Card as uint8")
			card := deck.Card(data[0])
			log.Println("FLIP_EXTRA: card is", card.String())
			return []uint8{}
		})

		p.On(game.SCORE_CHANGED, func(_ []uint8) []uint8 {
			var scores []uint8
			for _, p2 := range players {
				scores = append(scores, p2.Points)
			}

			log.Println("SCORE_CHANGED", scores)
			return []uint8{}
		})

		p.On(game.REQUEST_PLAY_CARD, func(data []uint8) []uint8 {
			assert.Assert(len(data) == 1, "expected REQUEST_PLAY_CARD to pass along the current count")

			count := data[0]
			playable := p.Hand.Playable(count)

			log.Println("Count is", count, ", playable cards are:", playable)

			// TODO: some correct implementation
			if len(playable) > 0 {
				return []uint8{uint8(playable[0])}
			}
			return []uint8{}
		})

		p.On(game.WAIT_FOR_PLAY_CARD, func(data []uint8) []uint8 {
			log.Println(p, "is waiting")
			return []uint8{}
		})

		p.On(game.UPDATE_COUNT, func(data []uint8) []uint8 {
			assert.Assert(len(data) == 2, "expected UPDATE_COUNT to receive [count, deck.Card]")
			count := data[0]
			card := deck.Card(data[1])
			log.Println(card, "was played, count is now", count)
			return []uint8{}
		})
	}

	ctx, cancel := context.WithCancel(context.Background())
	for _, p := range players {
		go p.Listen(ctx)
	}

	game := game.New(players)
	for game.Next() {
	}
	cancel()
}
