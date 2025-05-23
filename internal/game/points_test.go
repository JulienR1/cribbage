package game_test

import (
	"testing"

	"github.com/julienr1/cribbage/internal/deck"
	"github.com/julienr1/cribbage/internal/game"
	"github.com/stretchr/testify/assert"
)

type ComboTest struct {
	cards  []deck.Card
	points uint8
}

func TestCountIs(t *testing.T) {
	var sum uint8 = 0
	for i := range uint8(32) {
		sum += game.CountIs(15, i)
		sum += game.CountIs(31, i)
	}

	assert.Equal(t, uint8(4), sum)
	assert.Equal(t, uint8(2), game.CountIs(15, 15))
	assert.Equal(t, uint8(2), game.CountIs(31, 31))
}

func c(value uint8) deck.Card {
	card, _ := deck.NewCard(value, deck.CLUBS)
	return card
}

func TestFifteen(t *testing.T) {
	tests := []ComboTest{
		{cards: []deck.Card{}, points: 0},
		{cards: []deck.Card{c(deck.ACE), c(2), c(3), c(4)}, points: 0},
		{cards: []deck.Card{c(deck.ACE), c(2), c(10), c(5)}, points: 2},
		{cards: []deck.Card{c(deck.ACE), c(deck.JACK), c(10), c(5)}, points: 4},
		{cards: []deck.Card{c(deck.ACE), c(4), c(10), c(5)}, points: 4},
		{cards: []deck.Card{c(10), c(deck.JACK), c(deck.QUEEN), c(deck.KING), c(5)}, points: 8},
		{cards: []deck.Card{c(deck.ACE), c(4), c(10), c(deck.JACK)}, points: 4},
		{cards: []deck.Card{c(8), c(7), c(3), c(4)}, points: 4},
		{cards: []deck.Card{c(8), c(7), c(3), c(4), c(8)}, points: 8},
		{cards: []deck.Card{c(8), c(7), c(3), c(4), c(4)}, points: 8},
		{cards: []deck.Card{c(8), c(7), c(3), c(3), c(4)}, points: 6},
	}

	for i, test := range tests {
		points := game.Fifteen(test.cards)
		assert.Equalf(t, test.points, points, "tests[%d]", i)
	}
}

func TestLastPlayed(t *testing.T) {
	p1, p2 := game.NewPlayer("1"), game.NewPlayer("2")
	someCard, _ := deck.NewCard(deck.ACE, deck.CLUBS)

	tests := []struct {
		card          *deck.Card
		playing       *game.Player
		lastWhoPlayed *game.Player
		points        uint8
	}{
		{card: &someCard, playing: p1, lastWhoPlayed: p2, points: 0},
		{card: &someCard, playing: p1, lastWhoPlayed: p1, points: 0},
		{card: nil, playing: p1, lastWhoPlayed: p2, points: 0},
		{card: nil, playing: p2, lastWhoPlayed: p2, points: 1},
	}

	for _, test := range tests {
		points := game.LastPlayed(test.card, test.playing, test.lastWhoPlayed)
		assert.Equal(t, test.points, points)
	}
}

func TestTailgateSeries(t *testing.T) {
	c1, _ := deck.NewCard(deck.ACE, deck.CLUBS)
	c2, _ := deck.NewCard(2, deck.CLUBS)
	c3, _ := deck.NewCard(3, deck.HEARTS)
	c4, _ := deck.NewCard(4, deck.DIAMONDS)
	c5, _ := deck.NewCard(5, deck.SPADES)
	c6, _ := deck.NewCard(6, deck.CLUBS)
	c7, _ := deck.NewCard(7, deck.CLUBS)

	tests := []ComboTest{
		{cards: []deck.Card{c1}, points: 0},
		{cards: []deck.Card{c1, c2}, points: 0},
		{cards: []deck.Card{c2, c1}, points: 0},
		{cards: []deck.Card{c4, c2}, points: 0},
		{cards: []deck.Card{c1, c2, c3}, points: 3},
		{cards: []deck.Card{c1, c3, c2}, points: 3},
		{cards: []deck.Card{c2, c1, c3}, points: 3},
		{cards: []deck.Card{c2, c3, c1}, points: 3},
		{cards: []deck.Card{c3, c2, c1}, points: 3},
		{cards: []deck.Card{c1, c2, c3, c4}, points: 4},
		{cards: []deck.Card{c1, c2, c3, c4, c5}, points: 5},
		{cards: []deck.Card{c1, c2, c3, c4, c5, c6}, points: 6},
		{cards: []deck.Card{c1, c2, c3, c4, c5, c6, c7}, points: 7},
		{cards: []deck.Card{c3, c1, c4, c7, c5, c2, c6}, points: 7},
		{cards: []deck.Card{c1, c3, c4, c5, c7}, points: 0},
		{cards: []deck.Card{c1, c3, c5, c4, c7}, points: 0},
		{cards: []deck.Card{c3, c1, c5, c4, c7}, points: 0},
		{cards: []deck.Card{c3, c1, c5, c4}, points: 0},
		{cards: []deck.Card{c1, c3, c5, c4}, points: 3},
	}

	for i, test := range tests {
		points := game.TailgateSeries(test.cards)
		assert.Equalf(t, test.points, points, "tests[%d]", i)
	}
}

func TestAnySeries(t *testing.T) {
	c1, _ := deck.NewCard(deck.ACE, deck.CLUBS)
	c2, _ := deck.NewCard(2, deck.CLUBS)
	c3, _ := deck.NewCard(3, deck.HEARTS)
	c4, _ := deck.NewCard(4, deck.DIAMONDS)
	c5, _ := deck.NewCard(5, deck.SPADES)
	c6, _ := deck.NewCard(6, deck.CLUBS)
	c7, _ := deck.NewCard(7, deck.CLUBS)

	tests := []ComboTest{
		{cards: []deck.Card{c1}, points: 0},
		{cards: []deck.Card{c1, c2}, points: 0},
		{cards: []deck.Card{c2, c1}, points: 0},
		{cards: []deck.Card{c4, c2}, points: 0},
		{cards: []deck.Card{c1, c2, c3}, points: 3},
		{cards: []deck.Card{c1, c3, c2}, points: 3},
		{cards: []deck.Card{c2, c1, c3}, points: 3},
		{cards: []deck.Card{c2, c3, c1}, points: 3},
		{cards: []deck.Card{c3, c2, c1}, points: 3},
		{cards: []deck.Card{c1, c2, c3, c4}, points: 4},
		{cards: []deck.Card{c1, c2, c3, c4, c5}, points: 5},
		{cards: []deck.Card{c1, c2, c3, c4, c5, c6}, points: 6},
		{cards: []deck.Card{c1, c2, c3, c4, c5, c6, c7}, points: 7},
		{cards: []deck.Card{c3, c1, c4, c7, c5, c2, c6}, points: 7},
	}

	for i, test := range tests {
		points := game.AnySeries(test.cards)
		assert.Equalf(t, test.points, points, "tests[%d]", i)
	}
}

func TestTailgateRepetitions(t *testing.T) {
	c1, _ := deck.NewCard(2, deck.CLUBS)
	c2, _ := deck.NewCard(2, deck.DIAMONDS)
	c3, _ := deck.NewCard(2, deck.HEARTS)
	c4, _ := deck.NewCard(2, deck.SPADES)
	c5, _ := deck.NewCard(3, deck.CLUBS)

	tests := []ComboTest{
		{cards: []deck.Card{c1, c5}, points: 0},
		{cards: []deck.Card{c1, c2, c5}, points: 0},
		{cards: []deck.Card{c1, c2, c3, c5}, points: 0},
		{cards: []deck.Card{c1, c2, c3, c4, c5}, points: 0},
		{cards: []deck.Card{c1, c2}, points: 2},
		{cards: []deck.Card{c1, c2, c3}, points: 6},
		{cards: []deck.Card{c1, c2, c3, c4}, points: 12},
		{cards: []deck.Card{c4, c3, c2, c1}, points: 12},
	}

	for i, test := range tests {
		points := game.TailgateRepetitions(test.cards)
		assert.Equalf(t, test.points, points, "tests[%d]", i)
	}
}

func TestAnyRepetitions(t *testing.T) {
	c1, _ := deck.NewCard(2, deck.CLUBS)
	c2, _ := deck.NewCard(2, deck.DIAMONDS)
	c3, _ := deck.NewCard(2, deck.HEARTS)
	c4, _ := deck.NewCard(2, deck.SPADES)
	c5, _ := deck.NewCard(3, deck.CLUBS)
	c6, _ := deck.NewCard(3, deck.SPADES)
	c7, _ := deck.NewCard(3, deck.DIAMONDS)

	tests := []ComboTest{
		{cards: []deck.Card{c1, c5}, points: 0},
		{cards: []deck.Card{c1, c2, c5}, points: 2},
		{cards: []deck.Card{c1, c2, c3, c5}, points: 6},
		{cards: []deck.Card{c1, c2, c3, c4, c5}, points: 12},
		{cards: []deck.Card{c1, c2}, points: 2},
		{cards: []deck.Card{c1, c2, c3}, points: 6},
		{cards: []deck.Card{c1, c2, c3, c4}, points: 12},
		{cards: []deck.Card{c4, c3, c2, c1}, points: 12},
		{cards: []deck.Card{c1, c3, c6, c7}, points: 4},
		{cards: []deck.Card{c5, c1, c3, c6, c7}, points: 8},
	}

	for i, test := range tests {
		points := game.AnyRepetitions(test.cards)
		assert.Equalf(t, test.points, points, "tests[%d]", i)
	}
}

func TestFlush(t *testing.T) {
	c1, _ := deck.NewCard(deck.ACE, deck.CLUBS)
	c2, _ := deck.NewCard(2, deck.CLUBS)
	c3, _ := deck.NewCard(3, deck.CLUBS)
	c4, _ := deck.NewCard(4, deck.CLUBS)
	c5, _ := deck.NewCard(5, deck.CLUBS)

	c6, _ := deck.NewCard(5, deck.DIAMONDS)
	c7, _ := deck.NewCard(5, deck.HEARTS)
	c8, _ := deck.NewCard(5, deck.SPADES)

	tests := []struct {
		cards  []deck.Card
		extra  deck.Card
		isCrib bool
		points uint8
	}{
		{cards: []deck.Card{c1, c2, c6, c7}, extra: c8, isCrib: false, points: 0},
		{cards: []deck.Card{c1, c2, c3, c7}, extra: c8, isCrib: false, points: 0},
		{cards: []deck.Card{c1, c2, c3, c4}, extra: c8, isCrib: false, points: 4},
		{cards: []deck.Card{c1, c2, c3, c4}, extra: c5, isCrib: false, points: 5},
		{cards: []deck.Card{c1, c2, c3, c6}, extra: c5, isCrib: false, points: 0},
		{cards: []deck.Card{c1, c2, c3, c6}, extra: c5, isCrib: true, points: 0},
		{cards: []deck.Card{c1, c2, c3, c4}, extra: c6, isCrib: true, points: 0},
		{cards: []deck.Card{c1, c2, c3, c4}, extra: c5, isCrib: true, points: 5},
	}

	for i, test := range tests {
		points := game.Flush(test.cards, test.extra, test.isCrib)
		assert.Equalf(t, test.points, points, "tests[%d]", i)
	}
}

func TestHisHeels(t *testing.T) {
	c1, _ := deck.NewCard(deck.JACK, deck.CLUBS)
	c2, _ := deck.NewCard(deck.JACK, deck.DIAMONDS)
	c3, _ := deck.NewCard(deck.JACK, deck.SPADES)
	c4, _ := deck.NewCard(deck.JACK, deck.HEARTS)
	c5, _ := deck.NewCard(3, deck.HEARTS)

	for _, c := range []deck.Card{c1, c2, c3, c4} {
		assert.Equal(t, uint8(1), game.HisHeels(c))
	}
	assert.Equal(t, uint8(0), game.HisHeels(c5))
}

func TestHisNobs(t *testing.T) {
	j1, _ := deck.NewCard(deck.JACK, deck.CLUBS)
	j2, _ := deck.NewCard(deck.JACK, deck.HEARTS)
	j3, _ := deck.NewCard(deck.JACK, deck.SPADES)
	j4, _ := deck.NewCard(deck.JACK, deck.DIAMONDS)

	c1, _ := deck.NewCard(2, deck.CLUBS)
	c2, _ := deck.NewCard(2, deck.DIAMONDS)
	c3, _ := deck.NewCard(2, deck.HEARTS)
	c4, _ := deck.NewCard(2, deck.SPADES)
	c5, _ := deck.NewCard(3, deck.SPADES)

	tests := []struct {
		hand   []deck.Card
		extra  deck.Card
		points uint8
	}{
		{hand: []deck.Card{j1, j2, j3, j4}, extra: c1, points: 1},
		{hand: []deck.Card{j1, j2, j3, j4}, extra: c2, points: 1},
		{hand: []deck.Card{j1, j2, j3, j4}, extra: c3, points: 1},
		{hand: []deck.Card{j1, j2, j3, j4}, extra: c4, points: 1},
		{hand: []deck.Card{c1, c2, c3, c4}, extra: j1, points: 0},
		{hand: []deck.Card{j1, c2, c3, c4}, extra: j2, points: 0},
		{hand: []deck.Card{j1, c2, c3, c4}, extra: c5, points: 0},
		{hand: []deck.Card{c1, c2, c3, c4}, extra: c5, points: 0},
		{hand: []deck.Card{c1, c2, j3, c4}, extra: c5, points: 1},
	}

	for i, test := range tests {
		points := game.HisNobs(test.extra, test.hand)
		assert.Equalf(t, test.points, points, "tests[%d]", i)
	}
}
