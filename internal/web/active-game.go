package web

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
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

	log.Println("received connection for game with id", gameId)
	game.connections = append(game.connections, conn)
	game.cancelationId.Add(1)

	return game, nil
}

func (registry *GameRegistry) UnregisterConnection(g *ActiveGame, conn *websocket.Conn) {
	if i := slices.Index(g.connections, conn); i != -1 {
		g.connections = slices.Delete(g.connections, i, i+1)
	}

	if len(g.connections) == 0 {
		go func() {
			cancelationId := g.cancelationId.Load()
			log.Println("No more connections on game", g.id, ", scheduled to be deleted in 1 minute.")

			time.Sleep(time.Minute)

			if g.cancelationId.Load() == cancelationId {
				registry.Delete(g.id)
				log.Println("Game", g.id, "was deleted.")
			} else {
				log.Println("Game", g.id, "was not deleted as a new connection was detected meanwhile.")
			}
		}()
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
