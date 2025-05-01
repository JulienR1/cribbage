package game

type GameState int

const (
	deal GameState = iota
	crib
	extra
	play
	score
	done
)

func (state GameState) String() string {
	switch state {
	case deal:
		return "deal"
	case crib:
		return "crib"
	case extra:
		return "extra"
	case play:
		return "play"
	case score:
		return "score"
	default:
		return "unknown"
	}
}
