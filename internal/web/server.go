package web

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienr1/cribbage/internal/utils"
	activegame "github.com/julienr1/cribbage/internal/web/active-game"
	"github.com/julienr1/cribbage/internal/web/templates"
)

func cookie(id string) http.Cookie {
	return http.Cookie{
		Name:     "playerId",
		Value:    id,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func Run() {
	games := activegame.NewRegistry()
	games.Set("8uGAs", activegame.New("8uGAs"))

	var upgrader = websocket.Upgrader{}

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("GET /public/res/", http.StripPrefix("/public/res/", fs))

	http.HandleFunc("GET /{gameId}/players/{playerId}", func(w http.ResponseWriter, r *http.Request) {
		game, ok := games.Get(r.PathValue("gameId"))
		if ok == false {
			http.Error(w, "invalid game id", http.StatusBadRequest)
			return
		}

		c, _ := r.Cookie("playerId")
		fmt.Println("cookie:", c.Value)

		playerId := r.PathValue("playerId")
		for _, player := range game.Players {
			if player.Id == playerId {
				templates.Player(player, game.GetPlayerStatus(player.Id)).Render(context.Background(), w)
				return
			}
		}

		http.Error(w, "unknown player", http.StatusBadRequest)
	})

	http.HandleFunc("GET /{gameId}/ws", func(w http.ResponseWriter, r *http.Request) {
		gameId := r.PathValue("gameId")
		if len(gameId) == 0 || games.Contains(gameId) == false {
			http.Error(w, "invalid game id", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		defer conn.Close()
		playerId := r.URL.Query().Get("player-id")
		game, err := games.RegisterConnection(gameId, playerId, conn)
		if err != nil {
			log.Println(err)
			return
		}

		defer games.UnregisterConnection(game, conn)
		game.Handle(conn)
	})

	http.HandleFunc("GET /{gameId}", func(w http.ResponseWriter, r *http.Request) {
		gameId := r.PathValue("gameId")

		if len(gameId) == 0 || games.Contains(gameId) == false {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		ctx := context.WithValue(context.Background(), "player-id", r.Header.Get("X-player-id"))

		c := cookie(utils.Id(6))
		http.SetCookie(w, &c)

		game, _ := games.Get(gameId)
		templates.Game(game).Render(ctx, w)
	})

	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		id, err := utils.UniqueId(5, games)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		games.Set(id, activegame.New(id))
		http.Redirect(w, r, fmt.Sprintf("/%s", id), http.StatusFound)
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
