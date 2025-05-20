package web

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienr1/cribbage/internal/assert"
	"github.com/julienr1/cribbage/internal/game"
)

var IncompleteHandshakeErr = errors.New("could not complete game handshake")
var UnknownGameErr = errors.New("could not find game with specified id")

type ActiveGame struct {
	id   string
	game *game.Game

	ch            <-chan string
	connections   []*websocket.Conn
	cancelationId atomic.Int32
}

func (registry *GameRegistry) RegisterConnection(conn *websocket.Conn) (*ActiveGame, error) {
	messageType, message, err := conn.ReadMessage()
	if messageType != websocket.TextMessage || err != nil {
		return nil, IncompleteHandshakeErr
	}

	gameId := string(message)
	game, ok := registry.Get(gameId)
	if ok == false {
		return nil, UnknownGameErr
	}

	game.cancelationId.Add(1)
	game.connections = append(game.connections, conn)
	log.Println("received connection for game with id", gameId)

	game.OnPlayerCountChange()

	return game, nil
}

func (registry *GameRegistry) UnregisterConnection(game *ActiveGame, conn *websocket.Conn) {
	if i := slices.Index(game.connections, conn); i != -1 {
		game.connections = slices.Delete(game.connections, i, i+1)
		game.OnPlayerCountChange()
	}

	if len(game.connections) == 0 {
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

func (game *ActiveGame) OnPlayerCountChange() {
	for _, player := range game.connections {
		w, err := player.NextWriter(websocket.TextMessage)
		assert.AssertE(err)
		defer w.Close()
		fmt.Fprintf(w, "player-count:%d", len(game.connections))
	}
}

func (g *ActiveGame) Handle(conn *websocket.Conn) {
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		switch messageType {
		case websocket.BinaryMessage:
			fmt.Println("ws:", message)
		case websocket.TextMessage:
			fmt.Println("ws:", string(message))
		}
	}
}
