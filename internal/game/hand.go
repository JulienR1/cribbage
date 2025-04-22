package game

import (
	"github.com/julienr1/cribbage/internal/assert"
	"github.com/julienr1/cribbage/internal/deck"
)

type Hand struct {
	Cards []deck.Card

	player chan any
}

func NewHand() *Hand {
	return &Hand{}
}

func (h *Hand) SelectCribCard() deck.Card {
	h.player <- []byte{REQUEST_CRIB_CARD}

	var c = <-h.player
	card, ok := c.(deck.Card)
	assert.Assert(ok, "expected response to REQUEST_CRIB_CARD to be a deck.Card")
	return card
}
