package game

import (
	"fmt"
	"log"
	"slices"
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

	startingIndex int
	playingIndex  int
	lastWhoPlayed *Player
	count         uint8
	stack         []deck.Card
}

func New(players []*Player) *Game {
	return &Game{players: players}
}

func (g *Game) Next() bool {
	log.Println("Current game state:", g.state.String())

	switch g.state {
	case deal:
		g.startingIndex = (g.startingIndex + 1) % len(g.players)
		g.playingIndex = g.startingIndex
		g.lastWhoPlayed = nil

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
		g.score()
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

	for _, p := range g.players {
		p.OriginalHand = slices.Clone(p.Hand)
	}

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

	lastPlayerIndex := (g.playingIndex + len(g.players) - 1) % len(g.players)
	g.points(g.players[lastPlayerIndex], HisHeels(g.extra))

	g.state = play
}

func (g *Game) playNextCard() {
	g.playingIndex = (g.playingIndex + 1) % len(g.players)
	playing := g.players[g.playingIndex]

	log.Println(playing, "is playing")

	var played *deck.Card
	g.sync(func(p *Player) {
		if p == playing {
			played = p.PlayCard(g.count)
		} else {
			p.WatchPlayedCard()
		}
	})

	if played != nil {
		g.lastWhoPlayed = playing
		g.stack = append(g.stack, *played)
		g.count += played.Points()
		g.sync(func(p *Player) {
			p.UpdateCount(g.count, *played)
		})
		log.Println("Stack:", g.stack)
	}

	g.points(playing, CountIs(15, g.count))
	g.points(playing, CountIs(31, g.count))
	g.points(playing, TailgateSeries(g.stack))
	g.points(playing, TailgateRepetitions(g.stack))

	var lastPlayedPoints = LastPlayed(played, playing, g.lastWhoPlayed)
	g.points(playing, lastPlayedPoints)

	if g.count == 31 || lastPlayedPoints > 0 {
		g.count = 0
		g.lastWhoPlayed = nil
		g.stack = []deck.Card{}
	}

	// NOTE: bail early if some players still have cards to play
	for _, p := range g.players {
		if len(p.Hand) > 0 {
			return
		}
	}

	g.points(playing, LastPlayed(nil, playing, g.lastWhoPlayed))
	g.state = score
}

func (g *Game) score() {
	for _, p := range g.players {
		p.Hand = p.OriginalHand
	}

	for i := range len(g.players) {
		playerIndex := (g.startingIndex + 1 + i) % len(g.players)
		p := g.players[playerIndex]
		hand := append(p.Hand, g.extra)

		log.Println("Scoring", p, "w/ hand", hand)

		g.points(p, Fifteen(hand))
		g.points(p, AnyRepetitions(hand))
		g.points(p, AnySeriesLog(hand))

		g.points(p, Flush(p.Hand, g.extra, false))
		g.points(p, HisNobs(g.extra, p.Hand))
	}

	dealer := g.players[g.startingIndex]
	crib := append(g.crib.Cards[:], g.extra)

	log.Println("Scoring crib w/ hand", crib)

	g.points(dealer, Fifteen(crib))
	g.points(dealer, AnyRepetitions(crib))
	g.points(dealer, AnySeriesLog(crib))

	g.points(dealer, Flush(g.crib.Cards[:], g.extra, true))

	g.state = done
}

func (g *Game) points(p *Player, points uint8) {
	if points > 0 {
		p.Score(points)
		g.notify(SCORE_CHANGED)
	}
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
