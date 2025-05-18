package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienr1/cribbage/internal/game"
)

func Run() {
	games := make(map[string]*game.Game)

	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		id, err := UniqueId(5, games)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		games[id] = nil
		fmt.Fprint(w, id)
	})

	http.HandleFunc("GET /{gameId}", func(w http.ResponseWriter, r *http.Request) {
		gameId := r.PathValue("gameId")

		// TODO: check if game actually exists

		if len(gameId) == 0 {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		fmt.Fprint(w, gameId)
	})

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "todo: some page")
	})

	log.Println("Listening on http://localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
