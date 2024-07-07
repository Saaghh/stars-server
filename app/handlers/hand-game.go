package handlers

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"stars-server/app/generated/api-server"
	"time"
)

type game interface {
	GameTick(ctx context.Context, duration time.Duration) error
}

func (h *Handlers) GameTick(ctx context.Context, params api.GameTickParams) (api.GameTickRes, error) {
	res := fmt.Sprintf("%v", params.DateTime)

	zap.L().Info(res)

	return nil, nil
}
