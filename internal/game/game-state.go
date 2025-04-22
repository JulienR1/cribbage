package game

type GameState int

const (
	deal GameState = iota
	crib
	extra
	play
	score
)
