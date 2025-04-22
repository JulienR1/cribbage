package game

import (
	"fmt"
	"strings"

	"github.com/julienr1/cribbage/internal/assert"
	"github.com/julienr1/cribbage/internal/deck"
)

type Cards []deck.Card

type Hand struct {
	Cards  Cards
	player chan any
}

type Hands []*Hand

type CardStack interface {
	AddCard(card deck.Card) error
}

func NewHand() *Hand {
	return &Hand{}
}

func (h *Hand) SendToCrib(count uint8, crib CardStack) {
	h.player <- []byte{REQUEST_CRIB_CARD, count}

	for range count {
		card, ok := (<-h.player).(deck.Card)
		assert.Assert(ok, "expected response to REQUEST_CRIB_CARD to be a deck.Card")

		err := crib.AddCard(card)
		assert.AssertE(err)
	}
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
