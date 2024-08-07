package models

import (
	"context"
)

type GameInfoKey struct{}

func GetGameIDFromCtx(ctx context.Context) (int, error) {
	id, ok := ctx.Value(GameInfoKey{}).(int)
	if !ok {
		return 0, ErrNoDataInCtx
	}

	return id, nil
}
