package handlers

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
	"stars-server/app/generated/api-server"
)

type game interface {
	GameTick(ctx context.Context, duration time.Duration) error
}

func (h *Handlers) GameTick(ctx context.Context, params api.GameTickParams) (api.GameTickRes, error) {
	err := h.proc.GameTick(ctx, params.Duration)
	if err != nil {
		zap.L().With(zap.Error(err)).Error("GameTick/h.proc.GameTick")
		return &api.GameTickInternalServerError{
			Code:        http.StatusInternalServerError,
			Description: "internal error",
		}, nil
	}

	return &api.GameTickOK{}, nil
}
