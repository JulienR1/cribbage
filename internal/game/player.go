package game

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/julienr1/cribbage/internal/assert"
	"github.com/julienr1/cribbage/internal/deck"
)

type Player struct {
	Id   string
	Name string

	ch       chan []uint8
	handlers map[uint8](func(data []uint8) []uint8)

	Hand         Hand
	OriginalHand Hand

	Points uint8
}

type Players []*Player

func (players Players) Contains(id string) bool {
	for _, p := range players {
		if p.Id == id {
			return true
		}
	}
	return false
}

func NewPlayer(id string) *Player {
	ch := make(chan []uint8)
	handlers := make(map[uint8](func(data []uint8) []uint8))
	return &Player{Id: id, ch: ch, handlers: handlers}
}

func (p *Player) On(opcode uint8, callback func(data []uint8) []uint8) {
	p.handlers[opcode] = callback
}

func (p *Player) Listen(ctx context.Context) {
	for {
		select {
		case data := <-p.ch:
			assert.Assert(len(data) >= 1, "expected message to contain at least an opcode")
			opcode := data[0]

			callback, ok := p.handlers[opcode]
			assert.Assertf(ok, "opcode (%d) is not defined as a handler", opcode)
			p.ch <- callback(data[1:])
		case <-ctx.Done():
			return
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (p *Player) SendToCrib(count uint8, crib CardStack) {
	p.ch <- []uint8{REQUEST_CRIB_CARD, count}

	data := <-p.ch
	assert.Assert(len(data) == int(count), "expected response to REQUEST_CRIB_CARD to be []deck.Card converted to []uint8")

	for _, c := range data {
		index := slices.Index(p.Hand, deck.Card(c))
		assert.Assert(index >= 0, "expected card to be found in hand")
		p.Hand = slices.Delete(p.Hand, index, index+1)

		err := crib.AddCard(deck.Card(c))
		assert.AssertE(err)
	}
}

func (p *Player) SeeExtra(extra deck.Card) {
	p.shout([]uint8{FLIP_EXTRA, uint8(extra)}, "FLIP_EXTRA")
}

func (p *Player) PlayCard(count uint8) *deck.Card {
	p.ch <- []uint8{REQUEST_PLAY_CARD, count}
	data := <-p.ch

	playable := p.Hand.Playable(count)
	if len(playable) == 0 {
		return nil
	}

	assert.Assert(len(data) == 1, "expected REQUEST_PLAY_CARD to be answered with a deck.Card")
	played := deck.Card(data[0])

	index := slices.Index(p.Hand, played)
	assert.Assert(index >= 0, "expected played card to be in hand")
	p.Hand = slices.Delete(p.Hand, index, index+1)

	return &played
}

func (p *Player) WatchPlayedCard() {
	p.shout([]uint8{WAIT_FOR_PLAY_CARD}, "WAIT_FOR_PLAY_CARD ")
}

func (p *Player) UpdateCount(count uint8, played deck.Card) {
	p.shout([]uint8{UPDATE_COUNT, count, uint8(played)}, "UPDATE_COUNT")
}

func (p *Player) Score(points uint8) {
	assert.Assert(points > 0, "expected to score some points")
	p.Points += points
}

func (p *Player) shout(payload []uint8, label string) {
	p.ch <- payload
	data := <-p.ch
	assert.Assertf(len(data) == 0, "expected response to %s to be of length 0", label)
}

func (p *Player) String() string {
	return fmt.Sprintf("Player %s", p.Id)
}
