package deck

import (
	"errors"
	"fmt"
	"strings"
)

var InvalidCardErr = errors.New("invalid card")

type Color uint8

const (
	SPADES Color = iota
	CLUBS
	DIAMONDS
	HEARTS
)

// color --> 2 bits
// value: 1-13 --> 4 bits
// bitmap: cc--vvvv
type Card uint8

func NewCard(value uint8, color Color) (Card, error) {
	if value == 0 || value > 13 {
		return 0, InvalidCardErr
	}
	if color != SPADES && color != CLUBS && color != DIAMONDS && color != HEARTS {
		return 0, InvalidCardErr
	}

	c := uint8(color) & 0b11
	v := value & 0b1111
	return Card(c<<6 | v), nil
}

func (c Card) Color() Color {
	return Color(c >> 6)
}

func (c Card) Value() uint8 {
	return uint8(c & 0b1111)
}

func (c Card) String() string {
	var label strings.Builder

	switch c.Value() {
	case 1:
		label.WriteString("A")
	case 11:
		label.WriteString("J")
	case 12:
		label.WriteString("Q")
	case 13:
		label.WriteString("K")
	default:
		label.WriteString(fmt.Sprint(c.Value()))
	}

	label.WriteString(" of ")

	switch c.Color() {
	case SPADES:
		label.WriteString("spades")
	case CLUBS:
		label.WriteString("clubs")
	case DIAMONDS:
		label.WriteString("diamonds")
	case HEARTS:
		label.WriteString("hearts")
	}

	return label.String()
}
