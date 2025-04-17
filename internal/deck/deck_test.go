package deck_test

import (
	"testing"

	"github.com/julienr1/cribbage/internal/deck"
	"github.com/stretchr/testify/assert"
)

func TestShuffle(t *testing.T) {
	decksAreIdentical := true

	for range 3 {
		d1 := deck.New().Shuffle()
		d2 := deck.New().Shuffle()

		for i := range len(d1) {
			if d1[i] != d2[i] {
				decksAreIdentical = false
			}
		}

		if decksAreIdentical == false {
			break
		}
	}

	assert.False(t, decksAreIdentical, "expected both decks to be different when shuffled, got the same 3 times in a row (unlikely).")
}
