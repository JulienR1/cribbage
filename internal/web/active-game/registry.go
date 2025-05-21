package activegame

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/gorilla/websocket"
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
	if ok == false {
		return nil, UnknownGameErr
	}

	var player *game.Player = nil
	for _, p := range g.players {
		if p.Id == playerId {
			player = p
		}
	}

	if player == nil {
		player = game.NewPlayer(g.players)
		g.players = append(g.players, player)
	}

	g.cancelationId.Add(1)
	g.sessions[player.Id] = append(g.sessions[player.Id], conn)

	log.Println("player", player.Id, "connected to game", gameId)

	write(conn, fmt.Sprintf("game-id:%s", gameId))
	write(conn, fmt.Sprintf("player-id:%s", player.Id))
	g.OnPlayerCountChange()

	return g, nil
}

func (registry *GameRegistry) UnregisterConnection(game *ActiveGame, conn *websocket.Conn) {
	for playerId, connections := range game.sessions {
		if i := slices.Index(connections, conn); i != -1 {
			game.sessions[playerId] = slices.Delete(connections, i, i+1)
			game.OnPlayerCountChange()
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
			log.Println("No more connections on game", game.id, ", scheduled to be deleted in 1 minute.")

			time.Sleep(time.Minute)

			if game.cancelationId.Load() == cancelationId {
				registry.Delete(game.id)
				log.Println("Game", game.id, "was deleted.")
			} else {
				log.Println("Game", game.id, "was not deleted as a new connection was detected meanwhile.")
			}
		}()
	}
}
