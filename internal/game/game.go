package game

import (
	"fmt"

	"github.com/julienr1/cribbage/internal/assert"
	"github.com/julienr1/cribbage/internal/deck"
)

type Game struct {
	deck  *deck.Deck
	state GameState

	hands []*Hand
	crib  Crib
}

func New(players []*Hand) *Game {
	d := deck.New()
	return &Game{deck: d, hands: players}
}

func (g *Game) Next() {
	switch g.state {
	case deal:
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

	for i := range g.hands {
		n, err := g.deck.DrawN(handSize, g.hands[i].Cards)
		assert.Assert(n == handSize, fmt.Sprintf("expected deck to draw %d cards.", handSize))
		assert.AssertE(err)
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
	for _, hand := range g.hands {
		card := hand.SelectCribCard()
		err := g.crib.AddCard(card)
		assert.AssertE(err)
	}
	g.state = extra
}
