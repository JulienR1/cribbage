package game_test

import (
	"testing"

	"github.com/julienr1/cribbage/internal/deck"
	"github.com/julienr1/cribbage/internal/game"
	"github.com/stretchr/testify/assert"
)

func TestAddCard(t *testing.T) {
	var crib game.Crib
	card, err := deck.NewCard(2, deck.SPADES)

	assert.NoError(t, err)
	for range 4 {
		err = crib.AddCard(card)
		assert.NoError(t, err)
	}
	assert.NotEmpty(t, crib.Cards)

	err = crib.AddCard(card)
	assert.Error(t, err)
}
