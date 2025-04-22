package game

import (
	"fmt"
	"log"

	"github.com/julienr1/cribbage/internal/assert"
	"github.com/julienr1/cribbage/internal/deck"
)

type Game struct {
	deck  *deck.Deck
	state GameState

	hands Hands
	crib  Crib
}

func New(players Hands) *Game {
	d := deck.New()
	return &Game{deck: d, hands: players}
}

func (g *Game) Next() {
	log.Println("Current game state:", g.state.String())

	switch g.state {
	case deal:
		g.deck.Shuffle()
		g.deal()
	case crib:
		g.buildCrib()
	case extra:
	case play:
	case score:
	}
}

func (g *Game) deal() {
	var handSize = 6
	if len(g.hands) > 2 {
		handSize = 5
	}

	log.Printf("Dealing %d cards in a %d player game.\n", handSize, len(g.hands))

	for i := range g.hands {
		g.hands[i].Cards = make([]deck.Card, handSize)
		n, err := g.deck.DrawN(handSize, g.hands[i].Cards)
		assert.AssertE(err)
		assert.Assert(n == handSize, fmt.Sprintf("expected deck to draw %d cards.", handSize))
	}

	// If this is a 3 player game, the crib should get a single card when dealing the hands.
	// The 3 missing cards will be coming from the players.
	if len(g.hands) == 3 {
		c, err := g.deck.Draw()
		assert.AssertE(err)

		err = g.crib.AddCard(c)
		assert.AssertE(err)
	}

	g.state = crib
}

func (g *Game) buildCrib() {
	var count uint8 = 2
	if len(g.hands) > 2 {
		count = 1
	}

	log.Printf("Waiting for players to add cards (%d) to the crib.", count)
	for _, hand := range g.hands {
		hand.SendToCrib(count, &g.crib)
	}
	log.Println("Crib is now:", g.crib.String())

	g.state = extra
}

