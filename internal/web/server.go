package web

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienr1/cribbage/internal/web/templates"
)

func Run() {
	games := make(GameRegistry)
	games["8uGAs"] = &ActiveGame{id: "8uGAs"}

	var upgrader = websocket.Upgrader{}

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("GET /public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("GET /ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		defer conn.Close()
		game, err := games.RegisterConnection(conn)
		if err != nil {
			log.Println(err)
			return
		}

		defer games.UnregisterConnection(game, conn)
		game.Handle(conn)
	})

	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		id, err := UniqueId(5, games)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		games[id] = &ActiveGame{id: id}
		http.Redirect(w, r, fmt.Sprintf("/%s", id), http.StatusFound)
	})

	http.HandleFunc("GET /{gameId}", func(w http.ResponseWriter, r *http.Request) {
		gameId := r.PathValue("gameId")
		if _, ok := games[gameId]; ok == false || len(gameId) == 0 {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		templates.Game().Render(context.Background(), w)
	})

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		gameId := r.FormValue("game-id")
		if len(gameId) > 0 {
			http.Redirect(w, r, fmt.Sprintf("/%s", gameId), http.StatusFound)
			return
		}

		templates.Index().Render(context.Background(), w)
	})

	log.Println("Listening on http://localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
