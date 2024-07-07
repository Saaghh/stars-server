package processor

import (
	"context"
	"fmt"
	"time"

	"stars-server/app/models"
)

type game interface {
	TxGetStellarBodies(ctx context.Context, filter models.StellarBodyFilter) ([]models.StellarBody, error)
	TxGetSystems(ctx context.Context, filter models.StellarBodyFilter) ([]models.System, error)
	TxGetStellarBodyTypes(ctx context.Context) ([]models.StellarBodyType, error)
	TxUpdateStellarBodiesMovement(ctx context.Context, duration time.Duration, gameID int) error
}

func (p *Processor) GetStellarBodies(ctx context.Context, filter models.StellarBodyFilter) ([]models.StellarBody, error) {
	var (
		result []models.StellarBody
		err    error
	)

	if err = p.db.WithTx(ctx, func(ctx context.Context) error {
		result, err = p.db.TxGetStellarBodies(ctx, filter)
		if err != nil {
			return fmt.Errorf("p.GetStellarBodies: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("p.db.WithTx: %w", err)
	}

	return result, nil
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

		// TODO: add filter for game
		if err = p.db.TxUpdateStellarBodiesMovement(ctx, duration, gameID); err != nil {
			return fmt.Errorf(
				"p.db.TxUpdateStellarBodiesMovement(duration: %v, gameID: %d): %w", duration, gameID, err)
		}

		// подвигать экономику
		// закончить

		return nil
	}); err != nil {
		return fmt.Errorf("p.db.WithTx: %w", err)
	}

	return nil
}
