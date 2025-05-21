package activegame

import (
	"fmt"
	"log"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/julienr1/cribbage/internal/assert"
	"github.com/julienr1/cribbage/internal/game"
)

type ActiveGame struct {
	id      string
	game    *game.Game
	players game.Players

	ch            <-chan string
	cancelationId atomic.Int32
	sessions      map[string][]*websocket.Conn
}

func New(id string) *ActiveGame {
	sessions := make(map[string][]*websocket.Conn)
	return &ActiveGame{id: id, sessions: sessions}
}

func (game *ActiveGame) OnPlayerCountChange() {
	for _, player := range game.sessions {
		for _, connection := range player {
			write(connection, fmt.Sprintf("player-count:%d", len(game.players)))
		}
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

func write(conn *websocket.Conn, message string) {
	w, err := conn.NextWriter(websocket.TextMessage)
	assert.AssertE(err)
	defer w.Close()
	fmt.Fprint(w, message)
}
