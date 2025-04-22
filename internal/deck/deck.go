package deck

import (
	"errors"
	"math/rand/v2"

	"github.com/julienr1/cribbage/internal/assert"
)

var EmptyDeckErr = errors.New("deck is empty")
var NotEnoughCardsErr = errors.New("deck does not contain enough cards")

type Deck struct {
	drawIndex int
	cards     [52]Card
}

func New() *Deck {
	var d Deck
	for i, color := range []Color{SPADES, CLUBS, DIAMONDS, HEARTS} {
		for value := range 13 {
			c, err := NewCard(uint8(value)+1, color)
			assert.AssertE(err)
			d.cards[i*13+value] = c
		}
	}
	return &d
}

func (d *Deck) Shuffle() {
	for i := range len(d.cards) {
		target := i + rand.IntN(len(d.cards)-i)
		d.cards[target], d.cards[i] = d.cards[i], d.cards[target]
	}
}

func (d *Deck) Draw() (Card, error) {
	if d.drawIndex >= len(d.cards) {
		return 0, EmptyDeckErr
	}

	c := d.cards[d.drawIndex]
	d.drawIndex++
	return c, nil
}

func (d *Deck) DrawN(n int, out []Card) (int, error) {
	if d.drawIndex+n >= len(d.cards) {
		return 0, NotEnoughCardsErr
	}

	for i := range n {
		c, err := d.Draw()
		assert.AssertE(err)
		out[i] = c
	}
	return n, nil
}
