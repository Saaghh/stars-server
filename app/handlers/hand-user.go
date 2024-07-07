package handlers

import (
	"context"
	"net/http"
	"stars-server/app/generated/api-server"
	"stars-server/app/models"
)

type user interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
}

func (h *Handlers) CreateUser(ctx context.Context, req *api.UserCreate) (api.CreateUserRes, error) {
	user := models.User{
		Username:     req.Username,
		PasswordHash: req.Password,
	}

	user, err := h.proc.CreateUser(ctx, user)
	if err != nil {
		return &api.CreateUserInternalServerError{
			Code:        http.StatusInternalServerError,
			Description: "internal error",
		}, nil
	}

	return &api.CreateUserOK{}, nil
}
