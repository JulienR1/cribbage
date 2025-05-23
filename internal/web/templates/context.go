package templates

import "context"

const GAME_ID = "game-id"
const USER_ID = "user-id"
const PLAYER_ID = "player-id"

func GameId(ctx context.Context) string        { return str(ctx, GAME_ID) }
func CurrentUserId(ctx context.Context) string { return str(ctx, USER_ID) }
func PlayerId(ctx context.Context) string      { return str(ctx, PLAYER_ID) }

func str(ctx context.Context, key string) string {
	if v, ok := ctx.Value(key).(string); ok {
		return v
	}
	return ""
}
