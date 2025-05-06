package game

import (
	"log"
	"slices"

	"github.com/julienr1/cribbage/internal/assert"
	"github.com/julienr1/cribbage/internal/deck"
)

func CountIs(target, count uint8) uint8 {
	if count == target {
		log.Printf("%d: +2\n", count)
		return 2
	}
	return 0
}

func Fifteen(cards []deck.Card) (points uint8) {
	assert.Assert(len(cards) < 8, "too many cards in hand to calculate combinaisons")

	for combinaison := range uint(1 << len(cards)) {
		var sum, idx = uint8(0), uint(0)

		for (1<<idx) <= combinaison && idx < uint(len(cards)) {
			if combinaison&(1<<idx) > 0 {
				sum += cards[idx].Points()
			}
			if sum > 15 {
				break
			}

			idx++
		}

		if sum == 15 {
			points += 2
		}
	}

	if points > 0 {
		log.Printf("15 (%dx): +%d\n", points/2, points)
	}

	return points
}

func LastPlayed(playedCard *deck.Card, playing, lastWhoPlayed *Player) (points uint8) {
	if playedCard == nil && playing == lastWhoPlayed {
		log.Println("Go!: +1")
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

	if bestLength > 0 {
		log.Printf("Run (%s): +%d\n", cards[len(cards)-bestLength:], bestLength)
	}

	return uint8(bestLength)
}

func AnySeries(cards []deck.Card) (points uint8) {
	copiedCards := make([]deck.Card, len(cards))
	for i, c := range cards {
		copiedCards[i] = c
	}

	slices.SortFunc(copiedCards, func(a, b deck.Card) int {
		return int(a.Value()) - int(b.Value())
	})

	var bestLength = 0
	var left, right = 0, 1

	for ; right < len(copiedCards); right++ {
		if copiedCards[right-1].Value()+1 == copiedCards[right].Value() {
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

func AnySeriesLog(cards []deck.Card) (points uint8) {
	points = AnySeries(cards)
	if points > 0 {
		log.Printf("Series: +%d\n", points)
	}
	return points
}

func TailgateRepetitions(cards []deck.Card) (points uint8) {
	var count = 1
	for ; count < len(cards); count++ {
		if cards[len(cards)-1-count].Value() != cards[len(cards)-1].Value() {
			break
		}
	}

	if count > 1 {
		log.Printf("Repetitions (%d): +%d\n", count, count*(count-1))
	}

	return uint8(count * (count - 1))
}

func AnyRepetitions(cards []deck.Card) (points uint8) {
	var counts = make(map[uint8]int)
	for _, c := range cards {
		counts[c.Value()]++
	}

	for k, v := range counts {
		if v > 1 {
			points += uint8(v * (v - 1))
			log.Printf("Repetitions (%d %ds): +%d\n", v, k, v*(v-1))
		}
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
		log.Println("Flush (w/ extra): +5")
		return 5
	}
	if isCrib {
		return 0
	}

	log.Println("Flush: +4")
	return 4
}

func HisHeels(extra deck.Card) (points uint8) {
	if extra.Value() == deck.JACK {
		log.Println("His heels! +1")
		return 1
	}
	return 0
}

func HisNobs(extra deck.Card, cards []deck.Card) (points uint8) {
	for _, c := range cards {
		if c.Value() == deck.JACK && c.Color() == extra.Color() {
			log.Println("His nobs! +1")
			return 1
		}
	}
	return 0
}
