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

	players []*Player
	crib    Crib
	extra   deck.Card
}

func New(players []*Player) *Game {
	return &Game{players: players}
}

func (g *Game) Next() bool {
	log.Println("Current game state:", g.state.String())

	switch g.state {
	case deal:
		g.deck = deck.New()
		g.deck.Shuffle()
		g.deal()
	case crib:
		g.buildCrib()
	case extra:
		g.flipExtraCard()
	case play:
		g.playNextCard()
	case score:
	}

	return g.state != done
}

func (g *Game) deal() {
	var handSize = 6
	if len(g.players) > 2 {
		handSize = 5
	}

	log.Printf("Dealing %d cards in a %d player game.\n", handSize, len(g.players))

	for i := range g.players {
		g.players[i].Hand = make([]deck.Card, handSize)
		n, err := g.deck.DrawN(handSize, g.players[i].Hand)
		assert.AssertE(err)
		assert.Assert(n == handSize, fmt.Sprintf("expected deck to draw %d cards.", handSize))
	}

	// If this is a 3 player game, the crib should get a single card when dealing the hands.
	// The 3 missing cards will be coming from the players.
	if len(g.players) == 3 {
		c, err := g.deck.Draw()
		assert.AssertE(err)

		err = g.crib.AddCard(c)
		assert.AssertE(err)
	}

	g.notify(RECEIVE_HAND)
	g.state = crib
}

func (g *Game) buildCrib() {
	var count uint8 = 2
	if len(g.players) > 2 {
		count = 1
	}

	log.Printf("Waiting for players to add cards (%d) to the crib.", count)
	g.sync(func(player *Player) {
		player.SendToCrib(count, &g.crib)
	})
	log.Println("Crib is now:", g.crib.String())

	g.state = extra
}

func (g *Game) flipExtraCard() {
	c, err := g.deck.Draw()
	assert.AssertE(err)
	g.extra = c

	log.Println("Extra card is", c)

	g.sync(func(player *Player) {
		player.SeeExtra(g.extra)
	})

	if g.extra.Value() == deck.JACK {
		p := g.players[len(g.players)-1]
		g.points(p, 1)
		log.Println(*p, "scored a point.")
	}

	g.state = play
}

func (g *Game) playNextCard() {
	panic("not implemented")
}

func (g *Game) points(p *Player, points uint8) {
	p.Score(points)
	g.notify(SCORE_CHANGED)
}

func (g *Game) sync(callback func(p *Player)) {
	var wg sync.WaitGroup
	wg.Add(len(g.players))

	for _, player := range g.players {
		go func() {
			callback(player)
			wg.Done()
		}()
	}

	wg.Wait()
}

func (g *Game) notify(opcode uint8) {
	g.sync(func(p *Player) {
		p.ch <- []uint8{opcode}
		<-p.ch
	})
}
