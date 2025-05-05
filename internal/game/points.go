package game

import (
	"slices"

	"github.com/julienr1/cribbage/internal/assert"
	"github.com/julienr1/cribbage/internal/deck"
)

func Fifteen(count uint8) (points uint8) {
	if count == 15 {
		return 2
	}
	return 0
}

func ThirtyOne(count uint8) (points uint8) {
	if count == 31 {
		return 2
	}
	return 0
}

func LastPlayed(playedCard *deck.Card, playing, lastWhoPlayed *Player) (points uint8) {
	if playedCard == nil && playing == lastWhoPlayed {
		return 1
	}
	return 0
}

func TailgateSeries(cards []deck.Card) (points uint8) {
	var bestLength = 0
	for i := range len(cards) - 2 {
		length := i + 3
		subset := cards[len(cards)-length:]
		if int(AnySeries(subset)) == len(subset) {
			bestLength = len(subset)
		}
	}

	return uint8(bestLength)
}

func AnySeries(cards []deck.Card) (points uint8) {
	slices.SortFunc(cards, func(a, b deck.Card) int {
		return int(a.Value()) - int(b.Value())
	})

	var bestLength = 0
	var left, right = 0, 1

	for ; right < len(cards); right++ {
		if cards[right-1].Value()+1 == cards[right].Value() {
			bestLength = right - left + 1
		} else {
			left = right
		}
	}

	if bestLength < 3 {
		bestLength = 0
	}

	return uint8(bestLength)
}

func TailgateRepetitions(cards []deck.Card) (points uint8) {
	var count = 1
	for ; count < len(cards); count++ {
		if cards[len(cards)-1-count].Value() != cards[len(cards)-1].Value() {
			break
		}
	}
	return uint8(count * (count - 1))
}

func AnyRepetitions(cards []deck.Card) (points uint8) {
	var counts = make(map[uint8]int)
	for _, c := range cards {
		counts[c.Value()]++
	}

	for _, v := range counts {
		points += uint8(v * (v - 1))
	}
	return points
}

func Flush(cards []deck.Card, extra deck.Card, isCrib bool) (points uint8) {
	assert.Assert(len(cards) == 4, "expected flush hands to contain 4 cards")

	for _, c := range cards {
		if c.Color() != cards[0].Color() {
			return 0
		}
	}

	if extra.Color() == cards[0].Color() {
		return 5
	}
	if isCrib {
		return 0
	}

	return 4
}

func HisHeels(extra deck.Card) (points uint8) {
	if extra.Value() == deck.JACK {
		return 1
	}
	return 0
}

func HisNobs(extra deck.Card, cards []deck.Card) (points uint8) {
	for _, c := range cards {
		if c.Value() == deck.JACK && c.Color() == extra.Color() {
			return 1
		}
	}
	return 0
}
