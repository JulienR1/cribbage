package web

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienr1/cribbage/internal/assert"
	"github.com/julienr1/cribbage/internal/utils"
	activegame "github.com/julienr1/cribbage/internal/web/active-game"
	"github.com/julienr1/cribbage/internal/web/middleware"
	"github.com/julienr1/cribbage/internal/web/templates"
)

var games = activegame.NewRegistry()
var upgrader = websocket.Upgrader{}

func Run() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("GET /public/res/", http.StripPrefix("/public/res/", fs))

	middleware.New(logger, validateGame, authenticate, selectedPlayer).
		HandleFunc("GET /{gameId}/players/{playerId}", func(w http.ResponseWriter, r *http.Request) {
			playerId := templates.PlayerId(r.Context())
			game, _ := games.Get(templates.GameId(r.Context()))

			for _, player := range game.Players {
				if player.Id == playerId {
					templates.Player(player, game.GetPlayerStatus(player.Id)).Render(context.Background(), w)
					return
				}
			}

			http.Error(w, "unknown player", http.StatusBadRequest)
		})

	middleware.New(logger, validateGame, authenticate).
		HandleFunc("GET /{gameId}/ws", func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println(err.Error())
				return
			}
			defer conn.Close()

			gameId := templates.GameId(r.Context())
			userId := templates.CurrentUserId(r.Context())

			game, err := games.RegisterConnection(gameId, userId, conn)
			if err != nil {
				log.Println(err)
				return
			}

			defer games.UnregisterConnection(game, conn)
			game.Handle(conn)
		})

	middleware.New(logger, validateGame).
		HandleFunc("GET /{gameId}", func(w http.ResponseWriter, r *http.Request) {
			gameId := templates.GameId(r.Context())
			game, _ := games.Get(gameId)

			var playerId = ""
			if c, err := r.Cookie(gameId); errors.Is(err, http.ErrNoCookie) == false {
				playerId = c.Value
			} else {
				playerId, err = utils.UniqueId(8, game.Players)
				assert.AssertE(err)
			}

			game.WithPlayer(playerId)

			cookie := http.Cookie{
				Name:     gameId,
				Value:    playerId,
				Path:     fmt.Sprintf("/%s", gameId),
				Secure:   true,
				HttpOnly: true,
				MaxAge:   3600, // 1 hour
				SameSite: http.SameSiteStrictMode,
			}
			http.SetCookie(w, &cookie)

			w.Header().Set("Cache-Control", "no-store")
			ctx := context.WithValue(r.Context(), templates.USER_ID, playerId)
			templates.Game(game).Render(ctx, w)
		})

	middleware.New(logger).
		HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
			id, err := utils.UniqueId(5, games)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, err)
				return
			}

			games.Set(id, activegame.New(id))
			http.Redirect(w, r, fmt.Sprintf("/%s", id), http.StatusFound)
		})

	middleware.New(logger).
		HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			gameId := r.FormValue("game-id")
			if len(gameId) > 0 {
				http.Redirect(w, r, fmt.Sprintf("/%s", gameId), http.StatusFound)
				return
			}

			templates.Index().Render(r.Context(), w)
		})

	log.Println("Listening on http://localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
