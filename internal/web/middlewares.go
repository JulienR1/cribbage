package web

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/julienr1/cribbage/internal/web/templates"
)

func logger(w http.ResponseWriter, r *http.Request, next func(*http.Request)) {
	log.Println(r.Method, r.URL.Path)
	next(r)
}

func validateGame(w http.ResponseWriter, r *http.Request, next func(*http.Request)) {
	gameId := r.PathValue("gameId")
	if games.Contains(gameId) == false {
		http.Error(w, "game does not exist", http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(r.Context(), templates.GAME_ID, gameId)
	next(r.WithContext(ctx))
}

func authenticate(w http.ResponseWriter, r *http.Request, next func(*http.Request)) {
	gameId := templates.GameId(r.Context())
	if len(gameId) == 0 {
		http.Error(w, "no game is currently ongoing", http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie(gameId)
	if err != nil && errors.Is(err, http.ErrNoCookie) == false {
		http.Error(w, "user not detected", http.StatusUnauthorized)
		return
	}

	userId := cookie.Value
	game, _ := games.Get(gameId)
	if game.Players.Contains(userId) == false {
		http.Error(w, "player is not in specifed game", http.StatusForbidden)
		return
	}

	ctx := context.WithValue(r.Context(), templates.USER_ID, userId)
	next(r.WithContext(ctx))
}

func selectedPlayer(w http.ResponseWriter, r *http.Request, next func(*http.Request)) {
	playerId := r.PathValue("playerId")

	g, _ := games.Get(templates.GameId(r.Context()))
	if g.Players.Contains(playerId) == false {
		http.Error(w, "player is not in specified game", http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(r.Context(), templates.PLAYER_ID, playerId)
	next(r.WithContext(ctx))
}
