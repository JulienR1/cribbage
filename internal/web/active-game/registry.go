package activegame

import (
	"errors"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienr1/cribbage/internal/assert"
	"github.com/julienr1/cribbage/internal/game"
)

var IncompleteHandshakeErr = errors.New("could not complete game handshake")
var UnknownGameErr = errors.New("could not find game with specified id")

type GameRegistry struct {
	games map[string]*ActiveGame
	lock  sync.Mutex
}

func NewRegistry() *GameRegistry {
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

func (registry *GameRegistry) RegisterConnection(gameId, playerId string, conn *websocket.Conn) (*ActiveGame, error) {
	g, ok := registry.Get(gameId)
	assert.Assert(ok, UnknownGameErr)

	index := slices.IndexFunc(g.Players, func(p *game.Player) bool { return p.Id == playerId })
	assert.Assert(index >= 0, "expected player to be created in game before being accessed")
	player := g.Players[index]

	g.cancelationId.Add(1)
	g.sessions[player.Id] = append(g.sessions[player.Id], conn)

	log.Println("player", player.Id, "connected to game", gameId)

	write(conn, "game-id", gameId)
	write(conn, "player-id", player.Id)
	g.OnPlayerChange(player.Id)

	return g, nil
}

func (registry *GameRegistry) UnregisterConnection(game *ActiveGame, conn *websocket.Conn) {
	for playerId, connections := range game.sessions {
		if i := slices.Index(connections, conn); i != -1 {
			game.sessions[playerId] = slices.Delete(connections, i, i+1)
			game.OnPlayerChange(playerId)
			break
		}
	}

	var connectionCount = 0
	for _, connections := range game.sessions {
		connectionCount += len(connections)
	}

	if connectionCount == 0 {
		go func() {
			cancelationId := game.cancelationId.Load()
			log.Println("No more connections on game", game.Id, ", scheduled to be deleted in 1 minute.")

			time.Sleep(time.Minute)

			if game.cancelationId.Load() == cancelationId {
				registry.Delete(game.Id)
				log.Println("Game", game.Id, "was deleted.")
			} else {
				log.Println("Game", game.Id, "was not deleted as a new connection was detected meanwhile.")
			}
		}()
	}
}
