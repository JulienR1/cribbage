package deck_test

import (
	"testing"

	"github.com/julienr1/cribbage/internal/deck"
	"github.com/stretchr/testify/assert"
)

func TestDraw(t *testing.T) {
	d := deck.New()
	c, err := d.Draw()

	assert.NoError(t, err)
	expected, _ := deck.NewCard(1, deck.SPADES)
	assert.Equal(t, expected, c)
}

func TestEmptyDeck(t *testing.T) {
	d := deck.New()
	var previous deck.Card = 0

	for range 52 {
		current, err := d.Draw()
		assert.NoError(t, err)
		assert.NotEqual(t, previous, current)
		previous = current
	}

	_, err := d.Draw()
	assert.Error(t, err)
}

func TestDrawN(t *testing.T) {
	d := deck.New()

	cards := make([]deck.Card, 5)
	n, err := d.DrawN(5, cards)

	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Len(t, cards, 5)
	assert.NotEmpty(t, cards)
}

func TestDrawMore(t *testing.T) {
	d := deck.New()
	cards := make([]deck.Card, 52)

	n, err := d.DrawN(26, cards)
	assert.NoError(t, err)
	assert.Equal(t, 26, n)
	n, err = d.DrawN(40, cards)
	assert.Error(t, err)
	assert.Equal(t, 0, n)
}

func TestShuffle(t *testing.T) {
	decksAreIdentical := true

	for range 3 {
		d1 := deck.New()
		d2 := deck.New()

		d1.Shuffle()
		d2.Shuffle()

		for range 52 {
			c1, e1 := d1.Draw()
			c2, e2 := d2.Draw()

			assert.NoError(t, e1)
			assert.NoError(t, e2)

			if c1 != c2 {
				decksAreIdentical = false
			}
		}

		if decksAreIdentical == false {
			break
		}
	}

	assert.False(t, decksAreIdentical, "expected both decks to be different when shuffled, got the same 3 times in a row (unlikely).")
}
