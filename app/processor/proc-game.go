package processor

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"stars-server/app/models"
)

type game interface {
	TxGetStellarBodies(ctx context.Context, filter models.StellarBodyFilter) (map[uuid.UUID]*models.StellarBody, error)
	TxGetStellarBodiesStockpiles(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID][]models.Stockpile, error)
	TxGetSystems(ctx context.Context, filter models.StellarBodyFilter) ([]models.System, error)
	TxGetStellarBodyTypes(ctx context.Context) ([]models.StellarBodyType, error)
	TxGetGames(ctx context.Context) ([]models.DBGame, error)

	TxUpdateStellarBodiesMovement(ctx context.Context, duration time.Duration, gameID int) error
	TxUpdateWorldTime(ctx context.Context, duration time.Duration, gameID int) error

	DeleteGame(ctx context.Context, id int) error
}

func (p *Processor) GetStellarBodies(ctx context.Context, filter models.StellarBodyFilter) (map[uuid.UUID]*models.StellarBody, error) {
	var (
		err           error
		stellarBodies map[uuid.UUID]*models.StellarBody
		stockpiles    map[uuid.UUID][]models.Stockpile
	)

	if err = p.db.WithTx(ctx, func(ctx context.Context) error {
		if stellarBodies, err = p.db.TxGetStellarBodies(ctx, filter); err != nil {
			return fmt.Errorf("p.db.GetStellarBodies: %w", err)
		}

		ids := make([]uuid.UUID, 0, len(stellarBodies))
		for _, sb := range stellarBodies {
			ids = append(ids, sb.ID)
		}

		if stockpiles, err = p.db.TxGetStellarBodiesStockpiles(ctx, ids); err != nil {
			return fmt.Errorf("p.db.TxGetStellarBodiesStockpiles: %w", err)
		}

		for key, sp := range stockpiles {
			stellarBodies[key].Stockpiles = sp
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("p.db.WithTx: %w", err)
	}

	return stellarBodies, nil
}

func (p *Processor) GetSystems(ctx context.Context, filter models.StellarBodyFilter) ([]models.System, error) {
	var (
		result []models.System
		err    error
	)

	if err = p.db.WithTx(ctx, func(ctx context.Context) error {
		result, err = p.db.TxGetSystems(ctx, filter)
		if err != nil {
			return fmt.Errorf("p.db.TxGetSystems: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("p.db.WithTx: %w", err)
	}

	return result, nil
}

func (p *Processor) GetStellarBodyTypes(ctx context.Context) ([]models.StellarBodyType, error) {
	var (
		result []models.StellarBodyType
		err    error
	)

	if err = p.db.WithTx(ctx, func(ctx context.Context) error {
		result, err = p.db.TxGetStellarBodyTypes(ctx)
		if err != nil {
			return fmt.Errorf("p.db.TxGetStellarBodyTypes: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("p.db.WithTx: %w", err)
	}

	return result, nil
}

func (p *Processor) GetGames(ctx context.Context) ([]models.DBGame, error) {
	var (
		err    error
		result []models.DBGame
	)

	if err = p.db.WithTx(ctx, func(ctx context.Context) error {
		result, err = p.db.TxGetGames(ctx)
		if err != nil {
			return fmt.Errorf("p.db.TxGetGames: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("p.db.WithTx: %w", err)
	}

	return result, nil
}

func (p *Processor) GameTick(ctx context.Context, duration time.Duration) error {
	var (
		err    error
		gameID int
	)

	gameID, err = models.GetGameIDFromCtx(ctx)

	if err = p.db.WithTx(ctx, func(ctx context.Context) error {
		// получить список событий
		// выполнять по порядку. До завершения или события со стопом
		// подвигать stellar bodies
		if err = p.db.TxUpdateStellarBodiesMovement(ctx, duration, gameID); err != nil {
			return fmt.Errorf(
				"p.db.TxUpdateStellarBodiesMovement(duration: %v, gameID: %d): %w", duration, gameID, err)
		}

		// подвигать экономику
		// обновить время мира
		if err = p.db.TxUpdateWorldTime(ctx, duration, gameID); err != nil {
			return fmt.Errorf("p.db.TxUpdateWorldTime: %w", err)
		}
		// закончить

		return nil
	}); err != nil {
		return fmt.Errorf("p.db.WithTx: %w", err)
	}

	return nil
}

func (p *Processor) DeleteWholeGame(ctx context.Context, id int) error {
	if err := p.db.DeleteGame(ctx, id); err != nil {
		return fmt.Errorf("p.db.DeleteGame(%d): %w", id, err)
	}

	return nil
}
