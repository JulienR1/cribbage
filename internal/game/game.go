package game

import (
	"fmt"
	"log"
	"sync"

	"github.com/julienr1/cribbage/internal/assert"
	"github.com/julienr1/cribbage/internal/deck"
)

type Game struct {
	deck  *deck.Deck
	state GameState

	hands Hands
	crib  Crib
	extra deck.Card
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
		g.flipExtraCard()
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

	var wg sync.WaitGroup
	wg.Add(len(g.hands))

	log.Printf("Waiting for players to add cards (%d) to the crib.", count)
	for _, hand := range g.hands {
		go func() {
			hand.SendToCrib(count, &g.crib)
			wg.Done()
		}()
	}

	wg.Wait()
	log.Println("Crib is now:", g.crib.String())

	g.state = extra
}

func (g *Game) flipExtraCard() {
	c, err := g.deck.Draw()
	assert.AssertE(err)
	g.extra = c

	log.Println("Extra card is", c)

	if g.extra.Value() == deck.JACK {
		p := g.hands[len(g.hands)-1]
		p.Score(1)
		log.Println(*p, "scored a point.")
	}
	g.state = play
}
}
