// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// CreateUser implements createUser operation.
//
// POST /user
func (UnimplementedHandler) CreateUser(ctx context.Context, req *UserCreate) (r CreateUserRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GameTick implements gameTick operation.
//
// PUT /games/{game_id}/tick
func (UnimplementedHandler) GameTick(ctx context.Context, params GameTickParams) (r GameTickRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetBodies implements getBodies operation.
//
// GET /games/{game_id}/bodies
func (UnimplementedHandler) GetBodies(ctx context.Context, params GetBodiesParams) (r GetBodiesRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetGames implements getGames operation.
//
// GET /games
func (UnimplementedHandler) GetGames(ctx context.Context) (r GetGamesRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetSystems implements getSystems operation.
//
// GET /games/{game_id}/systems
func (UnimplementedHandler) GetSystems(ctx context.Context, params GetSystemsParams) (r GetSystemsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetTypes implements getTypes operation.
//
// GET /types
func (UnimplementedHandler) GetTypes(ctx context.Context) (r GetTypesRes, _ error) {
	return r, ht.ErrNotImplemented
}
