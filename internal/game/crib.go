package game

import (
	"errors"

	"github.com/julienr1/cribbage/internal/deck"
)

var FullCribErr = errors.New("crib is full")

type Crib struct {
	Cards [4]deck.Card
	count int
}

func (c *Crib) AddCard(card deck.Card) error {
	if c.count >= len(c.Cards) {
		return FullCribErr
	}

	c.Cards[c.count] = card
	c.count++
	return nil
}

func (c *Crib) String() string {
	return Cards(c.Cards[:]).String()
}
