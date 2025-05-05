package game

import (
	"fmt"
	"strings"

	"github.com/julienr1/cribbage/internal/deck"
)

type Hand []deck.Card

type CardStack interface {
	AddCard(card deck.Card) error
}

func (cards Hand) Playable(count uint8) (playable []deck.Card) {
	for _, c := range cards {
		if c.Points()+count <= 31 {
			playable = append(playable, c)
		}
	}
	return playable
}

func (cards Hand) String() string {
	var cardStrs []string
	for _, c := range cards {
		cardStrs = append(cardStrs, c.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(cardStrs, " | "))
}
