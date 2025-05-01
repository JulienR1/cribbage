package game

import (
	"context"
	"time"

	"github.com/julienr1/cribbage/internal/assert"
	"github.com/julienr1/cribbage/internal/deck"
)

var nextId int = 0

type Player struct {
	id       int
	ch       chan []uint8
	handlers map[uint8](func(data []uint8) []uint8)

	Hand   Hand
	points uint8
}

func NewPlayer() *Player {
	id := nextId
	nextId++

	ch := make(chan []uint8)
	handlers := make(map[uint8](func(data []uint8) []uint8))
	return &Player{id: id, ch: ch, handlers: handlers}
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
		err := crib.AddCard(deck.Card(c))
		assert.AssertE(err)
	}
}

func (p *Player) SeeExtra(extra deck.Card) {
	p.ch <- []uint8{FLIP_EXTRA, uint8(extra)}
	data := <-p.ch
	assert.Assert(len(data) == 0, "expected response to FLIP_EXTRA to be of length 0")
}

func (p *Player) Score(points uint8) {
	assert.Assert(points > 0, "expected to score some points")
	p.points += points
}
