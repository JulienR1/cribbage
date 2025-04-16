package deck_test

import (
	"testing"

	"github.com/julienr1/cribbage/internal/deck"
	"github.com/stretchr/testify/assert"
)

func TestCardCreation(t *testing.T) {
	tests := []struct {
		value    uint8
		color    deck.Color
		expected uint8
	}{
		{value: 2, color: deck.SPADES, expected: 0b00000010},
		{value: 2, color: deck.CLUBS, expected: 0b01000010},
		{value: 2, color: deck.DIAMONDS, expected: 0b10000010},
		{value: 2, color: deck.HEARTS, expected: 0b11000010},
		{value: 12, color: deck.HEARTS, expected: 0b11001100},
	}

	for i, test := range tests {
		c, err := deck.NewCard(test.value, test.color)
		assert.NoError(t, err)
		assert.Equalf(t, test.expected, uint8(c), "tests[%d]: expected card to be '%d', got '%d'.", i, test.expected, uint8(c))
		assert.Equalf(t, test.value, c.Value(), "tests[%d]: expected card have value '%d', got '%d'.", i, test.value, c.Value())
		assert.Equalf(t, test.color, c.Color(), "tests[%d]: expected card have color '%d', got '%d'.", i, test.color, c.Color())
	}
}

func TestInvalidCardCreation(t *testing.T) {
	tests := []struct {
		value uint8
		color deck.Color
	}{
		{value: 0, color: deck.SPADES},
		{value: 14, color: deck.SPADES},
		{value: 15, color: deck.SPADES},
	}
	for i := 4; i <= 15; i++ {
		tests = append(tests, struct {
			value uint8
			color deck.Color
		}{value: 3, color: deck.Color(i)})
	}

	for _, test := range tests {
		_, err := deck.NewCard(test.value, test.color)
		assert.Error(t, err)
	}
}

func TestCardString(t *testing.T) {
	tests := []struct {
		value    uint8
		color    deck.Color
		expected string
	}{
		{value: 2, color: deck.SPADES, expected: "2 of spades"},
		{value: 1, color: deck.SPADES, expected: "A of spades"},
		{value: 10, color: deck.SPADES, expected: "10 of spades"},
		{value: 11, color: deck.CLUBS, expected: "J of clubs"},
		{value: 12, color: deck.DIAMONDS, expected: "Q of diamonds"},
		{value: 13, color: deck.HEARTS, expected: "K of hearts"},
	}

	for _, test := range tests {
		c, err := deck.NewCard(test.value, test.color)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, c.String())
	}
}
