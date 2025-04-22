package main

import (
	"fmt"

	"github.com/julienr1/cribbage/internal/game"
)

func main() {
	fmt.Println("cribbage")

	players := []*game.Hand{
		game.NewHand(),
		game.NewHand(),
	}
	game := game.New(players)

	fmt.Println(game)
}
