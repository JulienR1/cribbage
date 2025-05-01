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

func (cards Hand) String() string {
	var cardStrs []string
	for _, c := range cards {
		cardStrs = append(cardStrs, c.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(cardStrs, " | "))
}
