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
	Id      string
	Players game.Players

	game *game.Game

	ch            <-chan string
	cancelationId atomic.Int32
	sessions      map[string][]*websocket.Conn
}

func New(id string) *ActiveGame {
	sessions := make(map[string][]*websocket.Conn)
	return &ActiveGame{Id: id, sessions: sessions}
}

func (game *ActiveGame) OnPlayerChange(playerId string) {
	for _, player := range game.sessions {
		for _, connection := range player {
			write(connection, "player-change", playerId)
		}
	}
}

func (g *ActiveGame) GetPlayerStatus(playerId string) PlayerStatus {
	if connections, ok := g.sessions[playerId]; ok == false {
		return Unknown
	} else if len(connections) == 0 {
		return Disconnected
	}
	return Connected
}

func (g *ActiveGame) Handle(conn *websocket.Conn) {
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if messageType == websocket.TextMessage {
			log.Println("Received text payload from websocket:", message)
			continue
		}

		assert.Assert(len(message) > 0, "websocket message was empty")
		opcode := message[0]
		data := message[1:]

		log.Println("ws:", opcode, data)

		switch opcode {
		case 0:
		}
	}
}

func write(conn *websocket.Conn, title, message string) {
	w, err := conn.NextWriter(websocket.TextMessage)
	assert.AssertE(err)
	defer w.Close()
	fmt.Fprint(w, fmt.Sprintf("%s:%s", title, message))
}
