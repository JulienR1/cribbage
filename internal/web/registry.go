package web

import "sync"

type GameRegistry struct {
	games map[string]*ActiveGame
	lock  sync.Mutex
}

func NewGameRegistry() *GameRegistry {
	return &GameRegistry{games: make(map[string]*ActiveGame)}
}

func (registry *GameRegistry) Set(gameId string, game *ActiveGame) {
	registry.lock.Lock()
	defer registry.lock.Unlock()
	registry.games[gameId] = game
}

func (registry *GameRegistry) Get(gameId string) (*ActiveGame, bool) {
	registry.lock.Lock()
	defer registry.lock.Unlock()

	game, ok := registry.games[gameId]
	return game, ok
}

func (registry *GameRegistry) Contains(gameId string) bool {
	_, ok := registry.Get(gameId)
	return ok
}

func (registry *GameRegistry) Delete(gameId string) {
	registry.lock.Lock()
	defer registry.lock.Unlock()
	delete(registry.games, gameId)
}
