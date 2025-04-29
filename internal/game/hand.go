package game

import (
	"context"
	"fmt"
	"strings"

	"github.com/julienr1/cribbage/internal/assert"
	"github.com/julienr1/cribbage/internal/deck"
)

type Cards []deck.Card

type Hand struct {
	Cards  Cards
	player chan []uint8

	handlers map[uint8](func(data []uint8) []uint8)
	points   uint8
}

type Hands []*Hand

type CardStack interface {
	AddCard(card deck.Card) error
}

func NewHand() *Hand {
	player := make(chan []uint8)
	handlers := make(map[uint8](func(data []uint8) []uint8))
	return &Hand{player: player, handlers: handlers}
}

func (h *Hand) On(opcode uint8, callback func(data []uint8) []uint8) {
	h.handlers[opcode] = callback
}

func (h *Hand) Listen(ctx context.Context) {
	for {
		select {
		case data := <-h.player:
			assert.Assert(len(data) >= 1, "expected message to contain at least an opcode")
			opcode := data[0]

			callback, ok := h.handlers[opcode]
			assert.Assertf(ok, "opcode (%d) is not defined as a handler", opcode)
			h.player <- callback(data[1:])
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (h *Hand) SendToCrib(count uint8, crib CardStack) {
	h.player <- []uint8{REQUEST_CRIB_CARD, count}

	data := <-h.player
	assert.Assert(len(data) == int(count), "expected response to REQUEST_CRIB_CARD to be []deck.Card converted to []uint8")

	for _, c := range data {
		err := crib.AddCard(deck.Card(c))
		assert.AssertE(err)
	}
}

func (h *Hand) Score(points uint8) {
	assert.Assert(points > 0, "expected to score some points")
	h.points += points
}

func (cards Cards) String() string {
	var cardStrs []string
	for _, c := range cards {
		cardStrs = append(cardStrs, c.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(cardStrs, " | "))
}

func (h Hands) String() string {
	var hands []string
	for _, p := range h {
		hands = append(hands, p.Cards.String())
	}
	return strings.Join(hands, "\n")
}
