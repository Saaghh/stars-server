package handlers

import (
	"context"
	"net/http"
	"stars-server/app/models"
	"time"

	"go.uber.org/zap"
	"stars-server/app/generated/api-server"
)

type game interface {
	GetGames(ctx context.Context) ([]models.DBGame, error)
	GetSystems(context.Context, models.StellarBodyFilter) ([]models.System, error)
	GetStellarBodyTypes(context.Context) ([]models.StellarBodyType, error)
	GetStellarBodies(context.Context, models.StellarBodyFilter) ([]models.StellarBody, error)

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

func (h *Handlers) GetGames(ctx context.Context) (api.GetGamesRes, error) {
	games, err := h.proc.GetGames(ctx)
	if err != nil {
		zap.L().With(zap.Error(err)).Error("GetGames/h.proc.GetGames")
		return &api.Error{
			Code:        http.StatusInternalServerError,
			Description: "internal error",
		}, nil
	}

	result := make(api.GetGamesOKApplicationJSON, 0, len(games))
	for _, gameItem := range games {
		result = append(result, api.Game{
			ID:        gameItem.ID,
			OwnerID:   gameItem.OwnerID,
			WorldTime: gameItem.WorldTime,
			Name:      gameItem.Name,
		})
	}

	return &result, nil
}

func (h *Handlers) GetSystems(ctx context.Context, params api.GetSystemsParams) (api.GetSystemsRes, error) {
	filter := models.StellarBodyFilter{
		Game: params.GameID,
	}

	systems, err := h.proc.GetSystems(ctx, filter)
	if err != nil {
		zap.L().With(zap.Error(err)).Error("GetSystems/h.proc.GetSystems")
		return &api.Error{
			Code:        http.StatusInternalServerError,
			Description: "internal error",
		}, nil
	}

	result := make(api.GetSystemsOKApplicationJSON, 0, len(systems))
	for _, system := range systems {
		result = append(result, api.System{
			ID:     system.ID,
			GameID: system.GameID,
			Name:   system.Name,
		})
	}

	return &result, nil
}

func (h *Handlers) GetTypes(ctx context.Context) (api.GetTypesRes, error) {
	types, err := h.proc.GetStellarBodyTypes(ctx)
	if err != nil {
		zap.L().With(zap.Error(err)).Error("GetTypes/h.proc.GetStellarBodyTypes")
		return &api.Error{
			Code:        http.StatusInternalServerError,
			Description: "internal error",
		}, nil
	}

	result := make(api.GetTypesOKApplicationJSON, 0, len(types))
	for _, typeItem := range types {
		result = append(result, api.StellarBodyType{
			ID:   typeItem.ID,
			Name: typeItem.Name,
		})
	}

	return &result, nil
}

func (h *Handlers) GetBodies(ctx context.Context, params api.GetBodiesParams) (api.GetBodiesRes, error) {
	filter := models.StellarBodyFilter{
		Game: params.GameID,
	}

	bodies, err := h.proc.GetStellarBodies(ctx, filter)
	if err != nil {
		zap.L().With(zap.Error(err)).Error("GetBodies/h.proc.GetStellarBodies")
		return &api.Error{
			Code:        http.StatusInternalServerError,
			Description: "internal error",
		}, nil
	}

	result := make(api.GetBodiesOKApplicationJSON, 0, len(bodies))
	for _, body := range bodies {
		result = append(result, stellarBodyToAPI(body))
	}

	return &result, nil
}

func stellarBodyToAPI(body models.StellarBody) api.StellarBody {
	apiBody := api.StellarBody{
		ID:       body.ID,
		SystemID: body.SystemID,
		Name:     body.Name,
		TypeID:   body.TypeID,
		Mass:     body.Mass,
		Diameter: body.Diameter,
		ParentBody: api.NilUUID{
			Null: body.ParentBodyID == nil,
		},
		OrbitalRadius: api.NilFloat64{
			Null: body.OrbitalRadius == nil,
		},
		Angle: api.NilFloat64{
			Null: body.Angle == nil,
		},
		AngleSpeed: api.NilFloat64{
			Null: body.AngleSpeed == nil,
		},
		LinearSpeed: api.NilFloat64{
			Null: body.LinearSpeed == nil,
		},
		CoordinateX: api.NilFloat64{
			Null: body.CoordinateX == nil,
		},
		CoordinateY: api.NilFloat64{
			Null: body.CoordinateY == nil,
		},
	}

	if !apiBody.ParentBody.Null {
		apiBody.ParentBody.Value = *body.ParentBodyID
	}

	if !apiBody.OrbitalRadius.Null {
		apiBody.OrbitalRadius.Value = *body.OrbitalRadius
	}

	if !apiBody.Angle.Null {
		apiBody.Angle.Value = *body.Angle
	}

	if !apiBody.AngleSpeed.Null {
		apiBody.AngleSpeed.Value = *body.AngleSpeed
	}

	if !apiBody.LinearSpeed.Null {
		apiBody.LinearSpeed.Value = *body.LinearSpeed
	}

	if !apiBody.CoordinateX.Null {
		apiBody.CoordinateX.Value = *body.CoordinateX
	}

	if !apiBody.CoordinateY.Null {
		apiBody.CoordinateY.Value = *body.CoordinateY
	}

	return apiBody
}
