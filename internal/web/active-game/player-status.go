package activegame

type PlayerStatus int

const (
	Connected PlayerStatus = iota
	Disconnected
	Unknown
)
