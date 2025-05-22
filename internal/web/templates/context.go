package templates

import "context"

func CurrentPlayerId(ctx context.Context) string {
	if playerId, ok := ctx.Value("player-id").(string); ok {
		return playerId
	}
	return ""
}
