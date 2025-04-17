package deck

import (
	"math/rand/v2"

	"github.com/julienr1/cribbage/internal/assert"
)

type Deck [52]Card

func New() Deck {
	var d Deck
	for i, color := range []Color{SPADES, CLUBS, DIAMONDS, HEARTS} {
		for value := range 13 {
			c, err := NewCard(uint8(value)+1, color)
			assert.AssertE(err)
			d[i*13+value] = c
		}
	}
	return d
}

func (d Deck) Shuffle() Deck {
	for i := range len(d) {
		target := i + rand.IntN(len(d)-i)
		d[target], d[i] = d[i], d[target]
	}
	return d
}
